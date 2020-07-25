package client

import "net/http"

//go:generate mockgen -source=$GOFILE -destination ./mocks/mock_client.go -package mocks Client
type Client interface {
	Request(source string, r *http.Request) (*http.Response, *[]byte, error)
}
