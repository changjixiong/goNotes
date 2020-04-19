package thread_func

import "time"

func Add(a, b int) int {
	time.Sleep(3 * time.Second)
	return a + b
}
