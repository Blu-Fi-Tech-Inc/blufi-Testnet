package types

import (
	"fmt"
	"reflect"
)

// List is a generic list structure.
type List[T any] struct {
	Data []T
}

// NewList creates a new instance of List.
func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

// Get retrieves the element at the specified index.
func (l *List[T]) Get(index int) T {
	if index < 0 || index >= len(l.Data) {
		err := fmt.Sprintf("index (%d) out of range (length: %d)", index, len(l.Data))
		panic(err)
	}
	return l.Data[index]
}

// Insert adds an element to the end of the list.
func (l *List[T]) Insert(v T) {
	l.Data = append(l.Data, v)
}

// Clear removes all elements from the list.
func (l *List[T]) Clear() {
	l.Data = []T{}
}

// GetIndex finds and returns the index of the specified element.
// Returns -1 if the element is not found.
func (l *List[T]) GetIndex(v T) int {
	for i, item := range l.Data {
		if reflect.DeepEqual(v, item) {
			return i
		}
	}
	return -1
}

// Remove removes the first occurrence of the specified element from the list.
func (l *List[T]) Remove(v T) {
	index := l.GetIndex(v)
	if index != -1 {
		l.Pop(index)
	}
}

// Pop removes the element at the specified index from the list.
func (l *List[T]) Pop(index int) {
	if index < 0 || index >= len(l.Data) {
		err := fmt.Sprintf("index (%d) out of range (length: %d)", index, len(l.Data))
		panic(err)
	}
	l.Data = append(l.Data[:index], l.Data[index+1:]...)
}

// Contains checks if the specified element exists in the list.
func (l *List[T]) Contains(v T) bool {
	for _, item := range l.Data {
		if reflect.DeepEqual(item, v) {
			return true
		}
	}
	return false
}

// Last returns the last element of the list.
func (l List[T]) Last() T {
	if len(l.Data) == 0 {
		var zeroValue T // Return zero value if list is empty
		return zeroValue
	}
	return l.Data[len(l.Data)-1]
}

// Len returns the number of elements in the list.
func (l *List[T]) Len() int {
	return len(l.Data)
}
