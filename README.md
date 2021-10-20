## USAGE
```go
package main

import (
	"fmt"
	"time"

	"github.com/liran/concurrent"
)

func main() {
	pool := concurrent.New(10, func(params ...interface{}) {
		n := params[0].(int)
		w := params[1].(string)
		time.Sleep(time.Second)
		fmt.Println(w, n)
	})
	defer pool.Close()

	for i := 0; i < 10; i++ {
		fmt.Println("i:", i)
		pool.Process(i, "hello")
	}

	pool.Wait()
}

```