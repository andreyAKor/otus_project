package cache

import (
	"io"
	"net/http"
)

type Key string

//go:generate mockgen -source=$GOFILE -destination ./mocks/mock_cache.go -package mocks Cache
type Cache interface {
	Get(key Key) (http.Header, *[]byte, bool, error)
	Set(key Key, header http.Header, body *[]byte) error
	Clear() error
	io.Closer
}
