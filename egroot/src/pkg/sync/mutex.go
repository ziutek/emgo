package sync

// A Mutex is a mutual exclusion lock.
// The zero value for a Mutex is an unlocked mutex.
type Mutex mutex

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {
	m.lock()
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
func (m *Mutex) Unlock() {
	m.unlock()
}