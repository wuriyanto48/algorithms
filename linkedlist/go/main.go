package main

import (
	"fmt"
)

// Node represent linkedlist node
type Node[T any] struct {
	Value T
	Prev  *Node[T]
	Next  *Node[T]
}

// NewNode func, the Node constructor
func NewNode[T any](value T) *Node[T] {
	return &Node[T]{
		Value: value,
		Prev:  nil,
		Next:  nil,
	}
}

// HasNext func,
func (n *Node[T]) HasNext() bool {
	return n.Next != nil
}

// HasPrev func
func (n *Node[T]) HasPrev() bool {
	return n.Prev != nil
}

// LinkedList, a doubly linkedlist
type LinkedList[T any] struct {
	len  uint
	head *Node[T]
	tail *Node[T]
}

// NewLinkedList the linkedlist constructor
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		len:  0,
		head: nil,
		tail: nil,
	}
}

// Len func, will return len of the linkedlist
func (l *LinkedList[T]) Len() uint {
	return l.len
}

// AddFront func, will add value to the front of the linkedlist
func (l *LinkedList[T]) AddFront(value T) {
	newNode := NewNode(value)
	if l.head == nil {
		l.tail = newNode
	} else {
		l.head.Prev = newNode
		newNode.Next = l.head
	}

	l.head = newNode
	l.len++
}

// AddLast func, will add value to the last of the linkedlist
func (l *LinkedList[T]) AddLast(value T) {
	newNode := NewNode(value)

	if l.head == nil {
		l.head = newNode
	} else {
		current := l.head
		for current.HasNext() {
			current = current.Next
		}

		current.Next = newNode
		newNode.Prev = current
	}

	l.tail = newNode
	l.len++
}

// AddAt func, will add value to the given pos of the linkedlist
func (l *LinkedList[T]) AddAt(value T, pos int) {
	newNode := NewNode(value)

	current := l.at(pos)
	if current == nil {
		return
	}

	if pos == 0 {
		if l.head == nil {
			l.head = newNode
			l.tail = newNode
		} else {
			newNode.Next = current
			current.Prev = newNode
			l.head = newNode
		}
	} else if uint(pos) == (l.len - 1) {
		if l.head == nil {
			l.head = newNode
			l.tail = newNode
		} else {
			newNode.Next = current
			newNode.Prev = current.Prev
			current.Prev.Next = newNode
			current.Prev = newNode
			l.tail = current
		}
	} else {
		if current.HasNext() {
			newNode.Next = current
			newNode.Prev = current.Prev
			current.Prev.Next = newNode
			current.Prev = newNode
		}
	}

	l.len++
}

// DeleteFirst func, will delete the first value of the linkedlist
func (l *LinkedList[T]) DeleteFirst() {
	if l.len <= 0 {
		return
	}

	tempHead := l.head
	if !l.head.HasNext() {
		l.head = nil
		l.tail = nil
	} else {
		l.head.Next.Prev = nil
	}

	l.head = tempHead.Next
	l.len--
}

// DeleteFirst func, will delete the last value of the linkedlist
func (l *LinkedList[T]) DeleteLast() {
	if l.len <= 0 {
		return
	}

	tempTail := l.tail
	if !l.tail.HasPrev() {
		l.head = nil
		l.tail = nil
	} else {
		l.tail.Prev.Next = nil
	}

	l.tail = tempTail.Prev
	l.len--
}

// IterForward func
func (l *LinkedList[T]) IterForward() <-chan *Node[T] {
	out := make(chan *Node[T])

	go func() {
		defer func() { close(out) }()
		current := l.head

		var i uint
		for i = 0; i < l.len; i++ {
			out <- current
			current = current.Next
		}
	}()

	return out
}

// IterBackward func
func (l *LinkedList[T]) IterBackward() <-chan *Node[T] {
	out := make(chan *Node[T])

	go func() {
		defer func() { close(out) }()
		current := l.tail

		var i uint
		for i = 0; i < l.len; i++ {
			out <- current
			current = current.Prev
		}
	}()

	return out
}

func (l *LinkedList[T]) at(pos int) *Node[T] {
	if l.len <= 0 {
		return nil
	}

	if uint(pos) > l.len {
		return nil
	}

	if pos < 0 {
		return nil
	}

	current := l.head
	for i := 0; i < pos; i++ {
		current = current.Next
	}

	return current
}

// DeleteAt func, will remove value of linkedlist at the given pos
func (l *LinkedList[T]) DeleteAt(pos int) {
	current := l.at(pos)
	if current == nil {
		return
	}

	if pos == 0 {
		current.Next.Prev = nil
		l.head = current.Next
	} else if uint(pos) == (l.len - 1) {
		current.Prev.Next = nil
		l.tail = current.Prev
	} else {
		if current.HasPrev() {
			current.Prev.Next = current.Next
		}
		
		if current.HasNext() {
			current.Next.Prev = current.Prev
		}
	}

	current.Next = nil
	current.Prev = nil
	l.len--
}

// At func, will return value of linkedlist at the given pos
func (l *LinkedList[T]) At(pos int) *T {
	current := l.at(pos)
	if current == nil {
		return nil
	}

	return &current.Value
}

func main() {
	l := NewLinkedList[string]()
	l.AddFront("Alex")
	l.AddFront("Bobby")
	l.AddFront("Cocy")
	l.AddFront("Dody")

	fmt.Println(l.Len())

	out := l.IterForward()
	for o := range out {
		fmt.Println(o)
	}

	l.DeleteAt(2)
	fmt.Println("--")

	l.AddFront("Wury")
	l.AddLast("Num")
	out = l.IterForward()
	for o := range out {
		fmt.Println(o)
	}

	fmt.Println("--")

	l.AddAt("Dul", 4)
	l.AddAt("Box", 0)
	out = l.IterForward()
	for o := range out {
		fmt.Println(o)
	}

	fmt.Println("--")
	v := l.At(6)
	if v != nil {
		fmt.Println(*v)
	}

	fmt.Println(l.Len())
}
