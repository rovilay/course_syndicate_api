package main

import "sync"

func main() {
	wg := sync.WaitGroup{}
	app := App{}
	app.initialize(1)

	wg.Add(2) // 2 is the number of go routines

	go func() {
		defer wg.Done()
		app.start()
	}()

	go func() {
		defer wg.Done()
		app.sc.SyncSchedulesWorker()
	}()

	wg.Wait()
}
