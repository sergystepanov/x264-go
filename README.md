## x264-go/v2

`x264-go/v2` provides H.264/MPEG-4 AVC codec encoder based on [x264](https://www.videolan.org/developers/x264.html) library and
original gen2brain/x264-go wrapper.

By default it will use installed in the system shared/static library.
If toy want to use old C source code included in the package then build with `-tags legacy`.

### Installation

    go get -u github.com/sergystepanov/x264-go/v2

### Examples

See [screengrab](examples/screengrab/screengrab.go) example.

### Usage

```go
package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"

	"github.com/sergystepanov/x264-go/v2"
)

func main() {
	buf := bytes.NewBuffer(make([]byte, 0))

	opts := &x264.Options{
		Width:     640,
		Height:    480,
		FrameRate: 25,
		Tune:      "zerolatency",
		Preset:    "veryfast",
		Profile:   "baseline",
		LogLevel:  x264.LogDebug,
	}

	enc, err := x264.NewEncoder(buf, opts)
	if err != nil {
		panic(err)
	}

	img := x264.NewYCbCr(image.Rect(0, 0, opts.Width, opts.Height))
	draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)

	for i := 0; i < opts.Width/2; i++ {
		img.Set(i, opts.Height/2, color.RGBA{255, 0, 0, 255})

		err = enc.Encode(img)
		if err != nil {
			panic(err)
		}
	}

	err = enc.Flush()
	if err != nil {
		panic(err)
	}

	err = enc.Close()
	if err != nil {
		panic(err)
	}
}
```
