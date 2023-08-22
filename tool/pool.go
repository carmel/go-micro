package tool

import (
	"sync"
)

// type Work func()

type WorkPool struct {
	queue chan struct{}
	wg    *sync.WaitGroup
}

func NewPool(size int, wg *sync.WaitGroup) *WorkPool {
	return &WorkPool{queue: make(chan struct{}, size), wg: wg}
}

func (p *WorkPool) Acquire() {
	defer p.wg.Add(1)
	p.queue <- struct{}{}
}

func (p *WorkPool) Release() {
	defer p.wg.Done()
	<-p.queue
}

func (p *WorkPool) Wait() {
	defer close(p.queue)
	p.wg.Wait()
}
