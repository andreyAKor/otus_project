package cache

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/andreyAKor/otus_project/internal/cache/list"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	cacheFileNamePattern = "image_previewer"
)

var (
	ErrKeyIsSet = errors.New("key is set")

	_ Cache = (*lruCache)(nil)
)

type lruCache struct {
	capacity int
	queue    list.List
	items    map[Key]*list.Item
	mux      sync.Mutex
}

type cacheItem struct {
	key    Key
	header http.Header
	file   string
}

func New(capacity int) (Cache, error) {
	return &lruCache{
		capacity: capacity,
		queue:    list.New(),
		items:    make(map[Key]*list.Item),
	}, nil
}

// Close cache.
func (l *lruCache) Close() error {
	return l.Clear()
}

// Returns a image data from the cache by key.
func (l *lruCache) Get(key Key) (http.Header, *[]byte, bool, error) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)

		content, err := ioutil.ReadFile(i.Value.(*cacheItem).file)
		if err != nil {
			return http.Header{}, nil, false, errors.Wrap(err, "reading cache-file fail")
		}

		return i.Value.(*cacheItem).header, &content, true, nil
	}

	return http.Header{}, nil, false, nil
}

// Adds a image data to the cache by key.
func (l *lruCache) Set(key Key, header http.Header, body *[]byte) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if _, ok := l.items[key]; ok {
		return ErrKeyIsSet
	}

	tmpfile, err := ioutil.TempFile("", cacheFileNamePattern)
	if err != nil {
		return errors.Wrap(err, "creating cache-file fail")
	}

	if _, err := tmpfile.Write(*body); err != nil {
		return errors.Wrap(err, "writing cache-file fail")
	}
	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "closing cache-file fail")
	}

	log.Info().Str("file", tmpfile.Name()).Msg("add cache-file")

	i := l.queue.PushFront(&cacheItem{key, header, tmpfile.Name()})
	l.items[key] = i

	if l.queue.Len() > l.capacity {
		i := l.queue.Back()
		l.queue.Remove(i)

		delete(l.items, i.Value.(*cacheItem).key)

		if err := l.remove(i.Value.(*cacheItem)); err != nil {
			return err
		}
	}

	return nil
}

// Clearing the cache.
func (l *lruCache) Clear() error {
	l.mux.Lock()
	defer l.mux.Unlock()

	for key, i := range l.items {
		l.queue.Remove(i)
		delete(l.items, key)

		if err := l.remove(i.Value.(*cacheItem)); err != nil {
			return err
		}
	}

	return nil
}

// Removig file from cache.
func (l *lruCache) remove(ci *cacheItem) error {
	if err := os.Remove(ci.file); err != nil {
		return errors.Wrap(err, "removing cache-file fail")
	}

	log.Info().Str("file", ci.file).Msg("remove cache-file")

	return nil
}
