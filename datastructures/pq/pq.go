package pq

import (
	"container/heap"
)

/*
	PQItem defines a node for a priority queue (max-heap)
*/
type PQItem struct {
	Value    interface{}
	Priority int64
	Index    int
}

/*
	Priority queue is based on a regular array
*/
type PriorityQueue []*PQItem

/*
	Len returns the length of a priority queue
*/
func (pq PriorityQueue) Len() int {
	return len(pq)
}

/*
	Less defines the order between two priority queue entries
*/
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

/*
	Swap is used for swapping two priority queue elements on sorting
*/
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

/*
	Push inserts an element to a priority queue. Complexity: O(logN)
*/
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PQItem)
	item.Index = n
	*pq = append(*pq, item)
}

/*
	Pop pulls the top (max priority) element from a priority queue. Complexity: O(1)
*/
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

/*
	Update allows to change a priority queue item's value and priority. Complexity: O(logN)
*/
func (pq *PriorityQueue) update(item *PQItem, value interface{}, priority int64) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}
