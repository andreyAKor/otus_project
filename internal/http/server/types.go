package server

// ImageRequest is the struct of request image for previewing.
type ImageRequest struct {
	Width, Height int
	Source        string
}
