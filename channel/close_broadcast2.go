//https://talks.golang.org/2012/10things.slide
package main

import (
	"fmt"
	"time"
)

func worker(i int, ch chan Work, quit chan struct{}) {
	var quitting bool
	for {
		select {
		case w := <-ch:
			if quitting {
				w.Refuse()
				fmt.Println("worker", i, "refused", w)
				break
			}
			w.Do()
			fmt.Println("worker", i, "processed", w)
		case <-quit:
			fmt.Println("worker", i, "quitting")
			quitting = true
		}
	}
}

func main() {
	ch, quit := make(chan Work), make(chan struct{})
	go makeWork(ch)
	for i := 0; i < 4; i++ {
		go worker(i, ch, quit)
	}
	time.Sleep(5 * time.Second)
	close(quit)
	time.Sleep(2 * time.Second)
}
