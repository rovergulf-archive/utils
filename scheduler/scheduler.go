package scheduler

import (
	"container/heap"
	"context"
	"github.com/rovergulf/utils/clog"
	"github.com/rovergulf/utils/datastructures/minheap"
	"github.com/rovergulf/utils/storages"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	dump      *storages.Dump
	eventHeap minheap.MinHeap
	lock      *sync.RWMutex
	onHandle  func(context.Context, []interface{})
}

func NewScheduler(period time.Duration, quit chan struct{}, onHandle func(context.Context, []interface{}), dumpFileName string, flushDelay time.Duration) *Scheduler {
	s := new(Scheduler)

	s.eventHeap = make(minheap.MinHeap, 0)
	heap.Init(&s.eventHeap)
	s.lock = new(sync.RWMutex)
	s.onHandle = onHandle
	s.dump = storages.NewDump(dumpFileName, flushDelay, s.Flush, s.OnFlushComplete)

	var version int64
	var recoverHeap minheap.MinHeap
	err := s.dump.Recover(&version, &recoverHeap)
	if err != nil {
		clog.Errorf("Unable to recover from dump: %s", err)
	} else {
		if recoverHeap != nil {
			s.eventHeap = recoverHeap
			heap.Init(&s.eventHeap)
		}
		log.Printf("Recovered %d objects", s.eventHeap.Len())
	}

	s.dump.StartFlushThread()

	log.Printf("Starting scheduler with %d period", period)
	ticker := time.NewTicker(period)
	go func() {
		for {
			select {
			case <-ticker.C:
				currentTime := time.Now().Unix()
				count := 0
				s.lock.Lock()

				heapLen := s.eventHeap.Len()
				if heapLen > 0 {
					log.Printf("Current event heap length: %d", heapLen)
				}

				var toProcess []interface{}
				for s.eventHeap.Len() > 0 {
					lastElem := s.eventHeap.Peek().(*minheap.PQItem)
					log.Printf("Top head element sent time: %d; current time: %d", lastElem.Priority, currentTime)

					if lastElem.Priority > currentTime {
						break
					}

					heap.Pop(&s.eventHeap)
					toProcess = append(toProcess, lastElem)
					count += 1
				}

				s.lock.Unlock()

				if toProcess != nil {
					onHandle(context.Background(), toProcess)
					toProcess = nil // wow, probably it's even not a kludge
				}

				if count > 0 {
					log.Printf("Scheduler processed this run: %d", count)
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return s
}

func (s *Scheduler) GetDump() *storages.Dump {
	return s.dump
}

func (s *Scheduler) NewEvent(executionTime int64, event interface{}) {

	if executionTime < time.Now().Unix() {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	heap.Push(&s.eventHeap, &minheap.PQItem{
		Value:    event,
		Priority: executionTime,
	})
}

// ???
func (s *Scheduler) RemoveEvent(executionTime int, event interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	heap.Remove(&s.eventHeap, executionTime)
}

// ???
func (s *Scheduler) UpdateEvent(originalEventTime int64, event interface{}, executionTime int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.eventHeap.Update(&minheap.PQItem{
		Value:    event,
		Priority: originalEventTime,
	}, event, executionTime)
}

func (s *Scheduler) Flush() {
	s.lock.Lock()
	s.dump.Flush(time.Now().Unix(), s.eventHeap)
	defer s.lock.Unlock()
}

func (s *Scheduler) OnFlushComplete() {}
