package network

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestStaticBackoffPolicy_Backoff(t *testing.T) {
	policy := staticBackoffPolicy{backoffDuration: 1 * time.Second}
	backoffDuration := policy.getBackoffDuration()
	assert.Equal(t, backoffDuration, 1*time.Second)
}

func TestIncrementalBackoffPolicy(t *testing.T) {
	policy := incrementalBackoffPolicy{
		backoffDuration:   3 * time.Second,
		incrementDuration: 5 * time.Second,
	}
	attempt0Duration := policy.getBackoffDuration(0)
	assert.Equal(t, 3*time.Second, attempt0Duration)
	attempt1Duration := policy.getBackoffDuration(1)
	assert.Equal(t, (3*time.Second)+(5*time.Second), attempt1Duration)
	attempt2Duration := policy.getBackoffDuration(2)
	assert.Equal(t, (3*time.Second)+(10*time.Second), attempt2Duration)
}

func submitConcurrentlyAndWait(fn func(ctx context.Context) error, times int) {
	goFn := func(w *sync.WaitGroup) {
		_ = fn(context.Background())
		w.Done()
	}

	wg := sync.WaitGroup{}
	wg.Add(times)
	for i := 0; i < times; i++ {
		go goFn(&wg)
	}
	wg.Wait()
}

func TestWaitingThrottler_Acquire(t *testing.T) {
	throttleLimit := 2
	thr := NewWaitingThrottler(throttleLimit)

	assertAcquire(t, thr, throttleLimit, 1*time.Second)
}

func TestBackoffThrottler_Acquire(t *testing.T) {
	backoffDuration := 1 * time.Second
	throttleLimit := 2
	thr := NewStaticBackoffThrottler(throttleLimit, backoffDuration)

	assertAcquire(t, thr, throttleLimit, backoffDuration)
}

func TestConcurrency(t *testing.T) {
	called := 0
	thr := NewWaitingThrottler(1)
	throttlingFn := func(ctx context.Context) {
		_ = thr.Acquire(ctx)
		if ctx.Err() == context.Canceled {
			fmt.Println(ctx.Err())
			return
		}
		called++
	}

	rootCtx := context.Background()
	var cancels []func()
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithCancel(rootCtx)
		cancels = append(cancels, cancel)
		go throttlingFn(ctx)
	}

	time.Sleep(500 * time.Millisecond)

	for _, cancelFn := range cancels {
		cancelFn()
	}

	assert.Equal(t, 1, called)

}

func assertAcquire(t *testing.T, thr Throttler, throttleLimit int, backoffDuration time.Duration) {
	t1 := time.Now()
	submitConcurrentlyAndWait(thr.Acquire, throttleLimit)
	t2 := time.Now()

	assert.Less(t, t2.Sub(t1), backoffDuration)

	time.Sleep(backoffDuration)

	// Create throttler with 2 aps limit and static backoff of 1 second
	// We create a waitgroup for 4 actions, submit all 4 using goroutines and wait for them
	// to complete. Since it is 2 actions allowed per second, the total time for all 4
	// concurrent requests should be around 1 second.
	t1 = time.Now()
	submitConcurrentlyAndWait(thr.Acquire, throttleLimit*2)
	t2 = time.Now()

	delayedDuration := t2.Sub(t1)

	assert.Greater(t, delayedDuration, backoffDuration)
	assert.Less(t, delayedDuration, 2*backoffDuration)

	time.Sleep(backoffDuration)

	t1 = time.Now()
	submitConcurrentlyAndWait(thr.Acquire, throttleLimit)
	t2 = time.Now()

	finalDuration := t2.Sub(t1)

	assert.Greater(t, delayedDuration, finalDuration)
}
