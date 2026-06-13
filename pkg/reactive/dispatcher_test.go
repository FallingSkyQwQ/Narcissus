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
	RunOnUI(func() {
		executed = true
	})

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
	RunOnUI(func() {
		executed = true
	})

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
	RunOnUI(func() {
		executed = true
	})

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

func TestWindowsDispatcher(t *testing.T) {
	timeout := time.AfterFunc(5*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	// Create Windows dispatcher on current thread
	dispatcher, err := NewWindowsDispatcher()
	if err != nil {
		t.Fatalf("NewWindowsDispatcher failed: %v", err)
	}

	// Should be on UI thread (since we created it on this thread)
	if !dispatcher.IsUIThread() {
		t.Error("Should be on UI thread immediately after creation")
	}

	// RunOnUI should execute immediately when on UI thread
	var executed bool
	err = dispatcher.RunOnUI(func() {
		executed = true
	})
	if err != nil {
		t.Errorf("RunOnUI failed: %v", err)
	}

	if !executed {
		t.Error("RunOnUI should execute function immediately when on UI thread")
	}
}

func TestWindowsDispatcherNonUIThread(t *testing.T) {
	timeout := time.AfterFunc(5*time.Second, func() {
		t.Fatal("test timed out")
	})
	defer timeout.Stop()

	// Create dispatcher on current thread
	dispatcher, err := NewWindowsDispatcher()
	if err != nil {
		t.Fatalf("NewWindowsDispatcher failed: %v", err)
	}

	// Run on a different goroutine (different thread)
	done := make(chan error, 1)
	go func() {
		err := dispatcher.RunOnUI(func() {
			t.Error("Function should not execute on non-UI thread")
		})
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Error("RunOnUI should return error when called from non-UI thread")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("test timed out waiting for goroutine")
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

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numCalls; j++ {
				RunOnUI(func() {})
			}
		}()
	}

	wg.Wait()

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

	// Save original dispatcher
	original := globalDispatcher
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
				RunOnUI(func() {})
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
}
