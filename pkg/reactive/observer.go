package reactive

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// source is a reactive value whose changes invalidate the dependents that read
// it. Both Signal and Computed implement it, which is what lets computeds be
// chained (a computed observes its sources and is itself a source).
type source interface {
	addDependent(d *dependent)
	removeDependent(d *dependent)
}

// dependent is a reactive computation: a function or expression that reads one
// or more sources and must be re-evaluated when any of them change.
//
// A dependent owns two sets:
//   - sources: the values it currently reads (rebuilt on every evaluation)
//   - dependents: the downstream computations that read *this* dependent
//
// Invalidations propagate downstream through dependents, so marking a signal
// dirty transitively marks every computed and effect that depends on it.
type dependent struct {
	mu         sync.Mutex
	sources    map[source]struct{}
	dependents map[*dependent]struct{}
	invalid    bool
	// flush runs the computation eagerly when it is invalidated. Effects
	// always set it; computeds set it only while they have subscribers, so an
	// unobserved computed stays lazy.
	flush func()
}

func newDependent() *dependent {
	return &dependent{
		sources:    make(map[source]struct{}),
		dependents: make(map[*dependent]struct{}),
	}
}

// beginTracking forgets the sources collected by the previous evaluation so a
// new evaluation builds a fresh dependency set.
func (d *dependent) beginTracking() {
	d.mu.Lock()
	sources := make([]source, 0, len(d.sources))
	for s := range d.sources {
		sources = append(sources, s)
	}
	d.sources = make(map[source]struct{})
	d.mu.Unlock()

	for _, s := range sources {
		s.removeDependent(d)
	}
}

// trackSource records that the current evaluation read s.
func (d *dependent) trackSource(s source) {
	d.mu.Lock()
	if _, ok := d.sources[s]; ok {
		d.mu.Unlock()
		return
	}
	d.sources[s] = struct{}{}
	d.mu.Unlock()

	s.addDependent(d)
}

// addDependent / removeDependent implement source, making a computed usable as
// a source for computations that read it.
func (d *dependent) addDependent(child *dependent) {
	d.mu.Lock()
	d.dependents[child] = struct{}{}
	d.mu.Unlock()
}

func (d *dependent) removeDependent(child *dependent) {
	d.mu.Lock()
	delete(d.dependents, child)
	d.mu.Unlock()
}

// invalidate marks d stale and propagates the change to everything downstream.
// The first invalidation wins, which keeps a diamond from being recomputed
// twice and stops cycles from looping forever.
func (d *dependent) invalidate() {
	d.mu.Lock()
	if d.invalid {
		d.mu.Unlock()
		return
	}
	d.invalid = true
	flush := d.flush
	dependents := make([]*dependent, 0, len(d.dependents))
	for child := range d.dependents {
		dependents = append(dependents, child)
	}
	d.mu.Unlock()

	if flush != nil {
		schedule(d)
	}
	for _, child := range dependents {
		child.invalidate()
	}
}

func (d *dependent) flushFn() func() {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.flush
}

func (d *dependent) clearInvalid() {
	d.mu.Lock()
	d.invalid = false
	d.mu.Unlock()
}

func (d *dependent) dispose() {
	d.beginTracking()

	d.mu.Lock()
	dependents := make([]*dependent, 0, len(d.dependents))
	for child := range d.dependents {
		dependents = append(dependents, child)
	}
	d.dependents = make(map[*dependent]struct{})
	d.flush = nil
	d.mu.Unlock()

	for _, child := range dependents {
		child.detachSource(d)
	}
}

// detachSource drops s from d's dependency set without touching d's other
// sources. It is used when a source is disposed.
func (d *dependent) detachSource(s source) {
	d.mu.Lock()
	delete(d.sources, s)
	d.mu.Unlock()
}

// --- Dependency collection context ---------------------------------------
//
// Implicit tracking needs to know which dependent is currently evaluating.
// The context is kept per goroutine so nested evaluations (a computed reading
// another computed) and concurrent evaluations do not clobber each other.

var (
	trackingMu     sync.Mutex
	trackingStacks = make(map[uint64][]*dependent)
)

func currentDependent() *dependent {
	id := goroutineID()
	trackingMu.Lock()
	defer trackingMu.Unlock()
	stack := trackingStacks[id]
	if len(stack) == 0 {
		return nil
	}
	return stack[len(stack)-1]
}

// track evaluates fn with d as the active dependent, collecting every source
// fn reads.
func track(d *dependent, fn func()) {
	id := goroutineID()
	d.beginTracking()

	trackingMu.Lock()
	trackingStacks[id] = append(trackingStacks[id], d)
	trackingMu.Unlock()

	defer func() {
		trackingMu.Lock()
		stack := trackingStacks[id]
		if len(stack) > 0 {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			delete(trackingStacks, id)
		} else {
			trackingStacks[id] = stack
		}
		trackingMu.Unlock()
	}()

	fn()
}

// --- Batched scheduling ---------------------------------------------------
//
// Invalidating an effect schedules it instead of running it inline, so a burst
// of Set calls runs the effect once (Batch makes that grouping explicit).

var (
	schedMu    sync.Mutex
	batchDepth int
	draining   bool
	queue      []*dependent
	queued     = make(map[*dependent]struct{})
)

// Batch groups several updates so effects and subscribed computeds run once at
// the end instead of after every change. Calls nest.
func Batch(fn func()) {
	if fn == nil {
		return
	}
	beginMutation()
	defer endMutation()
	fn()
}

// beginMutation / endMutation bracket an update so the eager queue drains once
// everything has been marked stale. Signal.Set uses them so that a single Set
// cannot run an effect while some of the signal's other dependents are still
// holding their previous values.
func beginMutation() {
	schedMu.Lock()
	batchDepth++
	schedMu.Unlock()
}

func endMutation() {
	schedMu.Lock()
	batchDepth--
	start := batchDepth == 0
	schedMu.Unlock()
	if start {
		flushQueue()
	}
}

func schedule(d *dependent) {
	schedMu.Lock()
	if _, ok := queued[d]; !ok {
		queued[d] = struct{}{}
		queue = append(queue, d)
	}
	schedMu.Unlock()

	flushQueue()
}

func flushQueue() {
	schedMu.Lock()
	if batchDepth > 0 || draining {
		schedMu.Unlock()
		return
	}
	draining = true
	schedMu.Unlock()

	for {
		schedMu.Lock()
		if len(queue) == 0 {
			draining = false
			schedMu.Unlock()
			return
		}
		batch := queue
		queue = nil
		queued = make(map[*dependent]struct{})
		schedMu.Unlock()

		for _, d := range batch {
			if fn := d.flushFn(); fn != nil {
				fn()
			}
		}
	}
}

// goroutineID extracts the current goroutine's id from the runtime stack.
func goroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	const prefix = "goroutine "
	stack := string(buf[:n])
	if !strings.HasPrefix(stack, prefix) {
		return 0
	}
	rest := stack[len(prefix):]
	end := strings.IndexByte(rest, ' ')
	if end < 0 {
		end = len(rest)
	}
	id, _ := strconv.ParseUint(rest[:end], 10, 64)
	return id
}
