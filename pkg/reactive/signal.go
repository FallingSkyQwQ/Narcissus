package reactive

import (
	"sync"
	"sync/atomic"
)

// Signal represents a reactive state container
type Signal[T any] struct {
	mu        sync.RWMutex
	value     T
	version   uint64
	observers map[uint64]func(T)
	nextID    atomic.Uint64
}

// NewSignal creates a new Signal with an initial value
func NewSignal[T any](initial T) *Signal[T] {
	return &Signal[T]{
		value:     initial,
		observers: make(map[uint64]func(T)),
	}
}

// Get returns the current value
func (s *Signal[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set updates the value and notifies observers
func (s *Signal[T]) Set(newValue T) {
	s.mu.Lock()
	s.value = newValue
	s.version++
	observers := make([]func(T), 0, len(s.observers))
	for _, observer := range s.observers {
		observers = append(observers, observer)
	}
	s.mu.Unlock()

	// Notify outside of lock to prevent deadlocks
	for _, observer := range observers {
		observer(newValue)
	}
}

// Subscribe registers an observer and returns an unsubscribe function.
// The observer is called immediately with the current value, and then on each Set.
// The initial callback is skipped if a newer Set has already occurred by the time
// the observer is registered (to prevent out-of-order delivery).
func (s *Signal[T]) Subscribe(observer func(T)) func() {
	// Generate unique ID using atomic increment (returns new value, so subtract 1 for 0-based index)
	id := s.nextID.Add(1) - 1

	s.mu.Lock()
	s.observers[id] = observer
	currentValue := s.value
	currentVersion := s.version
	s.mu.Unlock()

	// Call observer immediately with current value only if version hasn't changed
	// This prevents out-of-order delivery if Set occurs between registration and initial callback
	s.mu.RLock()
	stillCurrent := s.version == currentVersion
	s.mu.RUnlock()

	if stillCurrent {
		observer(currentValue)
	}

	return func() {
		s.mu.Lock()
		delete(s.observers, id)
		s.mu.Unlock()
	}
}
