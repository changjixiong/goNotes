package threadpool

import (
	"reflect"
	"sync"
)

type ThreadPool struct {
	function  interface{}
	tasks     chan []interface{}
	Result    chan interface{}
	ThreadNum int
	wg        sync.WaitGroup
	stop      chan interface{}
}

func (tp *ThreadPool) Init(threadNum, maxPendingTaskNum, reusltNum int, function interface{}) {
	tp.ThreadNum = threadNum
	tp.tasks = make(chan []interface{}, maxPendingTaskNum)
	tp.Result = make(chan interface{}, reusltNum)
	tp.stop = make(chan interface{})
	tp.function = function
}
func (tp *ThreadPool) AddTask(task []interface{}) {
	tp.tasks <- task
}
func (tp *ThreadPool) Start() {
	for i := 0; i < tp.ThreadNum; i++ {
		tp.wg.Add(1)
		go func() {
			defer tp.wg.Done()
			for {
				task, ok := <-tp.tasks
				if ok {
					obj := reflect.ValueOf(tp.function)
					if obj.Kind() == reflect.Func {
						params := make([]reflect.Value, len(task))
						for i, value := range task {
							params[i] = reflect.ValueOf(value)

						}
						tp.Result <- obj.Call(params)
					}
				} else {
					break
				}
			}
		}()
	}
}

func (tp *ThreadPool) Stop() {
	tp.stop <- nil
}

func (tp *ThreadPool) Wait() {
	<-tp.stop
	close(tp.tasks)
	tp.wg.Wait()
	close(tp.Result)
}
