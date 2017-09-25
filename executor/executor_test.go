package executor

import (
	"testing"
	"time"
)

func doSomething(t *testing.T, counter chan int) func(string) {
	t.Helper()
	return func(d string) {
		dur, err := time.ParseDuration(d)
		if err != nil {
			t.Error(err)
			return
		}

		time.Sleep(dur)
		counter <- 1
	}
}

func TestExecutorImpl_Do(t *testing.T) {
	exec := New(100)

	counter := make(chan int)
	result := make(chan int, 1)
	defer close(result)

	go func() {
		count := 0
		for c := range counter {
			count += c
		}
		result <- count
	}()

	for i := 0; i < 1000; i++ {
		sleepDur := time.Duration(i) * time.Millisecond
		exec.Do(doSomething(t, counter), sleepDur.String())
	}

	exec.Wait()

	close(counter)
	out := <-result

	if out != 1000 {
		t.Errorf(`expected result=1000, got - %d`, out)
	}
}
