package taskrunner

import (
	"time"
)

type Worker struct {
	ticker *time.Ticker
	runner *TaskRunner
}

func NewWorker(interval time.Duration, runner *TaskRunner) *Worker {
	worker := &Worker{}
	worker.ticker = time.NewTicker(interval * time.Second)
	worker.runner = runner
	return worker
}

func (w *Worker) start() {
	for {
		select {
		case <-w.ticker.C:
			// 每间隔指定时间启动一次删除任务
			go w.runner.start()
		}
	}
}

func StartWorker() {
	runner := NewTaskRunner(BATCH_DELETE_CNT, true, VideoClearDispatcher, VideoClearExecutor)
	worker := NewWorker(BATCH_DELETE_INTERVAL, runner)
	go worker.start()
}
