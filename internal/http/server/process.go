package server

import (
	"net/http"

	"github.com/andreyAKor/otus_project/internal/cache"

	"github.com/pkg/errors"
)

var ErrBadGateway = errors.New("bad gateway")

// Get image and process preparing preview image.
func (s *Server) process(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	res, ok, err := s.cache.Get(cache.Key(r.URL.RequestURI()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.Wrap(err, "getting image from cache fail")
	}

	if !ok {
		res, err = s.preparingImage(w, r)
		if err != nil {
			return nil, errors.Wrap(err, "preparing image fail")
		}
	}

	return res, nil
}

//nolint:bodyclose
func (s *Server) preparingImage(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	// Prepare request params
	ir, err := s.parseURI(r.URL.RequestURI())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.Wrap(err, "request uri parsing fail")
	}

	// Retrieving image content
	rsp, content, err := s.client.Request(ir.Source, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.Wrap(err, "client request fail")
	}
	if rsp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadGateway)
		return nil, ErrBadGateway
	}

	// Resizing an image
	content, err = s.image.Resize(content, ir.Width, ir.Height)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.Wrap(err, "image runner fail")
	}

	if err := s.cache.Set(cache.Key(r.URL.RequestURI()), content); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.Wrap(err, "setting image to cache fail")
	}

	return content, nil
}
