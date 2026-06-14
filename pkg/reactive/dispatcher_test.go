package reactive

import (
	"sync"
	"testing"
	"time"
)

// mockDispatcher is a mock implementation of Dispatcher for testing
type mockDispatcher struct {
	isUIThread bool
	scheduled  []func()
	mu         sync.Mutex
}

func (m *mockDispatcher) RunOnUI(fn func()) error {
	m.mu.Lock()
	m.scheduled = append(m.scheduled, fn)
	m.mu.Unlock()
	return nil
}

func (m *mockDispatcher) IsUIThread() bool {
	return m.isUIThread
}

func (m *mockDispatcher) getScheduledCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.scheduled)
}

func (m *mockDispatcher) executeScheduled() {
	m.mu.Lock()
	funcs := make([]func(), len(m.scheduled))
	copy(funcs, m.scheduled)
	m.scheduled = m.scheduled[:0]
	m.mu.Unlock()

	for _, fn := range funcs {
		fn()
	}
}

func TestDispatcher(t *testing.T) {
	timeout := time.AfterFunc(5*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	mock := &mockDispatcher{isUIThread: false}
	SetDispatcher(mock)

	var executed bool
	err := RunOnUI(func() {
		executed = true
	})
	if err != nil {
		t.Errorf("RunOnUI failed: %v", err)
	}

	if mock.getScheduledCount() != 1 {
		t.Errorf("expected 1 scheduled function, got %d", mock.getScheduledCount())
	}

	// Execute the scheduled function
	mock.executeScheduled()

	if !executed {
		t.Error("function was not executed")
	}
}

func TestDispatcherIsUIThread(t *testing.T) {
	timeout := time.AfterFunc(5*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	// Test with UI thread
	mock := &mockDispatcher{isUIThread: true}
	SetDispatcher(mock)

	var executed bool
	err := RunOnUI(func() {
		executed = true
	})
	if err != nil {
		t.Errorf("RunOnUI failed: %v", err)
	}

	// When on UI thread, function should execute immediately
	if mock.getScheduledCount() != 0 {
		t.Errorf("expected 0 scheduled functions when on UI thread, got %d", mock.getScheduledCount())
	}

	if !executed {
		t.Error("function should execute immediately on UI thread")
	}
}

func TestDispatcherNil(t *testing.T) {
	timeout := time.AfterFunc(5*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	// Test with nil dispatcher
	SetDispatcher(nil)

	var executed bool
	err := RunOnUI(func() {
		executed = true
	})
	if err != nil {
		t.Errorf("RunOnUI failed: %v", err)
	}

	// When dispatcher is nil, function should execute immediately
	if !executed {
		t.Error("function should execute immediately when dispatcher is nil")
	}
}

func TestIsUIThread(t *testing.T) {
	timeout := time.AfterFunc(5*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	// Test with mock dispatcher
	mock := &mockDispatcher{isUIThread: true}
	SetDispatcher(mock)

	if !IsUIThread() {
		t.Error("IsUIThread should return true when mock returns true")
	}

	mock.isUIThread = false
	if IsUIThread() {
		t.Error("IsUIThread should return false when mock returns false")
	}

	// Test with nil dispatcher
	SetDispatcher(nil)
	if !IsUIThread() {
		t.Error("IsUIThread should return true when dispatcher is nil")
	}
}

func TestDispatcherConcurrent(t *testing.T) {
	timeout := time.AfterFunc(10*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	mock := &mockDispatcher{isUIThread: false}
	SetDispatcher(mock)

	var wg sync.WaitGroup
	numGoroutines := 10
	numCalls := 100
	errorsCh := make(chan error, numGoroutines*numCalls)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numCalls; j++ {
				if err := RunOnUI(func() {}); err != nil {
					errorsCh <- err
				}
			}
		}()
	}

	wg.Wait()
	close(errorsCh)

	for err := range errorsCh {
		t.Errorf("RunOnUI failed: %v", err)
	}

	expectedCount := numGoroutines * numCalls
	if mock.getScheduledCount() != expectedCount {
		t.Errorf("expected %d scheduled functions, got %d", expectedCount, mock.getScheduledCount())
	}
}

func TestSetDispatcher(t *testing.T) {
	timeout := time.AfterFunc(5*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	// Save original dispatcher with lock
	globalDispatcherMu.RLock()
	original := globalDispatcher
	globalDispatcherMu.RUnlock()
	defer SetDispatcher(original)

	// Test setting a new dispatcher
	mock := &mockDispatcher{isUIThread: false}
	SetDispatcher(mock)

	globalDispatcherMu.RLock()
	current := globalDispatcher
	globalDispatcherMu.RUnlock()

	if current != mock {
		t.Error("SetDispatcher should set the global dispatcher")
	}
}

func TestDispatcherConcurrentSetAndGet(t *testing.T) {
	timeout := time.AfterFunc(10*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	var wg sync.WaitGroup
	numGoroutines := 10
	numIterations := 100
	errorsCh := make(chan error, numGoroutines*numIterations)

	// Concurrent SetDispatcher
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				mock := &mockDispatcher{isUIThread: id%2 == 0}
				SetDispatcher(mock)
			}
		}(i)
	}

	// Concurrent RunOnUI
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				if err := RunOnUI(func() {}); err != nil {
					errorsCh <- err
				}
			}
		}()
	}

	// Concurrent IsUIThread
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				_ = IsUIThread()
			}
		}()
	}

	wg.Wait()
	close(errorsCh)

	for err := range errorsCh {
		t.Errorf("RunOnUI failed: %v", err)
	}
}
