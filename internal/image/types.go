package image

//go:generate mockgen -source=$GOFILE -destination ./mocks/mock_image.go -package mocks Image
type Image interface {
	Resize(source []byte, width, height int) ([]byte, error)
}
