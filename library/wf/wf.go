package wf

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type PayloadInterface interface {
	Play()
}

type Job struct {
	Payload PayloadInterface
}

type Worker struct {
	WorkerId        string
	Workbench       chan Job
	GWorkbenchQueue chan chan Job
	Finished        chan bool
}

func NewWorker(WorkbenchQueue chan chan Job, Id string) *Worker {
	return &Worker{
		WorkerId:        Id,
		Workbench:       make(chan Job),
		GWorkbenchQueue: WorkbenchQueue,
		Finished:        make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.GWorkbenchQueue <- w.Workbench
			select {
			case wJob := <-w.Workbench:
				wJob.Payload.Play()
			case bFinished := <-w.Finished:
				if true == bFinished {
					return
				}
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.Finished <- true
	}()
}

type Dispatcher struct {
	DispatcherId    string
	MaxWorkers      int
	Workers         []*Worker
	Closed          chan bool
	EndDispatch     chan os.Signal
	GJobQueue       chan Job
	GWorkbenchQueue chan chan Job
}

func NewDispatcher(dispatcherId string, maxWorkers, maxQueue int) *Dispatcher {
	Closed := make(chan bool)
	EndDispatch := make(chan os.Signal)
	JobQueue := make(chan Job, maxQueue)
	WorkbenchQueue := make(chan chan Job, maxWorkers)
	signal.Notify(EndDispatch, syscall.SIGINT, syscall.SIGTERM)
	return &Dispatcher{
		DispatcherId:    dispatcherId,
		MaxWorkers:      maxWorkers,
		Closed:          Closed,
		EndDispatch:     EndDispatch,
		GJobQueue:       JobQueue,
		GWorkbenchQueue: WorkbenchQueue,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.GWorkbenchQueue, fmt.Sprintf("work-%s", strconv.Itoa(i)))
		d.Workers = append(d.Workers, worker)
		worker.Start()
	}
	go d.Dispatch()
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case <-d.EndDispatch:
			close(d.GJobQueue)
		case wJob, Ok := <-d.GJobQueue:
			if true == Ok {
				go func(wJob Job) {
					Workbench := <-d.GWorkbenchQueue
					Workbench <- wJob
				}(wJob)
			} else {
				for _, w := range d.Workers {
					w.Stop()
				}
				d.Closed <- true
				return
			}
		}
	}
}

type WorkFlow struct {
	GDispatch *Dispatcher
}

func Start(dispatchId string, maxWorkers, maxQueue int) *WorkFlow {
	var wf WorkFlow
	wf.GDispatch = NewDispatcher(dispatchId, maxWorkers, maxQueue)
	wf.GDispatch.Run()
	return &wf
}

func (wf *WorkFlow) AddJob(wJob Job) {
	wf.GDispatch.GJobQueue <- wJob
}

func (wf *WorkFlow) Close() {
	closed := <-wf.GDispatch.Closed
	if true == closed {
		g.Log().Println(fmt.Sprintf("Dispatch %s closed ... ", wf.GDispatch.DispatcherId))
	}
}
