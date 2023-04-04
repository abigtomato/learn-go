package container

import (
	"container/list"
	"fmt"
	"testing"
)

// 双端队列
type Deque[T any] struct {
	l *list.List
}

func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{
		l: list.New(),
	}
}

func (d *Deque[T]) AddFirst(elem T) {
	d.l.PushFront(elem)
}

func (d *Deque[T]) AddLast(elem T) {
	d.l.PushBack(elem)
}

func (d *Deque[T]) RemoveFirst() T {
	return d.l.Remove(d.l.Front()).(T)
}

func (d *Deque[T]) RemoveLast() T {
	return d.l.Remove(d.l.Back()).(T)
}

// 栈
type Stack[T any] struct {
	l *list.List
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		l: list.New(),
	}
}

func (s *Stack[T]) Push(elem T) {
	s.l.PushBack(elem)
}

func (s *Stack[T]) Pop() T {
	return s.l.Remove(s.l.Back()).(T)
}

func (s *Stack[T]) Peek() T {
	return s.l.Back().Value.(T)
}

type kv[K comparable, V any] struct {
	k K
	v V
}

// LRU
type LRU[K comparable, V any] struct {
	l    *list.List
	m    map[K]*list.Element
	size int
}

func NewLRU[K comparable, V any](size int) *LRU[K, V] {
	return &LRU[K, V]{
		l:    list.New(),
		m:    make(map[K]*list.Element, size),
		size: size,
	}
}

func (l *LRU[K, V]) Put(k K, v V) {
	// 如果k已经存在，直接把它移到最后面，然后设置新值
	if elem, ok := l.m[k]; ok {
		l.l.MoveToBack(elem)
		keyValue := elem.Value.(kv[K, V])
		keyValue.v = v
		return
	}

	// 如果已经到达最大尺寸，先剔除一个元素
	if l.l.Len() == l.size {
		front := l.l.Front()
		l.l.Remove(front)
		delete(l.m, front.Value.(kv[K, V]).k)
	}

	// 添加元素
	l.m[k] = l.l.PushBack(kv[K, V]{k, v})
}

func (l *LRU[K, V]) Get(k K) (V, bool) {
	// 如果存在移动到尾部，然后返回
	if elem, ok := l.m[k]; ok {
		l.l.MoveToBack(elem)
		return elem.Value.(kv[K, V]).v, true
	}
	// 不存在返回空值和false
	var v V
	return v, false
}

func TestList(t *testing.T) {
	// 使用list.New()直接初始化
	l := list.New()
	// 链表头部添加元素
	l.PushFront(1)
	fmt.Println(l.Front().Value)

	// 使用list.List{}延迟初始化
	// 在调用PushFront()、PushBack()、PushFrontList()、PushBackList() 时会调用 lazyInit() 检查是否已经初始化
	// 如果没有初始化则调用 Init() 进行初始化
	l2 := list.List{}
	// 链表尾部添加元素
	elem := l2.PushBack(2)
	elem2 := l2.PushBack(4)
	fmt.Println(l2.Back().Value)

	// 分别是获取头元素、获取尾元素、获取长度，都不会对链表修改
	fmt.Println(l.Front().Value)
	fmt.Println(l.Back().Value)
	fmt.Println(l.Len())

	// 分别是在链表头部或在尾部插入链表
	l2.PushFrontList(l)
	l2.PushBackList(l)

	// 分别是在某个元素前插入，在某个元素后插入
	l2.InsertBefore(2, elem)
	l2.InsertAfter(3, elem2)

	// 分别是移动元素到某个元素前面、移动元素到某个元素后面、移动元素到头部、移动元素到尾部
	l2.MoveAfter(elem2, elem)
	l2.MoveBefore(elem, elem2)
	l2.MoveToFront(elem)
	l2.MoveToBack(elem2)

	// 遍历链表
	// list是一个双向循环链表，尾节点的后继指针会指向头节点
	for cur := l2.Front(); cur != l2.Front().Prev(); cur = cur.Next() {
		fmt.Println(cur.Value)
	}
}
