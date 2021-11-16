package workers

import (
	"context"
	"encoding/json"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/tasks"
	"go.uber.org/zap"
	"sync"
	"time"
)

type WorkerPoolJobStoreInterface interface {
	GetJobToExecute(ctx context.Context, maxCount int) ([]models.JobStoreRow, error)
	ExecuteJob(ctx context.Context, jobID string) error
	IncreaseCounter(ctx context.Context, jobID string, count int) error
}

type Job struct {
	ID    string
	Func  func(ctx context.Context) error
	Count int
}

type WorkerPool struct {
	jobStore         WorkerPoolJobStoreInterface
	taskStore        *tasks.TaskStore
	numOfWorkers     int
	inputCh          chan Job
	log              *zap.SugaredLogger
	maxJobRetryCount int
}

func New(jobStore WorkerPoolJobStoreInterface, taskStore *tasks.TaskStore, cfg *configurations.ConfigWorkerPool,
	log *zap.SugaredLogger) *WorkerPool {
	wp := &WorkerPool{
		jobStore:         jobStore,
		taskStore:        taskStore,
		numOfWorkers:     cfg.NumOfWorkers,
		inputCh:          make(chan Job, cfg.PoolBuffer),
		log:              log,
		maxJobRetryCount: cfg.MaxJobRetryCount,
	}
	return wp
}

func (wp *WorkerPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for i := 0; i < wp.numOfWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			wp.log.Infof("Worker #%v start", i)
		outer:
			for {
				select {
				case job := <-wp.inputCh:
					err := job.Func(ctx)
					if err != nil {
						wp.log.Errorf("Error on worker #%v: %v\n", i, err.Error())
						err = wp.jobStore.IncreaseCounter(ctx, job.ID, job.Count)
						if err != nil {
							wp.log.Errorf("Error with increase job counter with job:%v error:%v", job.ID, err.Error())
						}
						continue
					}
					err = wp.jobStore.ExecuteJob(ctx, job.ID)
					if err != nil {
						wp.log.Errorf("Error with executing job: %v, error: %v", job.ID, err.Error())
					}
				case <-ctx.Done():
					break outer
				}

			}
			wp.log.Infof("Worker #%v close", i)
			wg.Done()
		}(i)
	}
	sch := wp.scheduler(ctx)
	wg.Wait()
	close(wp.inputCh)
	sch.Stop()
}

func (wp *WorkerPool) push(task Job) {
	wp.inputCh <- task
}

func (wp *WorkerPool) scheduler(ctx context.Context) *time.Ticker {
	ticker := time.NewTicker(time.Second * 5)
	wp.log.Info("start scheduler")
	go func() {
		for {
			select {
			case <-ticker.C:
				wp.transferTaskToWorkerPool(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
	return ticker
}

func (wp *WorkerPool) transferTaskToWorkerPool(ctx context.Context) {
	jobs, err := wp.jobStore.GetJobToExecute(ctx, wp.maxJobRetryCount)
	if err != nil {
		wp.log.Errorf("Error occured in getting task in worker pool: %v", err.Error())
		return
	}
	for _, job := range jobs {

		task, ok := wp.taskStore.MapOfTask[job.Type]

		if !ok {
			wp.log.Errorf("Get job of unknown type: %v", job.Type)
			continue
		}
		parameters := make(map[string]string)
		err := json.Unmarshal([]byte(job.Parameters), &parameters)
		if err != nil {
			wp.log.Errorf("Error with parce parameters, job_id: %v, err: %v", job.ID, err.Error())
			continue
		}
		function, err := task.CreateFunction(parameters)
		if err != nil {
			wp.log.Errorf("Wrong paramenters for function, job_id: %v, err: %v", job.ID, err.Error())
			continue
		}
		jobToPush := Job{
			ID:   job.ID,
			Func: function,
		}

		wp.push(jobToPush)
	}
}
