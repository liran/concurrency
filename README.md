## USAGE
```go
package main

import (
	"fmt"
	"time"

	"github.com/liran/concurrency"
)

func main() {
	// Create a thread pool that can allocate up to 10 threads
	pool := concurrency.New(10, func(params ...interface{}) {
		n := params[0].(int)
		w := params[1].(string)
		time.Sleep(time.Second)
		fmt.Println(w, n)
	})
	defer pool.Close()

	for i := 0; i < 100; i++ {
		fmt.Println("a:", i)
		pool.Process(i, "hello")
	}
	pool.Wait()

	for i := 0; i < 20; i++ {
		fmt.Println("b:", i)
		pool.Process(i, "world")
	}

	// Multiple threads will be notified at the same time
	go func() {
		pool.Wait()
		fmt.Println("wait on go thread ")
	}()

	pool.Wait()
	fmt.Println("wait on main thread")
}
```
