package utils_concurrency

import (
	"sync"
)

type ConditionalVariable struct {
	sync.Mutex
	Cond *sync.Cond
}

func NewConditionalVariable() *ConditionalVariable {
	cond := ConditionalVariable{}
	cond.Cond = sync.NewCond(&cond)
	return &cond
}

type MultipleRoutineWaitGroup struct {
	mut sync.Mutex
	wg  sync.WaitGroup
}

func NewMultipleRoutineWaitGroup() *MultipleRoutineWaitGroup {
	return &MultipleRoutineWaitGroup{}
}

func (mwg *MultipleRoutineWaitGroup) Add(delta int) {
	mwg.mut.Lock()
	mwg.wg.Add(delta)
	mwg.mut.Unlock()
}
func (mwg *MultipleRoutineWaitGroup) Done() {
	mwg.mut.Lock()
	mwg.wg.Done()
	mwg.mut.Unlock()
}
func (mwg *MultipleRoutineWaitGroup) Wait() {
	mwg.mut.Lock()
	mwg.wg.Wait()
	mwg.mut.Unlock()
}
