package workers

import (
	"context"
	"errors"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/customerrors"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/utils"
	"go.uber.org/zap"
	"sync"
	"time"
)

type WorkerPool struct {
	numOfWorkers int
	inputCh      chan func(ctx context.Context) error
	log          *zap.SugaredLogger
	taskToRepeat utils.FunctionStack
}

func New(numOfWorkers int, buffer int, log *zap.SugaredLogger) *WorkerPool {
	wp := &WorkerPool{
		numOfWorkers: numOfWorkers,
		inputCh:      make(chan func(ctx context.Context) error, buffer),
		log:          log,
		taskToRepeat: utils.NewFunctionStack(),
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
						var repeatError *customerrors.RepeatError
						if errors.As(err, &repeatError) {
							wp.taskToRepeat.Push(f)
						}
					}
				case <-ctx.Done():
					break outer
				}

			}
			wp.log.Infof("Worker #%v close\n", i)
			wg.Done()
		}(i)
	}
	sch := wp.scheduler(ctx)
	wg.Wait()
	close(wp.inputCh)
	sch.Stop()
}

func (wp *WorkerPool) Push(task func(ctx context.Context) error) {
	wp.inputCh <- task
}

func (wp *WorkerPool) scheduler(ctx context.Context) *time.Ticker {
	ticker := time.NewTicker(time.Second * 5)
	wp.log.Info("start scheduler")
	go func() {
		for {
			select {
			case <-ticker.C:
				wp.log.Info("ticker tick")
				wp.pushToRepeat()
			case <-ctx.Done():
				return
			}
		}
	}()
	return ticker
}

func (wp *WorkerPool) pushToRepeat() {
	for {
		f, ok := wp.taskToRepeat.Pop()
		if !ok {
			break
		}
		wp.inputCh <- f
	}
}
