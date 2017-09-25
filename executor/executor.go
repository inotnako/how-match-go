package executor

import "sync"

type Deltaer interface {
	Add(delta int64)
}

type noopDeltaer struct{}

func (n noopDeltaer) Add(delta int64) {}

var stats Deltaer = noopDeltaer{}

func New(maxSize int) Executor {
	return &executorImpl{
		wg:   &sync.WaitGroup{},
		pool: make(chan *struct{}, maxSize),
	}
}

type Executor interface {
	Do(func(string), string)
	Wait()
}

type executorImpl struct {
	wg   *sync.WaitGroup
	pool chan *struct{}
}

func (e *executorImpl) Do(f func(string), path string) {
	// if we can't send msg to pool, we wait
	e.pool <- &struct{}{}
	e.wg.Add(1)
	// for stats
	stats.Add(1)

	go func() {
		f(path)

		// free place from pool
		<-e.pool
		// for graceful shutdown all tasks
		e.wg.Done()

		stats.Add(-1)
	}()
}

func (e *executorImpl) Wait() {
	e.wg.Wait()
}
