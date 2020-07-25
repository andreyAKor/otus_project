package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeUrl(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		url := client.normalizeURL(client{}, "")
		require.Equal(t, defaultSchema, url)
	})
	t.Run("invalid URI format", func(t *testing.T) {
		domain := "lala.be"
		url := client.normalizeURL(client{}, domain)
		require.Equal(t, defaultSchema+domain, url)
	})
}
