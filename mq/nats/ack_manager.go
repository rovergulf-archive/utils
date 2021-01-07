package natsmq

import (
	"github.com/rovergulf/utils/storages"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type AckManager struct {
	Dump     *storages.Dump
	Version  int64
	Lock     *sync.RWMutex
	Rand     *rand.Rand
	sequence uint64
}

type ackTimestampState struct {
	Sequence uint64
}

func NewAckTimestampManager(dumpFileName string, flushDelay time.Duration) *AckManager {
	c := new(AckManager)
	c.Lock = new(sync.RWMutex)
	c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	c.Dump = storages.NewDump(dumpFileName, flushDelay, c.Flush, c.OnFlushComplete)

	var state ackTimestampState
	if err := c.Dump.Recover(&c.Version, &state); err != nil {
		log.Printf("Failed to recover ack timestamp dump: %s", err)
	}

	c.Lock.Lock()

	c.sequence = state.Sequence

	c.Lock.Unlock()

	c.Dump.StartFlushThread()

	return c
}

func (a *AckManager) Flush() {
	a.Lock.Lock()

	a.Dump.Flush(a.Version, ackTimestampState{
		Sequence: a.sequence,
	})

	a.Lock.Unlock()
}

func (a *AckManager) OnFlushComplete() {
}

func (a *AckManager) Set(sequence uint64) {
	atomic.AddInt64(&a.Version, 1)

	a.Lock.Lock()

	a.sequence = sequence + 1

	a.Lock.Unlock()
}

func (a *AckManager) Get() uint64 {
	var sequence uint64

	a.Lock.RLock()

	sequence = a.sequence

	a.Lock.RUnlock()

	return sequence
}
