package server

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestParseUri(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, err := Server.parseURI(Server{}, "")
		require.Equal(t, ErrInvalidURIFormat, errors.Cause(err))
	})
	t.Run("invalid URI format", func(t *testing.T) {
		t.Run("invalid width", func(t *testing.T) {
			_, err := Server.parseURI(Server{}, "//200/img.com/image.jpg")
			require.Equal(t, ErrInvalidURIFormat, errors.Cause(err))

			_, err = Server.parseURI(Server{}, "/1./200/img.com/image.jpg")
			require.Equal(t, ErrInvalidURIFormat, errors.Cause(err))
		})

		t.Run("invalid height", func(t *testing.T) {
			_, err := Server.parseURI(Server{}, "/300//img.com/image.jpg")
			require.Equal(t, ErrInvalidURIFormat, errors.Cause(err))

			_, err = Server.parseURI(Server{}, "/300/.100/img.com/image.jpg")
			require.Equal(t, ErrInvalidURIFormat, errors.Cause(err))
		})

		t.Run("invalid source", func(t *testing.T) {
			_, err := Server.parseURI(Server{}, "/300/200/")
			require.Equal(t, ErrInvalidURIFormat, errors.Cause(err))
		})
	})
}
