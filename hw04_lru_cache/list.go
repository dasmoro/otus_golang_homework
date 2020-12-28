package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	// Place your code here
	Len() int
	Front() *ListItem
	Back() *ListItem
	Remove(i *ListItem)
	PushFront(val interface{})
	PushBack(val interface{})
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value      interface{}
	Next, Prev *ListItem
	li         *list
}

type list struct {
	head  *ListItem
	tail  *ListItem
	count int
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(val interface{}) {
	if l.head == nil {
		l.head = &ListItem{Value: val, li: l}
		l.tail = l.head
	} else {
		l.head.Next = &ListItem{Value: val, Prev: l.head, Next: nil, li: l}
		l.head = l.head.Next
	}
	l.count++
}

func (l *list) PushBack(val interface{}) {
	if l.tail == nil {
		l.tail = &ListItem{Value: val, li: l}
		l.tail = l.head
	} else {
		l.tail.Prev = &ListItem{Value: val, Prev: nil, Next: l.tail, li: l}
		l.tail = l.tail.Prev
	}
	l.count++
}

func (l *list) Remove(i *ListItem) {
	if i.li != l {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.tail = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.head = i.Prev
	}
	if l.count >= 1 {
		l.count--
	}
	i.Next = nil
	i.Prev = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if i != l.head {
		l.Remove(i)
		l.head.Next = i
		i.Prev = l.head
		l.head = i
		l.count++
	}
}

func NewList() List {
	l := list{}
	l.count = 0
	return &l
}
