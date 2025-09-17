package or

import (
	"fmt"
	"testing"
	"time"
)

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func TestOrClosesOnFirst(t *testing.T) {
	start := time.Now()
	done := Or(sig(300*time.Millisecond), sig(150*time.Millisecond), sig(500*time.Millisecond))
	select {
	case <-done:
	case <-time.After(800 * time.Millisecond):
		t.Fatal("timeout waiting for Or")
	}
	elapsed := time.Since(start)
	if elapsed < 100*time.Millisecond || elapsed > 700*time.Millisecond {
		t.Fatalf("Or closed after %v; want around shortest (≈150ms)", elapsed)
	}
}

func TestOrZeroChannelsIsClosed(t *testing.T) {
	select {
	case <-Or():
	default:
		t.Fatal("Or() must return a closed channel")
	}
}

func TestOrSingleClosedChannel(t *testing.T) {
	ch := make(chan interface{})
	close(ch)
	select {
	case <-Or(ch):
	default:
		t.Fatal("Or(closed) must close immediately")
	}
}

func TestOrIgnoresNil(t *testing.T) {
	start := time.Now()
	<-Or(nil, sig(100*time.Millisecond), nil)
	if time.Since(start) < 80*time.Millisecond {
		t.Fatalf("Or returned too early")
	}
}

func ExampleOr() {
	<-Or(sig(200*time.Millisecond), sig(1*time.Second))
	fmt.Println("done")
}
