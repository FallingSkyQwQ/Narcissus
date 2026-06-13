package reactive

import (
	"sync"
	"testing"
	"time"
)

func TestSignalBasic(t *testing.T) {
	s := NewSignal(10)
	if s.Get() != 10 {
		t.Errorf("expected 10, got %d", s.Get())
	}

	s.Set(20)
	if s.Get() != 20 {
		t.Errorf("expected 20, got %d", s.Get())
	}
}

func TestSignalSubscribe(t *testing.T) {
	s := NewSignal(0)
	var received []int

	unsubscribe := s.Subscribe(func(v int) {
		received = append(received, v)
	})

	s.Set(1)
	s.Set(2)
	s.Set(3)

	if len(received) != 4 { // Initial + 3 updates
		t.Errorf("expected 4 values, got %d", len(received))
	}

	unsubscribe()
	s.Set(4)

	if len(received) != 4 {
		t.Errorf("expected 4 values after unsubscribe, got %d", len(received))
	}
}

func TestSignalMultipleSubscribers(t *testing.T) {
	s := NewSignal(0)
	var received1, received2 []int

	unsubscribe1 := s.Subscribe(func(v int) {
		received1 = append(received1, v)
	})

	unsubscribe2 := s.Subscribe(func(v int) {
		received2 = append(received2, v)
	})

	s.Set(1)
	s.Set(2)

	if len(received1) != 3 || len(received2) != 3 {
		t.Errorf("expected 3 values each, got %d and %d", len(received1), len(received2))
	}

	unsubscribe1()
	s.Set(3)

	if len(received1) != 3 {
		t.Errorf("expected 3 values for unsubscribed observer, got %d", len(received1))
	}
	if len(received2) != 4 {
		t.Errorf("expected 4 values for active observer, got %d", len(received2))
	}

	unsubscribe2()
}

func TestSignalString(t *testing.T) {
	s := NewSignal("hello")
	if s.Get() != "hello" {
		t.Errorf("expected 'hello', got %q", s.Get())
	}

	s.Set("world")
	if s.Get() != "world" {
		t.Errorf("expected 'world', got %q", s.Get())
	}
}

func TestSignalStruct(t *testing.T) {
	type Point struct {
		X, Y int
	}

	s := NewSignal(Point{X: 1, Y: 2})
	if s.Get().X != 1 || s.Get().Y != 2 {
		t.Errorf("expected Point{1, 2}, got %+v", s.Get())
	}

	s.Set(Point{X: 3, Y: 4})
	if s.Get().X != 3 || s.Get().Y != 4 {
		t.Errorf("expected Point{3, 4}, got %+v", s.Get())
	}
}

func TestSignalConcurrent(t *testing.T) {
	s := NewSignal(0)
	var wg sync.WaitGroup
	numGoroutines := 10
	numIterations := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				s.Set(id*numIterations + j)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				_ = s.Get()
			}
		}()
	}

	// Concurrent subscriptions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				unsub := s.Subscribe(func(v int) {})
				unsub()
			}
		}()
	}

	wg.Wait()
}

func TestComputed(t *testing.T) {
	a := NewSignal(2)
	b := NewSignal(3)

	sum := NewComputed(func() int {
		return a.Get() + b.Get()
	})

	if sum.Get() != 5 {
		t.Errorf("expected 5, got %d", sum.Get())
	}

	// Note: In this simplified implementation, Computed does not auto-update
	// The value is computed once at creation
	a.Set(5)
	// Value should still be 5 (the initial computed value)
	if sum.Get() != 5 {
		t.Errorf("expected 5 (initial computed value), got %d", sum.Get())
	}

	sum.Dispose()
}

func TestComputedRecompute(t *testing.T) {
	a := NewSignal(2)
	b := NewSignal(3)

	sum := NewComputed(func() int {
		return a.Get() + b.Get()
	})

	if sum.Get() != 5 {
		t.Errorf("expected 5, got %d", sum.Get())
	}

	a.Set(5)
	sum.Recompute()
	if sum.Get() != 8 {
		t.Errorf("expected 8 after Recompute, got %d", sum.Get())
	}

	b.Set(10)
	sum.Recompute()
	if sum.Get() != 15 {
		t.Errorf("expected 15 after Recompute, got %d", sum.Get())
	}

	sum.Dispose()
}

func TestComputedSubscribe(t *testing.T) {
	a := NewSignal(1)
	b := NewSignal(2)

	product := NewComputed(func() int {
		return a.Get() * b.Get()
	})

	var received []int
	unsubscribe := product.Subscribe(func(v int) {
		received = append(received, v)
	})

	// Should have initial value
	if len(received) != 1 || received[0] != 2 {
		t.Errorf("expected initial value [2], got %v", received)
	}

	// Note: In this simplified implementation, Computed does not auto-update
	// So changing source signals won't trigger updates
	a.Set(3)
	if len(received) != 1 {
		t.Errorf("expected no new values in simplified implementation, got %v", received)
	}

	unsubscribe()
	product.Dispose()
}

func TestEffectBasic(t *testing.T) {
	var callCount int
	s := NewSignal(0)

	effect := NewEffect(func() {
		callCount++
		_ = s.Get() // Access signal
	})

	// Effect should run immediately on creation
	if callCount != 1 {
		t.Errorf("expected effect to run once on creation, got %d calls", callCount)
	}

	// Note: In this simplified implementation, Effect does not auto-re-run on dependency change
	s.Set(1)
	if callCount != 1 {
		t.Errorf("effect should not auto-run in simplified implementation, got %d calls", callCount)
	}

	effect.Dispose()
}

func TestEffectDispose(t *testing.T) {
	var callCount int

	effect := NewEffect(func() {
		callCount++
	})

	if callCount != 1 {
		t.Errorf("expected effect to run once, got %d calls", callCount)
	}

	effect.Dispose()

	// After dispose, effect should not run again
	initialCount := callCount
	effect.Dispose() // Double dispose should be safe
	if callCount != initialCount {
		t.Errorf("effect should not run after dispose, got %d calls", callCount)
	}
}

func TestEffectCleanup(t *testing.T) {
	var cleanupCount int

	effect := NewEffect(func() {
		// Effect function
	})

	effect.SetCleanup(func() {
		cleanupCount++
	})

	// Cleanup should not be called yet
	if cleanupCount != 0 {
		t.Errorf("cleanup should not be called yet, got %d calls", cleanupCount)
	}

	effect.Dispose()

	// Cleanup should be called on dispose
	if cleanupCount != 1 {
		t.Errorf("cleanup should be called once on dispose, got %d calls", cleanupCount)
	}
}

func TestComputedSetCompute(t *testing.T) {
	a := NewSignal(2)
	b := NewSignal(3)

	// Initial compute: a + b
	sum := NewComputed(func() int {
		return a.Get() + b.Get()
	})

	if sum.Get() != 5 {
		t.Errorf("expected 5, got %d", sum.Get())
	}

	// Change to compute: a * b
	sum.SetCompute(func() int {
		return a.Get() * b.Get()
	})

	if sum.Get() != 6 {
		t.Errorf("expected 6 after SetCompute, got %d", sum.Get())
	}

	sum.Dispose()
}

func TestComputedConcurrent(t *testing.T) {
	s := NewSignal(10)
	c := NewComputed(func() int {
		return s.Get() * 2
	})

	var wg sync.WaitGroup
	numGoroutines := 10
	numIterations := 100

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				_ = c.Get()
			}
		}()
	}

	// Concurrent recomputes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				c.Recompute()
			}
		}()
	}

	// Concurrent signal changes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				s.Set(id*numIterations + j)
			}
		}(i)
	}

	wg.Wait()
	c.Dispose()
}

func TestEffectConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	numEffects := 10

	for i := 0; i < numEffects; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			effect := NewEffect(func() {
				// Effect function
			})
			effect.SetCleanup(func() {
				// Cleanup function
			})
			effect.Dispose()
		}()
	}

	wg.Wait()
}

func TestSignalConcurrentSubscribeSet(t *testing.T) {
	s := NewSignal(0)
	var wg sync.WaitGroup
	numSubscribers := 10
	numSets := 20

	// Track deliveries from each subscriber
	type delivery struct {
		subscriberID int
		value        int
	}
	deliveries := make(chan delivery, numSubscribers*numSets*2)

	// Start subscribers
	for i := 0; i < numSubscribers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var mu sync.Mutex
			received := []int{}
			unsub := s.Subscribe(func(v int) {
				mu.Lock()
				received = append(received, v)
				mu.Unlock()
				deliveries <- delivery{subscriberID: id, value: v}
			})
			defer unsub()
			// Keep subscriber alive during the test
			wg.Wait()
		}(i)
	}

	// Let subscribers register
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 1; j <= numSets; j++ {
			s.Set(j)
		}
	}()

	wg.Wait()
	close(deliveries)

	// Verify: for each subscriber, values should be in ascending order
	// (no newer value delivered before an older initial value)
	subscriberValues := make(map[int][]int)
	for d := range deliveries {
		subscriberValues[d.subscriberID] = append(subscriberValues[d.subscriberID], d.value)
	}

	for id, values := range subscriberValues {
		for i := 1; i < len(values); i++ {
			if values[i] < values[i-1] {
				t.Errorf("subscriber %d: value %d delivered before %d (out of order)", id, values[i], values[i-1])
			}
		}
	}
}

func TestComputedReentrancyNoDeadlock(t *testing.T) {
	s := NewSignal(1)
	var c *Computed[int]

	// Computed that calls itself via Get (reentrancy)
	c = NewComputed(func() int {
		val := s.Get()
		if val > 0 {
			// Reentrant call to Get
			_ = c.Get()
		}
		return val * 2
	})

	// Test Recompute doesn't deadlock with reentrancy
	done := make(chan bool, 1)
	go func() {
		s.Set(2)
		c.Recompute()
		done <- true
	}()

	select {
	case <-done:
		// Success, no deadlock
	case <-time.After(2 * time.Second):
		t.Fatal("Recompute deadlocked with reentrant Get call")
	}

	c.Dispose()
}
