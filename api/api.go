package api

import (
	"errors"
	"image"
)

var (
	FormatInvalidErr = errors.New("format:invalid")
	ImageTooWideErr  = errors.New("image:too_wide")
	ImageTooLongErr  = errors.New("image:too_long")
)

type IImageProcessor interface {
	Process(image.Image, string) (image.Image, error)
}
