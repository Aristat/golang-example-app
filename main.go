package main

import (
	"fmt"
	"sync"
)

//
//func goRoutineA(a <-chan int) {
//	val := <-a
//	fmt.Println("goRoutineA received the data", val)
//}
//
//func goRoutineB(a <-chan int) {
//	val := <-a
//	fmt.Println("goRoutineB received the data", val)
//}
//
//func main() {
//	ch := make(chan int)
//	go goRoutineA(ch)
//	go goRoutineB(ch)
//	ch <- 3
//	time.Sleep(time.Second * 1)
//}

//func main() {
//	cmd.Execute()
//}

//func main() {
//	ch := make(chan int)
//	//done := make(chan struct{}) // create DONE channel
//
//	var wg sync.WaitGroup
//
//	numbers := []int{1, 2, 3}
//
//	wg.Add(len(numbers))
//
//	for _, n := range numbers {
//		go func(n int) {
//			defer wg.Done()
//			ch <- n
//		}(n)
//	}
//
//	go func() {
//		//defer close(done) // close done channel to tell that all jobs is done
//		for c := range ch {
//			fmt.Printf("routine start %v\n", c)
//			time.Sleep(1 * time.Second) // for better understanding
//			fmt.Printf("routine done %v\n", c)
//		}
//	}()
//
//	wg.Wait()
//	//close(ch) // after all routines push payload then we close channel to release RANGE
//	//<-done // wait when all jobs is done
//}

func main() {
	//b := &goback.SimpleBackoff{
	//	Min:    100 * time.Millisecond,
	//	Max:    60 * time.Second,
	//	Factor: 2,
	//}
	//cb := b
	var wg sync.WaitGroup

	theMine := [5]string{"ore1", "ore2", "ore3", "ore4", "ore5"}
	oreChan := make(chan string)
	done := make(chan struct{})

	wg.Add(1)

	go func(mine [5]string) {
		defer wg.Done()

		for _, item := range mine {
			fmt.Println("send mine", item)
			oreChan <- item //отправка
		}
	}(theMine)

	go func() {
		defer close(done)
		for foundOre := range oreChan {
			fmt.Printf("Miner: Received %v from finder\n", foundOre)
		}
	}()

	//go func() {
	//	for {
	//		fmt.Println("call")
	//
	//		d, e := cb.NextAttempt()
	//		if e != nil {
	//			fmt.Printf("backoff error: %v", logger.Args(e))
	//		}
	//
	//		time.Sleep(d)
	//		continue
	//	}
	//}()

	wg.Wait()
	close(oreChan)
	<-done
}
