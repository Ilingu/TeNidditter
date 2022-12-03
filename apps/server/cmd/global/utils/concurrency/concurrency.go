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

// WaitGroup wrapper with a mutex
func NewMultipleRoutineWaitGroup() *MultipleRoutineWaitGroup {
	return &MultipleRoutineWaitGroup{}
}

// 1) Lock Mutex
//
// 2) Add delta to wait group (see func (*sync.WaitGroup).Add(delta int) docs)
//
// 3) Unlock Mutex
func (mwg *MultipleRoutineWaitGroup) Add(delta int) {
	mwg.mut.Lock()
	mwg.wg.Add(delta)
	mwg.mut.Unlock()
}

// 1) Lock Mutex
//
// 2) func (*sync.WaitGroup).Done()
//
// 3) Unlock Mutex
func (mwg *MultipleRoutineWaitGroup) Done() {
	mwg.mut.Lock()
	mwg.wg.Done()
	mwg.mut.Unlock()
}

// 1) Lock Mutex
//
// 2) func (*sync.WaitGroup).Wait()
//
// 3) Unlock Mutex
func (mwg *MultipleRoutineWaitGroup) Wait() {
	mwg.mut.Lock()
	mwg.wg.Wait()
	mwg.mut.Unlock()
}
