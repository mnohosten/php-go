package spl

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// SplStack Tests
// ============================================================================

func TestSplStackPushPop(t *testing.T) {
	stack := NewSplStack()

	// Test empty stack
	if !stack.IsEmpty() {
		t.Error("New stack should be empty")
	}

	if stack.Count() != 0 {
		t.Errorf("Expected count 0, got %d", stack.Count())
	}

	// Push elements
	stack.Push(types.NewInt(1))
	stack.Push(types.NewInt(2))
	stack.Push(types.NewInt(3))

	if stack.Count() != 3 {
		t.Errorf("Expected count 3, got %d", stack.Count())
	}

	// Pop elements (LIFO)
	value, ok := stack.Pop()
	if !ok || value.ToInt() != 3 {
		t.Errorf("Expected 3, got %d", value.ToInt())
	}

	value, ok = stack.Pop()
	if !ok || value.ToInt() != 2 {
		t.Errorf("Expected 2, got %d", value.ToInt())
	}

	value, ok = stack.Pop()
	if !ok || value.ToInt() != 1 {
		t.Errorf("Expected 1, got %d", value.ToInt())
	}

	// Pop from empty stack
	_, ok = stack.Pop()
	if ok {
		t.Error("Expected false when popping from empty stack")
	}
}

func TestSplStackTop(t *testing.T) {
	stack := NewSplStack()

	// Top on empty stack
	_, ok := stack.Top()
	if ok {
		t.Error("Expected false when getting top of empty stack")
	}

	stack.Push(types.NewInt(10))
	stack.Push(types.NewInt(20))

	// Top should not remove element
	value, ok := stack.Top()
	if !ok || value.ToInt() != 20 {
		t.Errorf("Expected 20, got %d", value.ToInt())
	}

	if stack.Count() != 2 {
		t.Error("Top should not remove elements")
	}
}

// ============================================================================
// SplQueue Tests
// ============================================================================

func TestSplQueueEnqueueDequeue(t *testing.T) {
	queue := NewSplQueue()

	// Test empty queue
	if !queue.IsEmpty() {
		t.Error("New queue should be empty")
	}

	// Enqueue elements
	queue.Enqueue(types.NewInt(1))
	queue.Enqueue(types.NewInt(2))
	queue.Enqueue(types.NewInt(3))

	if queue.Count() != 3 {
		t.Errorf("Expected count 3, got %d", queue.Count())
	}

	// Dequeue elements (FIFO)
	value, ok := queue.Dequeue()
	if !ok || value.ToInt() != 1 {
		t.Errorf("Expected 1, got %d", value.ToInt())
	}

	value, ok = queue.Dequeue()
	if !ok || value.ToInt() != 2 {
		t.Errorf("Expected 2, got %d", value.ToInt())
	}

	value, ok = queue.Dequeue()
	if !ok || value.ToInt() != 3 {
		t.Errorf("Expected 3, got %d", value.ToInt())
	}

	// Dequeue from empty queue
	_, ok = queue.Dequeue()
	if ok {
		t.Error("Expected false when dequeuing from empty queue")
	}
}

func TestSplQueuePeek(t *testing.T) {
	queue := NewSplQueue()

	// Peek on empty queue
	_, ok := queue.Peek()
	if ok {
		t.Error("Expected false when peeking empty queue")
	}

	queue.Enqueue(types.NewInt(10))
	queue.Enqueue(types.NewInt(20))

	// Peek should not remove element
	value, ok := queue.Peek()
	if !ok || value.ToInt() != 10 {
		t.Errorf("Expected 10, got %d", value.ToInt())
	}

	if queue.Count() != 2 {
		t.Error("Peek should not remove elements")
	}
}

// ============================================================================
// SplFixedArray Tests
// ============================================================================

func TestSplFixedArrayBasic(t *testing.T) {
	arr := NewSplFixedArray(5)

	if arr.GetSize() != 5 {
		t.Errorf("Expected size 5, got %d", arr.GetSize())
	}

	// Set values
	arr.Set(0, types.NewInt(10))
	arr.Set(2, types.NewInt(20))
	arr.Set(4, types.NewInt(30))

	// Get values
	value, ok := arr.Get(0)
	if !ok || value.ToInt() != 10 {
		t.Errorf("Expected 10, got %d", value.ToInt())
	}

	value, ok = arr.Get(2)
	if !ok || value.ToInt() != 20 {
		t.Errorf("Expected 20, got %d", value.ToInt())
	}

	// Unset index should be null
	value, ok = arr.Get(1)
	if !ok || value.Type() != types.TypeNull {
		t.Error("Expected null for unset index")
	}
}

func TestSplFixedArrayBounds(t *testing.T) {
	arr := NewSplFixedArray(3)

	// Out of bounds get
	_, ok := arr.Get(-1)
	if ok {
		t.Error("Expected false for negative index")
	}

	_, ok = arr.Get(10)
	if ok {
		t.Error("Expected false for out of bounds index")
	}

	// Out of bounds set
	ok = arr.Set(-1, types.NewInt(10))
	if ok {
		t.Error("Expected false when setting negative index")
	}

	ok = arr.Set(10, types.NewInt(10))
	if ok {
		t.Error("Expected false when setting out of bounds index")
	}
}

func TestSplFixedArrayResize(t *testing.T) {
	arr := NewSplFixedArray(3)
	arr.Set(0, types.NewInt(10))
	arr.Set(1, types.NewInt(20))
	arr.Set(2, types.NewInt(30))

	// Resize larger
	arr.SetSize(5)
	if arr.GetSize() != 5 {
		t.Errorf("Expected size 5, got %d", arr.GetSize())
	}

	// Old values should be preserved
	value, _ := arr.Get(0)
	if value.ToInt() != 10 {
		t.Error("Values should be preserved after resize")
	}

	// New slots should be null
	value, _ = arr.Get(4)
	if value.Type() != types.TypeNull {
		t.Error("New slots should be null")
	}

	// Resize smaller
	arr.SetSize(2)
	if arr.GetSize() != 2 {
		t.Errorf("Expected size 2, got %d", arr.GetSize())
	}
}

func TestSplFixedArrayToArray(t *testing.T) {
	fixed := NewSplFixedArray(3)
	fixed.Set(0, types.NewInt(10))
	fixed.Set(1, types.NewInt(20))
	fixed.Set(2, types.NewInt(30))

	arr := fixed.ToArray()
	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}

	value, _ := arr.Get(types.NewInt(1))
	if value.ToInt() != 20 {
		t.Errorf("Expected 20, got %d", value.ToInt())
	}
}

// ============================================================================
// SplDoublyLinkedList Tests
// ============================================================================

func TestSplDoublyLinkedListPushPop(t *testing.T) {
	list := NewSplDoublyLinkedList()

	if !list.IsEmpty() {
		t.Error("New list should be empty")
	}

	// Push elements
	list.Push(types.NewInt(1))
	list.Push(types.NewInt(2))
	list.Push(types.NewInt(3))

	if list.Count() != 3 {
		t.Errorf("Expected count 3, got %d", list.Count())
	}

	// Pop from end
	value, ok := list.Pop()
	if !ok || value.ToInt() != 3 {
		t.Errorf("Expected 3, got %d", value.ToInt())
	}

	if list.Count() != 2 {
		t.Errorf("Expected count 2, got %d", list.Count())
	}
}

func TestSplDoublyLinkedListUnshiftShift(t *testing.T) {
	list := NewSplDoublyLinkedList()

	// Unshift elements (add to beginning)
	list.Unshift(types.NewInt(3))
	list.Unshift(types.NewInt(2))
	list.Unshift(types.NewInt(1))

	if list.Count() != 3 {
		t.Errorf("Expected count 3, got %d", list.Count())
	}

	// Shift from beginning
	value, ok := list.Shift()
	if !ok || value.ToInt() != 1 {
		t.Errorf("Expected 1, got %d", value.ToInt())
	}

	value, ok = list.Shift()
	if !ok || value.ToInt() != 2 {
		t.Errorf("Expected 2, got %d", value.ToInt())
	}
}

func TestSplDoublyLinkedListTopBottom(t *testing.T) {
	list := NewSplDoublyLinkedList()

	// Empty list
	_, ok := list.Top()
	if ok {
		t.Error("Expected false for top of empty list")
	}

	_, ok = list.Bottom()
	if ok {
		t.Error("Expected false for bottom of empty list")
	}

	list.Push(types.NewInt(10))
	list.Push(types.NewInt(20))
	list.Push(types.NewInt(30))

	// Top (end)
	value, ok := list.Top()
	if !ok || value.ToInt() != 30 {
		t.Errorf("Expected 30 at top, got %d", value.ToInt())
	}

	// Bottom (beginning)
	value, ok = list.Bottom()
	if !ok || value.ToInt() != 10 {
		t.Errorf("Expected 10 at bottom, got %d", value.ToInt())
	}

	// Should not modify list
	if list.Count() != 3 {
		t.Error("Top/Bottom should not modify list")
	}
}

func TestSplDoublyLinkedListMixed(t *testing.T) {
	list := NewSplDoublyLinkedList()

	// Mix push and unshift
	list.Push(types.NewInt(2))
	list.Unshift(types.NewInt(1))
	list.Push(types.NewInt(3))

	// Should be: 1, 2, 3
	value, _ := list.Shift()
	if value.ToInt() != 1 {
		t.Errorf("Expected 1, got %d", value.ToInt())
	}

	value, _ = list.Shift()
	if value.ToInt() != 2 {
		t.Errorf("Expected 2, got %d", value.ToInt())
	}

	value, _ = list.Shift()
	if value.ToInt() != 3 {
		t.Errorf("Expected 3, got %d", value.ToInt())
	}
}

// ============================================================================
// SplMaxHeap Tests
// ============================================================================

func TestSplMaxHeapBasic(t *testing.T) {
	heap := NewSplMaxHeap()

	if !heap.IsEmpty() {
		t.Error("New heap should be empty")
	}

	// Insert elements
	heap.Insert(types.NewInt(5))
	heap.Insert(types.NewInt(3))
	heap.Insert(types.NewInt(7))
	heap.Insert(types.NewInt(1))
	heap.Insert(types.NewInt(9))

	if heap.Count() != 5 {
		t.Errorf("Expected count 5, got %d", heap.Count())
	}

	// Max heap should extract in descending order
	value, ok := heap.Extract()
	if !ok || value.ToInt() != 9 {
		t.Errorf("Expected 9, got %d", value.ToInt())
	}

	value, ok = heap.Extract()
	if !ok || value.ToInt() != 7 {
		t.Errorf("Expected 7, got %d", value.ToInt())
	}

	value, ok = heap.Extract()
	if !ok || value.ToInt() != 5 {
		t.Errorf("Expected 5, got %d", value.ToInt())
	}
}

func TestSplMaxHeapTop(t *testing.T) {
	heap := NewSplMaxHeap()

	// Top on empty heap
	_, ok := heap.Top()
	if ok {
		t.Error("Expected false for top of empty heap")
	}

	heap.Insert(types.NewInt(10))
	heap.Insert(types.NewInt(5))
	heap.Insert(types.NewInt(15))

	// Top should return max without removing
	value, ok := heap.Top()
	if !ok || value.ToInt() != 15 {
		t.Errorf("Expected 15, got %d", value.ToInt())
	}

	if heap.Count() != 3 {
		t.Error("Top should not remove elements")
	}
}

// ============================================================================
// SplMinHeap Tests
// ============================================================================

func TestSplMinHeapBasic(t *testing.T) {
	heap := NewSplMinHeap()

	// Insert elements
	heap.Insert(types.NewInt(5))
	heap.Insert(types.NewInt(3))
	heap.Insert(types.NewInt(7))
	heap.Insert(types.NewInt(1))
	heap.Insert(types.NewInt(9))

	// Min heap should extract in ascending order
	value, ok := heap.Extract()
	if !ok || value.ToInt() != 1 {
		t.Errorf("Expected 1, got %d", value.ToInt())
	}

	value, ok = heap.Extract()
	if !ok || value.ToInt() != 3 {
		t.Errorf("Expected 3, got %d", value.ToInt())
	}

	value, ok = heap.Extract()
	if !ok || value.ToInt() != 5 {
		t.Errorf("Expected 5, got %d", value.ToInt())
	}
}

func TestSplMinHeapTop(t *testing.T) {
	heap := NewSplMinHeap()

	heap.Insert(types.NewInt(10))
	heap.Insert(types.NewInt(5))
	heap.Insert(types.NewInt(15))

	// Top should return min without removing
	value, ok := heap.Top()
	if !ok || value.ToInt() != 5 {
		t.Errorf("Expected 5, got %d", value.ToInt())
	}

	if heap.Count() != 3 {
		t.Error("Top should not remove elements")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestSplStackEmpty(t *testing.T) {
	stack := NewSplStack()

	_, ok := stack.Pop()
	if ok {
		t.Error("Pop on empty stack should return false")
	}

	_, ok = stack.Top()
	if ok {
		t.Error("Top on empty stack should return false")
	}
}

func TestSplQueueEmpty(t *testing.T) {
	queue := NewSplQueue()

	_, ok := queue.Dequeue()
	if ok {
		t.Error("Dequeue on empty queue should return false")
	}

	_, ok = queue.Peek()
	if ok {
		t.Error("Peek on empty queue should return false")
	}
}

func TestSplDoublyLinkedListEmpty(t *testing.T) {
	list := NewSplDoublyLinkedList()

	_, ok := list.Pop()
	if ok {
		t.Error("Pop on empty list should return false")
	}

	_, ok = list.Shift()
	if ok {
		t.Error("Shift on empty list should return false")
	}
}

func TestSplHeapEmpty(t *testing.T) {
	heap := NewSplMaxHeap()

	_, ok := heap.Extract()
	if ok {
		t.Error("Extract on empty heap should return false")
	}

	_, ok = heap.Top()
	if ok {
		t.Error("Top on empty heap should return false")
	}
}

func TestSplFixedArrayNegativeSize(t *testing.T) {
	arr := NewSplFixedArray(5)
	originalSize := arr.GetSize()

	// Negative size should be ignored
	arr.SetSize(-1)
	if arr.GetSize() != originalSize {
		t.Error("Negative size should be ignored")
	}
}
