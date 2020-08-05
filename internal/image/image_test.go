package image

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcOffsets(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		offsetX, offsetY, sizeX, sizeY := image.calcOffsets(image{}, 1000, 800, 0, 0)
		require.Equal(t, 0.0, offsetX)
		require.Equal(t, 0.0, offsetY)
		require.Equal(t, 0.0, sizeX)
		require.Equal(t, 0.0, sizeY)

		offsetX, offsetY, sizeX, sizeY = image.calcOffsets(image{}, 0, 0, 500, 200)
		require.Equal(t, 0.0, offsetX)
		require.Equal(t, 0.0, offsetY)
		require.Equal(t, 0.0, sizeX)
		require.Equal(t, 0.0, sizeY)

		offsetX, offsetY, sizeX, sizeY = image.calcOffsets(image{}, 0, 0, 0, 0)
		require.Equal(t, 0.0, offsetX)
		require.Equal(t, 0.0, offsetY)
		require.Equal(t, 0.0, sizeX)
		require.Equal(t, 0.0, sizeY)
	})
	t.Run("resizing", func(t *testing.T) {
		t.Run("from 1000x800 to 500x200", func(t *testing.T) {
			offsetX, offsetY, sizeX, sizeY := image.calcOffsets(image{}, 1000, 800, 500, 200)
			require.Equal(t, 0.0, offsetX)
			require.Equal(t, 100.0, offsetY)
			require.Equal(t, 2.0, sizeX)
			require.Equal(t, 2.0, sizeY)
		})
		t.Run("from 1000x800 to 200x500", func(t *testing.T) {
			offsetX, offsetY, sizeX, sizeY := image.calcOffsets(image{}, 1000, 800, 200, 500)
			require.Equal(t, 212.5, offsetX)
			require.Equal(t, 0.0, offsetY)
			require.Equal(t, 1.6, sizeX)
			require.Equal(t, 1.6, sizeY)
		})
		t.Run("from 1000x800 to 2000x1600", func(t *testing.T) {
			offsetX, offsetY, sizeX, sizeY := image.calcOffsets(image{}, 1000, 800, 2000, 1600)
			require.Equal(t, 0.0, offsetX)
			require.Equal(t, 0.0, offsetY)
			require.Equal(t, 0.5, sizeX)
			require.Equal(t, 0.5, sizeY)
		})
	})
}

func TestValidateImageSize(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		i := &image{0, 0}
		res := i.ValidateImageSize(0, 0)
		require.True(t, res)

		res = i.ValidateImageSize(1, 1)
		require.False(t, res)
	})
	t.Run("normal", func(t *testing.T) {
		i := &image{1, 1}
		res := i.ValidateImageSize(0, 0)
		require.True(t, res)

		res = i.ValidateImageSize(1, 1)
		require.True(t, res)

		res = i.ValidateImageSize(2, 1)
		require.False(t, res)

		res = i.ValidateImageSize(1, 2)
		require.False(t, res)

		res = i.ValidateImageSize(2, 2)
		require.False(t, res)
	})
}
