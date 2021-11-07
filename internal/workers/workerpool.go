package workers

import (
	"context"
	"go.uber.org/zap"
	"sync"
)

type WorkerPool struct {
	numOfWorkers int
	inputCh      chan func(ctx context.Context) error
	log          *zap.SugaredLogger
}

func New(numOfWorkers int, buffer int, log *zap.SugaredLogger) *WorkerPool {
	wp := &WorkerPool{
		numOfWorkers: numOfWorkers,
		inputCh:      make(chan func(ctx context.Context) error, buffer),
		log:          log,
	}
	return wp
}

func (wp *WorkerPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for i := 0; i < wp.numOfWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			wp.log.Infof("Worker #%v start \n", i)
		outer:
			for {
				select {
				case f := <-wp.inputCh:
					err := f(ctx)
					if err != nil {
						wp.log.Errorf("Error on worker #%v: %v\n", i, err.Error())
					}
				case <-ctx.Done():
					break outer
				}

			}
			wp.log.Infof("Worker #%v close\n", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(wp.inputCh)
}

func (wp *WorkerPool) Push(task func(ctx context.Context) error) {
	wp.inputCh <- task
}
