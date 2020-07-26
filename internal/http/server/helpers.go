package server

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

var (
	ErrInvalidURIFormat = errors.New("invalid URI format")

	reParseURI = regexp.MustCompile(`/([\d]+)/([\d]+)/(.+)`)
)

// Prepare ImageRequest scructure from uri-data.
func (s Server) parseURI(uri string) (ImageRequest, error) {
	var (
		ir  ImageRequest
		err error
	)

	list := reParseURI.FindStringSubmatch(uri)
	if len(list) == 0 || len(list) < 4 {
		return ir, ErrInvalidURIFormat
	}

	ir.Width, err = strconv.Atoi(list[1])
	if err != nil {
		return ir, errors.Wrap(err, "can't convert width")
	}

	ir.Height, err = strconv.Atoi(list[2])
	if err != nil {
		return ir, errors.Wrap(err, "can't convert height")
	}

	ir.Source = list[3]

	return ir, nil
}
