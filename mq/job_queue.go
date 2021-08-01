package mq

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/liucxer/courier/courier"
	"github.com/liucxer/courier/metax"
	"github.com/liucxer/courier/mq/worker"
)

func NewJobBoard(taskMgr TaskMgr) *JobBoard {
	return &JobBoard{
		taskMgr: taskMgr,
	}
}

type JobBoard struct {
	taskMgr TaskMgr
}

func (b *JobBoard) Dispatch(channel string, task *Task) error {
	if task == nil {
		return nil
	}
	return b.taskMgr.Push(channel, task)
}

type JobWorkerOpts struct {
	Channel    string
	NumWorkers int
	OnFinish   func(ctx context.Context, task *Task)
}

func NewJobWorker(taskMgr TaskMgr, opts JobWorkerOpts) *JobWorker {
	return &JobWorker{
		JobWorkerOpts: opts,
		taskMgr:       taskMgr,
	}
}

type JobWorker struct {
	JobWorkerOpts
	operators       sync.Map
	taskMgr         TaskMgr
	worker          *worker.Worker
	contextInjector func(ctx context.Context) context.Context
}

func (w *JobWorker) Context() context.Context {
	ctx := context.Background()
	if w.contextInjector != nil {
		return w.contextInjector(context.Background())
	}
	return ctx
}
func (w *JobWorker) getOperatorMeta(typ string) (*courier.OperatorMeta, error) {
	op, ok := w.operators.Load(typ)
	if !ok {
		return nil, fmt.Errorf("missing operator %s", typ)
	}
	return op.(*courier.OperatorMeta), nil
}

func (w *JobWorker) WithContextInjector(contextInjector func(ctx context.Context) context.Context) *JobWorker {
	return &JobWorker{
		JobWorkerOpts:   w.JobWorkerOpts,
		operators:       sync.Map{},
		taskMgr:         w.taskMgr,
		worker:          w.worker,
		contextInjector: contextInjector,
	}
}

func (w *JobWorker) Serve(router *courier.Router) error {
	w.Register(router)

	chStop := make(chan os.Signal, 1)
	signal.Notify(chStop, os.Interrupt, syscall.SIGTERM)
	w.worker = worker.NewWorker(w.process, w.NumWorkers)
	go func() {
		w.worker.Start(w.Context())
	}()

	<-chStop

	return http.ErrServerClosed
}

func (w *JobWorker) Register(router *courier.Router) {
	for _, route := range router.Routes() {
		factories := route.OperatorFactories()
		if len(factories) != 1 {
			continue
		}
		f := factories[0]
		w.operators.Store(f.Type.Name(), f)
	}
}

func (w *JobWorker) process(ctx context.Context) (err error) {
	task, err := w.taskMgr.Shift(w.Channel)
	if err != nil || task == nil {
		return nil
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %s; calltrace:%s", fmt.Sprint(e), string(debug.Stack()))
		}

		if err != nil {
			task.Stage = STAGE_FAILED
		} else {
			task.Stage = STAGE_SUCCESS
		}

		if w.OnFinish != nil {
			w.OnFinish(ctx, task)
		}
	}()

	opMeta, e := w.getOperatorMeta(task.Subject)
	if e != nil {
		err = e
		return
	}

	op := opMeta.New()

	if task.Argv != nil {
		if writer, ok := op.(io.Writer); ok {
			if _, e := writer.Write(task.Argv); e != nil {
				err = e
				return
			}
		}
	}

	meta := metax.ParseMeta(task.Id)
	meta.Add("task", w.Channel+"#"+task.Subject)

	ctx = metax.ContextWithMeta(ctx, meta)

	if _, e := op.Output(ctx); e != nil {
		err = e
		return
	}
	return
}
