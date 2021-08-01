package mq_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/liucxer/courier/metax"
	"github.com/stretchr/testify/require"

	"github.com/liucxer/courier/courier"
	"github.com/liucxer/courier/mq/memtaskmgr"
	"github.com/liucxer/courier/mq/redistaskmgr"
	"github.com/gomodule/redigo/redis"

	"github.com/liucxer/courier/mq"
)

var taskMgr = memtaskmgr.NewMemTaskMgr()
var taskMgrRedis = redistaskmgr.NewRedisTaskMgr(r)

var channel = "TEST"

var r = redistaskmgr.RedisOperatorFromPool(&redis.Pool{
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379", redis.DialDatabase(3))
	},
	MaxIdle:     10,
	MaxActive:   10,
	IdleTimeout: 10 * time.Second,
	Wait:        true,
})

func init() {
	_ = taskMgr.Destroy(channel)
	_ = taskMgrRedis.Destroy(channel)
}

type A struct {
}

func (a A) Output(ctx context.Context) (interface{}, error) {
	fmt.Println(metax.MetaFromContext(ctx))
	return nil, nil
}

type B struct {
	bytes.Buffer
}

func (b B) Output(ctx context.Context) (interface{}, error) {
	fmt.Println(metax.MetaFromContext(ctx), b.String())
	return nil, nil
}

var router = courier.NewRouter()

func TestJobQueue(t *testing.T) {
	list := []mq.TaskMgr{
		taskMgr,
		taskMgrRedis,
	}

	for i := range list {
		taskMgr := list[i]

		jobBoard := mq.NewJobBoard(taskMgr)

		n := 100

		for i := 0; i < n; i++ {
			for j := 0; j < 5; j++ {
				_ = jobBoard.Dispatch("TEST", mq.NewTask("A", []byte("A"), fmt.Sprintf("A%d", i)))
				_ = jobBoard.Dispatch("TEST", mq.NewTask("B", []byte("B"), fmt.Sprintf("B%d", i)))
			}
		}

		router.Register(courier.NewRouter(&A{}))
		router.Register(courier.NewRouter(&B{}))

		jobWorker := mq.NewJobWorker(taskMgr, mq.JobWorkerOpts{
			Channel:    "TEST",
			NumWorkers: 2,
			OnFinish: func(ctx context.Context, task *mq.Task) {
				require.Equal(t, mq.STAGE_SUCCESS, task.Stage)
			},
		})

		go func() {
			err := jobWorker.Serve(router)
			fmt.Println(err)
		}()

		time.Sleep(400 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)

		time.Sleep(500 * time.Millisecond)
	}
}
