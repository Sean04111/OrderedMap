package OrderedMap

import "sync"

type OrderMapImp struct {
	_ nocopy

	mu       sync.Mutex
	queryMap map[string]*node
	head     *node
	tail     *node
}

type node struct {
	key   string
	value []byte
	next  *node
	pre   *node
}

func (om *OrderMapImp) Init() {
	om.queryMap = make(map[string]*node)
	om.head = &node{}
	om.tail = &node{}
	om.head.next = om.tail
	om.tail.pre = om.head
}

func (om *OrderMapImp) Add(key, value string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, ok := om.queryMap[key]; ok {
		return
	}

	newNode := &node{}
	newNode.key = key
	newNode.value = []byte(value)

	om.tail.pre.next = newNode
	newNode.pre = om.tail.pre
	newNode.next = om.tail
	om.tail.pre = newNode

	om.queryMap[key] = newNode
}

func (om *OrderMapImp) Del(key string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	node, ok := om.queryMap[key]
	if !ok {
		return
	}

	node.pre.next = node.next
	node.next.pre = node.pre

	delete(om.queryMap, key)
}

func (om *OrderMapImp) Range() []string {
	om.mu.Lock()
	defer om.mu.Unlock()
	var ret []string
	for p := om.head.next; p != om.tail; p = p.next {
		ret = append(ret, string(p.value))
	}
	return ret
}
