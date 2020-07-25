package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := New()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := New()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, l.Len(), 3)

		middle := l.Back().Next // 20
		l.Remove(middle)        // [10, 30]
		require.Equal(t, l.Len(), 2)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, l.Len(), 7)
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{50, 30, 10, 40, 60, 80, 70}, elems)
	})

	t.Run("additional complex", func(t *testing.T) {
		l := New()

		l.Remove(l.Back()) // []
		require.Equal(t, 0, l.Len())

		l.Remove(l.Front()) // []
		require.Equal(t, 0, l.Len())

		l.PushBack(nil)    // [nil]
		l.PushFront(10)    // [10, nil]
		l.PushFront(20.56) // [20.56, 10, nil]
		l.PushBack('\n')   // [20.56, 10, nil, '\n']
		l.PushBack("mumu") // [20.56, 10, nil, '\n', "mumu"]
		require.Equal(t, 5, l.Len())

		l.Remove(nil)
		require.Equal(t, 5, l.Len())

		l.MoveToFront(nil)
		require.Equal(t, 5, l.Len())
		require.Equal(t, 20.56, l.Front().Value)
		require.Equal(t, "mumu", l.Back().Value)

		l.Remove(l.Back()) // [20.56, 10, nil, '\n']
		require.Equal(t, 4, l.Len())

		l.Remove(l.Front()) // [10, nil, '\n']
		require.Equal(t, 3, l.Len())

		l.Remove(l.Back()) // [10, nil]
		require.Equal(t, 2, l.Len())

		l.Remove(l.Front()) // [nil]
		require.Equal(t, 1, l.Len())

		l.Remove(l.Back()) // []
		require.Equal(t, 0, l.Len())
	})
}
