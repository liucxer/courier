package worker

import (
	"context"
	"sync"
)

func NewWorker(process func(ctx context.Context) error, numWorkers int) *Worker {
	return &Worker{
		numWorkers: numWorkers,
		process:    process,
	}
}

type Worker struct {
	numWorkers int
	process    func(ctx context.Context) error
	wg         sync.WaitGroup
}

func (mq *Worker) Start(ctx context.Context) {
	mq.wg.Add(mq.numWorkers)

	for i := 0; i < mq.numWorkers; i++ {
		go func(workerID int) {
			defer mq.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					if err := mq.process(ctx); err != nil {
						continue
					}
				}
			}
		}(i)
	}

	mq.wg.Wait()
}
