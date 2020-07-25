package list

var _ List = (*list)(nil)

// The double-linked list implementation.
type list struct {
	front *Item // The first item in the list
	back  *Item // The last item in the list
	size  int   // List size
}

func New() List {
	return &list{}
}

// List length.
func (l *list) Len() int {
	return l.size
}

// The first item in the list.
func (l *list) Front() *Item {
	return l.front
}

// The last item in the list.
func (l *list) Back() *Item {
	return l.back
}

// Adds a value to the beginning of the list.
func (l *list) PushFront(v interface{}) *Item {
	if l.front == nil {
		l.front = &Item{
			Value: v,
		}
		l.back = l.front
	} else {
		l.front.Next = &Item{
			Value: v,
			Prev:  l.front,
		}

		l.front = l.front.Next
	}

	l.size++

	return l.front
}

// Adds a value to the end of the list.
func (l *list) PushBack(v interface{}) *Item {
	if l.back == nil {
		l.back = &Item{
			Value: v,
		}
		l.front = l.back
	} else {
		l.back.Prev = &Item{
			Value: v,
			Next:  l.back,
		}

		l.back = l.back.Prev
	}

	l.size++

	return l.back
}

// Removes an item from the list.
func (l *list) Remove(i *Item) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.back = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.front = i.Prev
	}

	l.size--
}

// Moves the item to the beginning of the list.
func (l *list) MoveToFront(i *Item) {
	if i == nil || i == l.front {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}
