package main

import "sync"

func main() {
	wg := sync.WaitGroup{}
	app := App{}
	app.initialize(1)

	wg.Add(2)

	go func() {
		defer wg.Done()
		app.start(&wg)
	}()

	go func() {
		defer wg.Done()
		app.sc.SyncSchedulesWorker()
	}()

	wg.Wait()
}
