package list

type List interface {
	Len() int                      // list length
	Front() *Item                  // first Item
	Back() *Item                   // last Item
	PushFront(v interface{}) *Item // add value to the beginning
	PushBack(v interface{}) *Item  // add value to the end
	Remove(i *Item)                // remove item
	MoveToFront(i *Item)           // move element to start
}

type Item struct {
	Value interface{} // value
	Next  *Item       // next element
	Prev  *Item       // previous item
}
