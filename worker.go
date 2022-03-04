package forrest

import (
	"log"
	"sync/atomic"
	"time"
)

type Worker struct {
	taskChan chan func()
	pool *Pool
	touchFishTime time.Time
}

func NewWorker(taskChan chan func(), pool *Pool) *Worker {
	w := &Worker{
		taskChan: taskChan,
		pool: pool,
	}
	w.run()
	return w
}

func (w *Worker) run()  {
	go func() {
		atomic.AddInt32(&w.pool.runningWorker, 1)
		defer func() {
			atomic.AddInt32(&w.pool.runningWorker, -1)
			if err := recover(); err != nil {
				log.Printf("worker exit with panic : %s", err)
			}
		}()
		for f := range w.taskChan {
			if f == nil {
				return
			}
			f()
			w.touchFishTime = time.Now()
			w.pool.revertWorker(w)
		}
	}()
}