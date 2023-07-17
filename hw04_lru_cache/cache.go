package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*Node
}

func (l lruCache) Set(key Key, value any) bool {
	_, exist := l.Get(key)
	if exist {
		node := l.items[key]
		node.Value = value
		l.queue.MoveToFront(node)

		return true
	}

	if l.queue.Len() > l.capacity {
		last := l.queue.Back()
		l.queue.Remove(last)
		// TODO: delete(l.items, key???)
	}

	node := l.queue.PushFront(value)
	l.items[key] = node

	return false
}

func (l lruCache) Get(key Key) (any, bool) {
	//TODO implement me
	panic("implement me")
}

func (l lruCache) Clear() {
	//TODO implement me
	panic("implement me")
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*Node, capacity),
	}
}
