package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type fooBarBaz struct {
	foo string
	bar string
	baz string
}

func main() {

	var (
		startTime = time.Now().UnixNano()
	)

	// Example function that calls three other long running functions concurrently. Calling each function sequentially would take ~1500ms. This function will return in ~500ms, the time it takes the longest of the three to return.
	data, err := callFooBarBaz()
	if err != nil {
		// Handle the error however you normally would
		fmt.Println(err)
		return
	}

	// This code will run if there are no errors.
	fmt.Println("do somthing important")

	// On success all values will be shown and the time will be ~500ms showing all three functions executed at the same time.
	fmt.Printf("Foo: %s, Bar: %s, Baz: %s, Time: %vms.\n", data.foo, data.bar, data.baz, (time.Now().UnixNano()-startTime)/1000000)

}

func callFooBarBaz() (fooBarBaz, error) {

	var (
		data = fooBarBaz{}
		// Mutex scoped to this function
		mutex = &sync.Mutex{}
		// Wait group blocking the main thread until other threads have finished
		wg = sync.WaitGroup{}
		// Channel used to send errors
		errChan = make(chan error)
		// This channel will be closed when the wait group finishes indicating all functions returned without error
		noErr = make(chan struct{})
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// The return values are scoped to the anonymous function to prevent a data race if all three go routines are writing to the struct or a shared err.
		foo, e := fooFunc()
		if e != nil {
			// If there is an error, pass it through the channel. This will cause the function to return. See the select statement below.
			errChan <- e
		}
		// Lock data while writing to prevent a data race
		mutex.Lock()
		data.foo = foo
		mutex.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		bar, e := barFunc()
		if e != nil {
			errChan <- e
		}
		mutex.Lock()
		data.bar = bar
		mutex.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		baz, e := bazFunc()
		if e != nil {
			errChan <- e
		}
		mutex.Lock()
		data.baz = baz
		mutex.Unlock()
	}()

	// Put the wait group on its own go routine to prevent it from blocking the select statement. When the wait group finishes the noErr channel closes which triggers its case in the select statement.
	go func() {
		wg.Wait()
		close(noErr)
	}()

	// Select blocks until one of the conditions is met. In this example, it blocks until one call returns with an error or all calls finish without error.
	select {
	case err := <-errChan:
		// If more than one of the functions has an error, only the first one will be returned.
		return fooBarBaz{}, err
	case <-noErr:
	}

	return data, nil
}

func fooFunc() (string, error) {
	time.Sleep(500 * time.Millisecond)
	return returnFunc("fooFunc")
}

func barFunc() (string, error) {
	time.Sleep(500 * time.Millisecond)
	return returnFunc("barFunc")
}

func bazFunc() (string, error) {
	time.Sleep(500 * time.Millisecond)
	return returnFunc("bazFunc")
}

// returnFunc randomly send back a success or error response
func returnFunc(funcName string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	if n > 80 {
		return "fail", fmt.Errorf("%s error", funcName)
	}
	return "success", nil
}
