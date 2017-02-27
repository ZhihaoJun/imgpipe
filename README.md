# imgpipe
image pipline

## dependencies
[github.com/anthonynsimon/bild/imgio](github.com/anthonynsimon/bild/imgio)

## usages and example
using echo framework

``` golang
import "zhihaojun.com/imgpipe"

func (a *App) apiUpload(c echo.Context) error {
	if c.Request().ContentLength > a.allowSize {
		return a.response(c, http.StatusConflict, map[string]interface{}{
			"error": "file:too_big",
			"msg":   "file is too big",
		})
	}

	file, err := c.FormFile("file-0")
	if err != nil {
		log.Println("[apiUpload] get form file failed", err)
		return a.errorHandle(c, err)
	}
	src, err := file.Open()
	if err != nil {
		log.Println("[apiUpload] file open failed", err)
		return a.errorHandle(c, err)
	}
	defer src.Close()

	pipeline := imgpipe.NewImageProcessPipeline()
	sizeDeny := imgpipe.NewImageSizeDeny(4500, 4500)
	typeDeny := imgpipe.NewImageTypeDeny("png", "jpg", "jpeg")
	saver := imgpipe.NewImageSaver(a.uploadRoot)
	pipeline.AddProcessor(typeDeny)
	pipeline.AddProcessor(sizeDeny)
	pipeline.AddProcessor(saver)

	if err := pipeline.Process(src); err != nil {
		log.Println("[apiUpload] pipeline process failed", err)
		return a.errorHandle(c, err)
	}

	return a.response(c, http.StatusOK, map[string]interface{}{
		"error": "ok",
		"msg":   "upload success",
		"url":   saver.SavedPath(),
	})
}
```

## custom processor
implement go interface

``` golang
type IImageProcessor interface {
	Process(image.Image, string) (image.Image, error)
}
```

for example

``` golang
type MyProcessor struct {

}

func (mp *MyProcessor) Process(img image.Image, string) (image.Image, error) {

}

func main() {
  p := &MyProcessor{}
  pipeline := imgpipe.NewImageProcessPipeline()
  pipeline.AddProcessor(p)
}
```