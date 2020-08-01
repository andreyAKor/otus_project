package cache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const multithreadingItems = 1_000

//nolint:funlen
func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c, _ := New(10)
		defer c.Close()

		_, ok, _ := c.Get("aaa")
		require.False(t, ok)

		_, ok, _ = c.Get("bbb")
		require.False(t, ok)
	})
	t.Run("simple", func(t *testing.T) {
		c, _ := New(5)
		defer c.Close()

		v1 := []byte("100")
		err := c.Set("aaa", v1)
		require.NoError(t, err)

		v2 := []byte("200")
		err = c.Set("bbb", v2)
		require.NoError(t, err)

		val, ok, _ := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, v1, val)

		val, ok, _ = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, v2, val)

		v3 := []byte("300")
		err = c.Set("aaa", v3)
		require.Equal(t, ErrKeyIsSet, errors.Cause(err))

		val, ok, _ = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, v1, val)

		val, ok, _ = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
	t.Run("purge logic", func(t *testing.T) {
		c, _ := New(5)
		defer c.Close()

		v1 := []byte("100")
		err := c.Set("aaa", v1)
		require.NoError(t, err)

		// [aaa] => [9, 8, 7, 6, 5]
		for i := 0; i < 10; i++ {
			s := strconv.Itoa(i)
			v := []byte(s)
			_ = c.Set(Key(s), v)
		}

		val, ok, _ := c.Get("9")
		require.True(t, ok)
		require.Equal(t, []byte("9"), val)

		val, ok, _ = c.Get("5")
		require.True(t, ok)
		require.Equal(t, []byte("5"), val)

		val, ok, _ = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		// [9, 8, 7, 6, 5] => [6, 5, 9, 8, 7]
		for i := 5; i < 7; i++ {
			_, _, _ = c.Get(Key(strconv.Itoa(i)))
		}

		// [6, 5, 9, 8, 7] => [0, 1, 2, 6, 5]
		for i := 0; i < 3; i++ {
			s := strconv.Itoa(i)
			v := []byte(s)
			_ = c.Set(Key(s), v)
		}

		val, ok, _ = c.Get("0")
		require.True(t, ok)
		require.Equal(t, []byte("0"), val)

		val, ok, _ = c.Get("5")
		require.True(t, ok)
		require.Equal(t, []byte("5"), val)
	})
	t.Run("additional logic", func(t *testing.T) {
		c, _ := New(5)
		defer c.Close()

		// [99, 98, 97, 96, 95]
		for i := 0; i < 100; i++ {
			s := strconv.Itoa(i)
			v := []byte(s)
			_ = c.Set(Key(s), v)
		}

		val, ok, _ := c.Get("95")
		require.True(t, ok)
		require.Equal(t, []byte("95"), val)

		// [99, 98, 97, 96, 95] => []
		_ = c.Clear()

		val, ok, _ = c.Get("95")
		require.False(t, ok)
		require.Nil(t, val)

		for i := 100; i < 200; i++ {
			s := strconv.Itoa(i)
			v := []byte(s)
			_ = c.Set(Key(s), v)
		}

		val, ok, _ = c.Get("195")
		require.True(t, ok)
		require.Equal(t, []byte("195"), val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c, _ := New(10)
	defer c.Close()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < multithreadingItems; i++ {
			s := strconv.Itoa(i)
			v := []byte(s)
			_ = c.Set(Key(s), v)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < multithreadingItems; i++ {
			_, _, _ = c.Get(Key(strconv.Itoa(rand.Intn(multithreadingItems))))
		}
	}()

	wg.Wait()
}
