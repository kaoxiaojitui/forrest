package forrest

import (
	"sync/atomic"
	"time"
)

const (
	Working int32 = iota
	Closed
)

const (
	DefaultCoreSize = 10
	DefaultCapacity = 100
	LayoffsTime = time.Second * 2
)

type Pool struct {
	status int32

	core int32
	capacity int32
	runningWorker int32

	workers chan *Worker

	lockChan chan struct{}
}

func NewDefaultPool() *Pool {
	return NewPool(DefaultCoreSize, DefaultCapacity)
}

func NewPool(core, capacity int32) *Pool {
	if core < 0 {
		core = DefaultCoreSize
	}
	if capacity < 0 {
		capacity = DefaultCapacity
	}
	if core > capacity {
		core = capacity
	}
	p := &Pool{
		core: core,
		capacity: capacity,
		workers: make(chan *Worker, capacity),
		status: Working,
		lockChan: make(chan struct{}, 1),
	}
	go p.winterComing()
	return p
}

func (p *Pool) Running() int32 {
	return atomic.LoadInt32(&p.runningWorker)
}

func (p *Pool) lock()  {
	p.lockChan <- struct{}{}
}

func (p *Pool) unlock() {
	<- p.lockChan
}

func (p *Pool) setPoolStatus(status int32)  {
	atomic.StoreInt32(&p.status, status)
}

func (p *Pool) getPoolStatus() int32 {
	return atomic.LoadInt32(&p.status)
}

func (p *Pool) Submit(f func()) error {
	if p.getPoolStatus() == Closed {
		return PoolHasBeenClosedError
	}

	// get worker
	worker := p.retrieveWorker()
	if worker == nil {
		return TooManyGoroutineBlockError
	}
	worker.taskChan <- f
	return nil
}

func (p *Pool) retrieveWorker() *Worker {
	for {
		p.lock()
		select {
		case worker, ok := <- p.workers:
			p.unlock()
			if !ok {
				p.setPoolStatus(Closed)
				return nil
			}
			return worker
		default:
			p.unlock()
			if p.Running() < p.capacity {
				return NewWorker(make(chan func()), p)
			}
		}
	}
}

func (p *Pool) revertWorker(worker *Worker)  {
	if p.getPoolStatus() == Closed {
		return
	}
	p.lock()
	p.workers <- worker
	p.unlock()
}

func (p *Pool) winterComing() {
	timer := time.NewTicker(LayoffsTime)
	for range timer.C {
		if p.getPoolStatus() == Closed {
			return
		}
		if p.Running() > p.core {
			p.lock()
			select {
			case worker, ok := <-p.workers:
				if !ok {
					p.setPoolStatus(Closed)
					p.unlock()
					return
				}
				if time.Now().Add(-LayoffsTime).Before(worker.touchFishTime) {
					p.unlock()
					p.revertWorker(worker)
				} else {
					atomic.AddInt32(&p.runningWorker, -1)
					worker = nil
					p.unlock()
				}
			}
		}
	}
}




