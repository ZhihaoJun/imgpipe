package imgpipe

import (
	"image"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"zhihaojun.com/imgpipe/api"

	"github.com/anthonynsimon/bild/imgio"
)

type ImageProcessPipeline struct {
	processors []api.IImageProcessor
}

func NewImageProcessPipeline() *ImageProcessPipeline {
	return &ImageProcessPipeline{
		processors: []api.IImageProcessor{},
	}
}

func (ipp *ImageProcessPipeline) AddProcessor(p api.IImageProcessor) {
	ipp.processors = append(ipp.processors, p)
}

func (ipp *ImageProcessPipeline) Process(reader io.Reader) error {
	src, format, err := image.Decode(reader)
	log.Println("[image format]", format)
	if err != nil {
		return err
	}
	for _, processor := range ipp.processors {
		src, err = processor.Process(src, format)
		if err != nil {
			return err
		}
	}
	return nil
}

type ImageTypeDeny struct {
	allowFormats []string
}

func NewImageTypeDeny(allowFormats ...string) *ImageTypeDeny {
	return &ImageTypeDeny{
		allowFormats: allowFormats,
	}
}

func (itd *ImageTypeDeny) inStrs(s string, strs []string) bool {
	for _, str := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func (itd *ImageTypeDeny) Process(src image.Image, format string) (image.Image, error) {
	if itd.inStrs(format, itd.allowFormats) == false {
		return src, api.FormatInvalidErr
	}
	return src, nil
}

type ImageSizeDeny struct {
	allowWidth  int
	allowHeight int
}

func NewImageSizeDeny(width, height int) *ImageSizeDeny {
	return &ImageSizeDeny{
		allowWidth:  width,
		allowHeight: height,
	}
}

func (isd *ImageSizeDeny) Process(src image.Image, format string) (image.Image, error) {
	bounds := src.Bounds()
	size := bounds.Size()
	if size.X > isd.allowWidth {
		return src, api.ImageTooWideErr
	} else if size.Y > isd.allowHeight {
		return src, api.ImageTooLongErr
	}
	return src, nil
}

type ImageSizeReducer struct {
}

func (isr *ImageSizeReducer) Process(src image.Image, format string) (image.Image, error) {
	return src, nil
}

type ImageSaver struct {
	path      string
	savedPath string
}

func NewImageSaver(path string) *ImageSaver {
	return &ImageSaver{
		path: path,
	}
}

func (is *ImageSaver) generateFilepath(src image.Image, format string) (string, error) {
	r := randStringRunes(8)
	fname := strings.Join([]string{r, ".", format}, "")
	return filepath.Abs(filepath.Join(is.path, r[:2], r[2:4], fname))
}

func (is *ImageSaver) formatConvert(format string) imgio.Format {
	switch format {
	case "jpeg":
		return imgio.JPEG
	case "jpg":
		return imgio.JPEG
	case "png":
		return imgio.PNG
	}
	return imgio.PNG
}

func (is *ImageSaver) Process(src image.Image, format string) (image.Image, error) {
	fullpath, err := is.generateFilepath(src, format)
	if err != nil {
		return src, err
	}
	log.Printf("[ImageSaver] fullpath %s\n", fullpath)
	// mkdir
	dir := filepath.Dir(fullpath)
	if err := os.MkdirAll(dir, os.ModeDir|0775); err != nil {
		if strings.HasSuffix(err.Error(), "file exists") == false {
			return src, err
		}
	}
	is.savedPath = fullpath

	// save file
	f, err := os.Create(fullpath)
	if err != nil {
		return src, err
	}
	defer f.Close()
	err = imgio.Encode(f, src, is.formatConvert(format))
	return src, err
}

func (is *ImageSaver) SavedPath() string {
	return is.savedPath
}

func (is *ImageSaver) Reset() {
	is.savedPath = ""
}
