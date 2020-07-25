package client

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const defaultSchema = "http://"

var _ Client = (*client)(nil)

type client struct {
	timeout time.Duration
}

func New(timeout string) (Client, error) {
	timeoutDur, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "request timeout parsing fail (%s)", timeout)
	}

	return &client{timeoutDur}, nil
}

// Make request to source.
func (c *client) Request(source string, r *http.Request) (*http.Response, *[]byte, error) {
	log.Info().
		Str("source", source).
		Msg("before request")

	cl := c.prepareClient()

	req, err := c.prepareRequest(source, r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "preparing request fail")
	}

	rsp, content, err := c.request(cl, req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "request fail")
	}

	log.Info().
		Int("length", len(*content)).
		Str("source", source).
		Msg("after request")

	return rsp, content, nil
}

// Preparing client.
func (c *client) prepareClient() *http.Client {
	return &http.Client{
		Timeout: c.timeout,
	}
}

// Preparing request.
func (c *client) prepareRequest(source string, r *http.Request) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		r.Context(),
		r.Method,
		c.normalizeURL(source),
		r.Body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "init request with context fail")
	}

	// Copy headers
	req.Header = r.Header

	return req, nil
}

// Make request.
func (c *client) request(client *http.Client, req *http.Request) (*http.Response, *[]byte, error) {
	rsp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "client request fail")
	}
	defer rsp.Body.Close()

	content, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "reading from client fail")
	}
	if _, err := io.Copy(ioutil.Discard, rsp.Body); err != nil {
		return nil, nil, errors.Wrap(err, "copying from response body fail")
	}

	return rsp, &content, nil
}

func (c client) normalizeURL(url string) string {
	return defaultSchema + url
}
