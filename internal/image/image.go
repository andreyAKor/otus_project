package image

import (
	"bufio"
	"bytes"
	img "image"
	"image/jpeg"

	"github.com/pkg/errors"
)

var _ Image = (*image)(nil)

type image struct {
}

func New() (Image, error) {
	return &image{}, nil
}

func (i image) Resize(source []byte, width, height int) ([]byte, error) {
	config, _, err := img.DecodeConfig(bytes.NewReader(source))
	if err != nil {
		return nil, errors.Wrap(err, "image config decoding fail")
	}

	// Old image
	imgOld, _, err := img.Decode(bytes.NewReader(source))
	if err != nil {
		return nil, errors.Wrap(err, "image decoding fail")
	}

	oldSize := img.Point{config.Width, config.Height}
	newSize := img.Point{width, height}

	// New image
	newImg := img.NewRGBA(img.Rectangle{
		Max: newSize,
	})

	i.resizing(imgOld, newImg, oldSize, newSize)

	// Make new image
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)

	if err := jpeg.Encode(bw, newImg, nil); err != nil {
		return nil, errors.Wrap(err, "image encoding fail")
	}

	return buf.Bytes(), nil
}

func (i image) resizing(imgOld img.Image, newImg *img.RGBA, oldSize, newSize img.Point) {
	offsetX, offsetY, sizeX, sizeY := i.calcOffsets(
		float64(oldSize.X),
		float64(oldSize.Y),
		float64(newSize.X),
		float64(newSize.Y),
	)

	for y := 0; y < newSize.Y; y++ {
		for x := 0; x < newSize.X; x++ {
			oldX := int((float64(x) + offsetX) * sizeX)
			oldY := int((float64(y) + offsetY) * sizeY)

			newImg.Set(x, y, imgOld.At(oldX, oldY))
		}
	}
}

func (i image) calcOffsets(
	oldWidth, oldHeight, newWidth, newHeight float64,
) (
	offsetX, offsetY, sizeX, sizeY float64,
) {
	innerWidth := newWidth
	innerHeight := newHeight

	// Resizing to inner size
	if newWidth < oldWidth || newHeight < oldHeight {
		ratioOld := oldWidth / oldHeight

		if (newWidth / newHeight) < ratioOld {
			innerWidth = newHeight * ratioOld
		} else {
			innerHeight = newWidth / ratioOld
		}
	}

	offsetX = (innerWidth - newWidth) / 2
	offsetY = (innerHeight - newHeight) / 2

	if oldWidth != 0 && innerWidth != 0 {
		sizeX = oldWidth / innerWidth
	}
	if oldHeight != 0 && innerHeight != 0 {
		sizeY = oldHeight / innerHeight
	}

	return
}
