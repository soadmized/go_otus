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

type cacheNode struct {
	key   Key
	value any
}

func (l lruCache) Set(key Key, value any) bool {
	_, exist := l.Get(key)
	if exist {
		node := l.items[key]
		node.Value = cacheNode{
			key:   key,
			value: value,
		}
		l.queue.MoveToFront(node)

		return true
	}

	if l.queue.Len() >= l.capacity {
		last := l.queue.Back()
		l.queue.Remove(last)

		delete(l.items, last.Value.(cacheNode).key)
	}

	node := l.queue.PushFront(cacheNode{key: key, value: value})
	l.items[key] = node

	return false
}

func (l lruCache) Get(key Key) (any, bool) {
	node, exist := l.items[key]
	if exist {
		l.queue.MoveToFront(node)

		return (node.Value).(cacheNode).value, true
	}

	return nil, false
}

func (l lruCache) Clear() {
	newQueue := NewList()
	newItems := make(map[Key]*Node)

	l.queue = newQueue
	l.items = newItems
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*Node, capacity),
	}
}
