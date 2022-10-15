package concurrency

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSingleTask(t *testing.T) {
	uptime := time.Now()
	pool := New(10, func(params ...any) {
		task := params[0].(int)
		time.Sleep(time.Second)
		fmt.Println("run", task)
	})
	defer pool.Close()
	for i := 0; i < 10; i++ {
		pool.Process(i)
	}

	go func() {
		pool.Wait()
		log.Println("goroutine")
	}()

	pool.Wait()
	pool.Wait()
	pool.Close()
	for i := 0; i < 10; i++ {
		pool.Process(i)
	}
	pool.Wait()
	pool.Close()
	log.Printf("complete, uptime: %s", time.Since(uptime))
}

func TestQ(t *testing.T) {
	q := make(chan any)
	go func() {
		log.Println("b")
		log.Println(len(q))

		<-q
	}()

	log.Println("a")
	q <- 1

	close(q)
}
