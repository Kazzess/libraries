package safe

import "sync"

// WaitGroup wraps sync.WaitGroup for safe concurrent usage.
type WaitGroup struct {
	wg *sync.WaitGroup
	mu sync.RWMutex
}

// NewWaitGroup creates a new *WaitGroup, wg maybe nil.
func NewWaitGroup(wg *sync.WaitGroup) *WaitGroup {
	w := &WaitGroup{
		wg: wg,
		mu: sync.RWMutex{},
	}

	if w.wg == nil {
		w.wg = new(sync.WaitGroup)
	}

	return w
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// This adds locks to call (sync.WaitGroup).Add.
func (wg *WaitGroup) Add(delta int) {
	wg.mu.RLock()
	defer wg.mu.RUnlock()

	wg.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
// This adds locks to call (sync.WaitGroup).Wait.
func (wg *WaitGroup) Wait() {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	wg.wg.Wait()
}
