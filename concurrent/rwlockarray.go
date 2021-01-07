package concurrent

import "sync"

type RWLockArray struct {
	Locks []*sync.RWMutex
	Size  int
}

func NewRWLockArray(size int) *RWLockArray {
	a := new(RWLockArray)
	a.Locks = make([]*sync.RWMutex, size)
	for i := 0; i < size; i++ {
		a.Locks[i] = new(sync.RWMutex)
	}
	a.Size = size
	return a
}

func (a *RWLockArray) Get(i int) *sync.RWMutex {
	return a.Locks[i]
}

func (a *RWLockArray) LockAll() {
	for i := 0; i < a.Size; i++ {
		a.Locks[i].Lock()
	}
}

func (a *RWLockArray) UnlockAll() {
	for i := 0; i < a.Size; i++ {
		a.Locks[i].Unlock()
	}
}

func (a *RWLockArray) RLockAll() {
	for i := 0; i < a.Size; i++ {
		a.Locks[i].RLock()
	}
}

func (a *RWLockArray) RUnlockAll() {
	for i := 0; i < a.Size; i++ {
		a.Locks[i].RUnlock()
	}
}
