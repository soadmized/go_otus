package hw04lrucache

type List interface {
	Len() int
	Front() *Node
	Back() *Node
	PushFront(v interface{}) *Node
	PushBack(v interface{}) *Node
	Remove(i *Node)
	MoveToFront(i *Node)
}

type Node struct {
	Value interface{}
	Next  *Node
	Prev  *Node
}

type list struct {
	head *Node
	tail *Node
	len  int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *Node {
	return l.head
}

func (l *list) Back() *Node {
	return l.tail
}

func (l *list) PushFront(v interface{}) *Node {
	node := &Node{Value: v}

	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		node.Next = l.head
		l.head.Prev = node
		l.head = node
	}

	l.len++

	return node
}

func (l *list) PushBack(v interface{}) *Node {
	node := &Node{Value: v}

	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		node.Prev = l.tail
		l.tail.Next = node
		l.tail = node
	}

	l.len++

	return node
}

func (l *list) Remove(node *Node) {
	if node == l.head {
		l.head = l.head.Next
		l.head.Prev = nil
		l.len--
		node.Next = nil

		return
	}

	if node == l.tail {
		l.tail = l.tail.Prev
		l.tail.Next = nil
		l.len--
		node.Prev = nil

		return
	}

	newPrev := node.Prev.Next.Prev
	newNext := node.Next.Prev.Next

	node.Prev.Next = newNext
	node.Next.Prev = newPrev

	node.Next = nil
	node.Prev = nil

	l.len--
}

func (l *list) MoveToFront(node *Node) {
	if node == l.head {
		return
	}

	if node == l.tail {
		l.tail = l.tail.Prev
		l.tail.Next = nil

		node.Prev = nil
		node.Next = l.head

		l.head.Prev = node
		l.head = node

		return
	}

	newPrev := node.Prev.Next.Prev
	newNext := node.Next.Prev.Next

	node.Prev.Next = newNext
	node.Next.Prev = newPrev

	node.Prev = nil
	node.Next = l.head

	l.head.Prev = node
	l.head = node
}

func NewList() List {
	return new(list)
}
