package cache

import (
	"io"
)

type Key string

//go:generate mockgen -source=$GOFILE -destination ./mocks/mock_cache.go -package mocks Cache
type Cache interface {
	Get(key Key) ([]byte, bool, error)
	Set(key Key, body []byte) error
	Clear() error
	io.Closer
}
