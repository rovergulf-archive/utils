package minheap

import "container/heap"

/*
	PQItem defines a node for a priority queue
*/
type PQItem struct {
	Value    interface{}
	Priority int64
	Index    int
}

/*
	Priority queue is based on a regular array
*/
type MinHeap []*PQItem

/*
	Len returns the length of a priority queue
*/
func (h MinHeap) Len() int {
	return len(h)
}

/*
	Less defines the order between two priority queue entries
*/
func (h MinHeap) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}

/*
	Swap is used for swapping two priority queue elements on sorting
*/
func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

/*
	Push inserts an element to a priority queue. Complexity: O(logN)
*/
func (h *MinHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*PQItem)
	item.Index = n
	*h = append(*h, item)
}

/*
	Peek returns the top element (min priority)
*/
func (h *MinHeap) Peek() interface{} {
	arr := *h
	n := len(arr)
	if n == 0 {
		return nil
	}
	return arr[0]
}

/*
	Pop pulls the top (max priority) element from a priority queue. Complexity: O(1)
*/
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*h = old[0 : n-1]
	return item
}

/*
	Update allows to change a priority queue item's value and priority. Complexity: O(logN)
*/
func (h *MinHeap) Update(item *PQItem, value interface{}, priority int64) {
	item.Value = value
	item.Priority = priority
	heap.Fix(h, item.Index)
}
