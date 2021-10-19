package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func TestConcurrent(t *testing.T) {
	pool := New(10, func(input interface{}) {
		n := input.(int)
		time.Sleep(time.Second)
		fmt.Println("n:", n)
	})
	defer pool.Close()

	for i := 0; i < 10; i++ {
		fmt.Println("i:", i)
		pool.Process(i)
	}
	time.Sleep(3 * time.Second)
}
