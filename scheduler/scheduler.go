package scheduler

import (
	"container/heap"
	"context"
	"github.com/rovergulf/utils/datastructures/minheap"
	"github.com/rovergulf/utils/storages"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Scheduler struct {
	dump      *storages.Dump
	eventHeap minheap.MinHeap
	lock      *sync.RWMutex
	onHandle  func(context.Context, []interface{})
	logger    *zap.SugaredLogger
}

func NewScheduler(lg *zap.SugaredLogger, period time.Duration, quit chan struct{}, onHandle func(context.Context, []interface{}), dumpFileName string, flushDelay time.Duration) *Scheduler {
	s := new(Scheduler)

	s.eventHeap = make(minheap.MinHeap, 0)
	s.logger = lg.Named("scheduler")
	heap.Init(&s.eventHeap)
	s.lock = new(sync.RWMutex)
	s.onHandle = onHandle
	s.dump = storages.NewDump(s.logger, dumpFileName, flushDelay, s.Flush, s.OnFlushComplete)

	var version int64
	var recoverHeap minheap.MinHeap
	err := s.dump.Recover(&version, &recoverHeap)
	if err != nil {
		s.logger.Errorf("Unable to recover from dump: %s", err)
	} else {
		if recoverHeap != nil {
			s.eventHeap = recoverHeap
			heap.Init(&s.eventHeap)
		}
		s.logger.Infof("Recovered %d objects", s.eventHeap.Len())
	}

	s.dump.StartFlushThread()

	s.logger.Infof("Starting scheduler with %d period", period)
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
					s.logger.Infof("Current event heap length: %d", heapLen)
				}

				var toProcess []interface{}
				for s.eventHeap.Len() > 0 {
					lastElem := s.eventHeap.Peek().(*minheap.PQItem)
					s.logger.Infof("Top head element sent time: %d; current time: %d", lastElem.Priority, currentTime)

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
					s.logger.Infof("Scheduler processed this run: %d", count)
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
