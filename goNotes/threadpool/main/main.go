package main

import (
	"fmt"
	"goNotes/threadpool"
	"goNotes/threadpool/thread_func"
	"reflect"
)

func main() {
	tp := &threadpool.ThreadPool{}
	tp.Init(3, 10, 6, thread_func.Add)
	tp.Start()
	go func() {
		for i := 0; i < 30; i++ {
			tp.AddTask([]interface{}{1, i})
		}
		tp.Stop()
	}()

	go func() {
		tp.Wait()
	}()

	for a := range tp.Result {
		fmt.Println(a.([]reflect.Value)[0].Int())
	}

}
