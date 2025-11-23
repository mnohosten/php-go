package spl

import (
	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// SplStack - Stack (LIFO) implementation
// ============================================================================

// SplStack represents a stack data structure (Last In, First Out)
type SplStack struct {
	items []*types.Value
}

// NewSplStack creates a new stack
func NewSplStack() *SplStack {
	return &SplStack{
		items: make([]*types.Value, 0),
	}
}

// Push adds an element to the top of the stack
func (s *SplStack) Push(value *types.Value) {
	s.items = append(s.items, value)
}

// Pop removes and returns the element at the top of the stack
func (s *SplStack) Pop() (*types.Value, bool) {
	if len(s.items) == 0 {
		return types.NewNull(), false
	}

	index := len(s.items) - 1
	value := s.items[index]
	s.items = s.items[:index]
	return value, true
}

// Top returns the element at the top without removing it
func (s *SplStack) Top() (*types.Value, bool) {
	if len(s.items) == 0 {
		return types.NewNull(), false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty returns true if the stack is empty
func (s *SplStack) IsEmpty() bool {
	return len(s.items) == 0
}

// Count returns the number of elements in the stack
func (s *SplStack) Count() int {
	return len(s.items)
}

// ============================================================================
// SplQueue - Queue (FIFO) implementation
// ============================================================================

// SplQueue represents a queue data structure (First In, First Out)
type SplQueue struct {
	items []*types.Value
}

// NewSplQueue creates a new queue
func NewSplQueue() *SplQueue {
	return &SplQueue{
		items: make([]*types.Value, 0),
	}
}

// Enqueue adds an element to the end of the queue
func (q *SplQueue) Enqueue(value *types.Value) {
	q.items = append(q.items, value)
}

// Dequeue removes and returns the element at the front of the queue
func (q *SplQueue) Dequeue() (*types.Value, bool) {
	if len(q.items) == 0 {
		return types.NewNull(), false
	}

	value := q.items[0]
	q.items = q.items[1:]
	return value, true
}

// Peek returns the element at the front without removing it
func (q *SplQueue) Peek() (*types.Value, bool) {
	if len(q.items) == 0 {
		return types.NewNull(), false
	}
	return q.items[0], true
}

// IsEmpty returns true if the queue is empty
func (q *SplQueue) IsEmpty() bool {
	return len(q.items) == 0
}

// Count returns the number of elements in the queue
func (q *SplQueue) Count() int {
	return len(q.items)
}

// ============================================================================
// SplFixedArray - Fixed-size array implementation
// ============================================================================

// SplFixedArray represents a fixed-size array
type SplFixedArray struct {
	items []*types.Value
	size  int
}

// NewSplFixedArray creates a new fixed array with the given size
func NewSplFixedArray(size int) *SplFixedArray {
	items := make([]*types.Value, size)
	// Initialize with null values
	for i := 0; i < size; i++ {
		items[i] = types.NewNull()
	}
	return &SplFixedArray{
		items: items,
		size:  size,
	}
}

// Get returns the value at the given index
func (a *SplFixedArray) Get(index int) (*types.Value, bool) {
	if index < 0 || index >= a.size {
		return types.NewNull(), false
	}
	return a.items[index], true
}

// Set sets the value at the given index
func (a *SplFixedArray) Set(index int, value *types.Value) bool {
	if index < 0 || index >= a.size {
		return false
	}
	a.items[index] = value
	return true
}

// GetSize returns the size of the array
func (a *SplFixedArray) GetSize() int {
	return a.size
}

// SetSize resizes the array
func (a *SplFixedArray) SetSize(newSize int) {
	if newSize < 0 {
		return
	}

	newItems := make([]*types.Value, newSize)
	copyLen := newSize
	if a.size < newSize {
		copyLen = a.size
	}

	// Copy existing items
	copy(newItems, a.items[:copyLen])

	// Initialize new slots with null
	for i := copyLen; i < newSize; i++ {
		newItems[i] = types.NewNull()
	}

	a.items = newItems
	a.size = newSize
}

// ToArray converts the fixed array to a regular array
func (a *SplFixedArray) ToArray() *types.Array {
	arr := types.NewEmptyArray()
	for i, item := range a.items {
		arr.Set(types.NewInt(int64(i)), item)
	}
	return arr
}

// ============================================================================
// SplDoublyLinkedList - Doubly linked list implementation
// ============================================================================

// SplDoublyLinkedListNode represents a node in the doubly linked list
type SplDoublyLinkedListNode struct {
	value *types.Value
	prev  *SplDoublyLinkedListNode
	next  *SplDoublyLinkedListNode
}

// SplDoublyLinkedList represents a doubly linked list
type SplDoublyLinkedList struct {
	head  *SplDoublyLinkedListNode
	tail  *SplDoublyLinkedListNode
	count int
}

// NewSplDoublyLinkedList creates a new doubly linked list
func NewSplDoublyLinkedList() *SplDoublyLinkedList {
	return &SplDoublyLinkedList{
		head:  nil,
		tail:  nil,
		count: 0,
	}
}

// Push adds an element to the end of the list
func (l *SplDoublyLinkedList) Push(value *types.Value) {
	node := &SplDoublyLinkedListNode{
		value: value,
		prev:  l.tail,
		next:  nil,
	}

	if l.tail != nil {
		l.tail.next = node
	}

	l.tail = node

	if l.head == nil {
		l.head = node
	}

	l.count++
}

// Pop removes and returns the element at the end of the list
func (l *SplDoublyLinkedList) Pop() (*types.Value, bool) {
	if l.tail == nil {
		return types.NewNull(), false
	}

	value := l.tail.value
	l.tail = l.tail.prev

	if l.tail != nil {
		l.tail.next = nil
	} else {
		l.head = nil
	}

	l.count--
	return value, true
}

// Unshift adds an element to the beginning of the list
func (l *SplDoublyLinkedList) Unshift(value *types.Value) {
	node := &SplDoublyLinkedListNode{
		value: value,
		prev:  nil,
		next:  l.head,
	}

	if l.head != nil {
		l.head.prev = node
	}

	l.head = node

	if l.tail == nil {
		l.tail = node
	}

	l.count++
}

// Shift removes and returns the element at the beginning of the list
func (l *SplDoublyLinkedList) Shift() (*types.Value, bool) {
	if l.head == nil {
		return types.NewNull(), false
	}

	value := l.head.value
	l.head = l.head.next

	if l.head != nil {
		l.head.prev = nil
	} else {
		l.tail = nil
	}

	l.count--
	return value, true
}

// Top returns the element at the end without removing it
func (l *SplDoublyLinkedList) Top() (*types.Value, bool) {
	if l.tail == nil {
		return types.NewNull(), false
	}
	return l.tail.value, true
}

// Bottom returns the element at the beginning without removing it
func (l *SplDoublyLinkedList) Bottom() (*types.Value, bool) {
	if l.head == nil {
		return types.NewNull(), false
	}
	return l.head.value, true
}

// IsEmpty returns true if the list is empty
func (l *SplDoublyLinkedList) IsEmpty() bool {
	return l.count == 0
}

// Count returns the number of elements in the list
func (l *SplDoublyLinkedList) Count() int {
	return l.count
}

// ============================================================================
// SplHeap - Abstract heap implementation
// ============================================================================

// SplHeap represents a heap data structure
type SplHeap struct {
	items   []*types.Value
	compare func(a, b *types.Value) int
}

// NewSplHeap creates a new heap with a custom compare function
func NewSplHeap(compare func(a, b *types.Value) int) *SplHeap {
	return &SplHeap{
		items:   make([]*types.Value, 0),
		compare: compare,
	}
}

// Insert adds an element to the heap
func (h *SplHeap) Insert(value *types.Value) {
	h.items = append(h.items, value)
	h.heapifyUp(len(h.items) - 1)
}

// Extract removes and returns the top element
func (h *SplHeap) Extract() (*types.Value, bool) {
	if len(h.items) == 0 {
		return types.NewNull(), false
	}

	value := h.items[0]
	lastIdx := len(h.items) - 1

	h.items[0] = h.items[lastIdx]
	h.items = h.items[:lastIdx]

	if len(h.items) > 0 {
		h.heapifyDown(0)
	}

	return value, true
}

// Top returns the top element without removing it
func (h *SplHeap) Top() (*types.Value, bool) {
	if len(h.items) == 0 {
		return types.NewNull(), false
	}
	return h.items[0], true
}

// IsEmpty returns true if the heap is empty
func (h *SplHeap) IsEmpty() bool {
	return len(h.items) == 0
}

// Count returns the number of elements in the heap
func (h *SplHeap) Count() int {
	return len(h.items)
}

// heapifyUp maintains heap property from bottom to top
func (h *SplHeap) heapifyUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2

		if h.compare(h.items[index], h.items[parent]) <= 0 {
			break
		}

		h.items[index], h.items[parent] = h.items[parent], h.items[index]
		index = parent
	}
}

// heapifyDown maintains heap property from top to bottom
func (h *SplHeap) heapifyDown(index int) {
	size := len(h.items)

	for {
		largest := index
		left := 2*index + 1
		right := 2*index + 2

		if left < size && h.compare(h.items[left], h.items[largest]) > 0 {
			largest = left
		}

		if right < size && h.compare(h.items[right], h.items[largest]) > 0 {
			largest = right
		}

		if largest == index {
			break
		}

		h.items[index], h.items[largest] = h.items[largest], h.items[index]
		index = largest
	}
}

// ============================================================================
// SplMaxHeap - Max heap (largest element at top)
// ============================================================================

// NewSplMaxHeap creates a new max heap
func NewSplMaxHeap() *SplHeap {
	return NewSplHeap(func(a, b *types.Value) int {
		aNum := a.ToFloat()
		bNum := b.ToFloat()

		if aNum > bNum {
			return 1
		} else if aNum < bNum {
			return -1
		}
		return 0
	})
}

// ============================================================================
// SplMinHeap - Min heap (smallest element at top)
// ============================================================================

// NewSplMinHeap creates a new min heap
func NewSplMinHeap() *SplHeap {
	return NewSplHeap(func(a, b *types.Value) int {
		aNum := a.ToFloat()
		bNum := b.ToFloat()

		// Reverse comparison for min heap
		if aNum < bNum {
			return 1
		} else if aNum > bNum {
			return -1
		}
		return 0
	})
}
