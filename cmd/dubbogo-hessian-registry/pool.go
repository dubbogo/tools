package main

import (
	"sync"
)

// Task 任务
type Task interface {
	Execute()
}

// Pool 线程池
type Pool struct {
	wg *sync.WaitGroup

	taskQueue chan Task
}

func NewPool(max int) *Pool {
	return &Pool{
		wg:        &sync.WaitGroup{},
		taskQueue: make(chan Task, max),
	}
}

func (p Pool) Execute(t Task) {
	p.wg.Add(1)
	p.taskQueue <- t
	go func() {
		t.Execute()
		<-p.taskQueue
		p.wg.Done()
	}()
}

func (p Pool) Wait() {
	p.wg.Wait()
}
