// +build !legacy

package x264

import (
	"fmt"
	"github.com/sergystepanov/x264-go/v2/x264c/color"
	x264c "github.com/sergystepanov/x264-go/v2/x264c/external"
	"image"
	"io"
	"log"
)

/*
#include <stdlib.h>
*/
import "C"

// Options represent encoding options.
type Options struct {
	// Frame width.
	Width int
	// Frame height.
	Height int
	// Frame rate.
	FrameRate int
	// Tunings: film, animation, grain, stillimage, psnr, ssim, fastdecode, zerolatency.
	Tune string
	// Presets: ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow, placebo.
	Preset string
	// Profiles: baseline, main, high, high10, high422, high444.
	Profile string
	// Log level.
	LogLevel int32
}

// Encoder type.
type Encoder struct {
	e *x264c.T
	w io.Writer

	img  *color.YCbCr
	opts *Options

	csp int32
	pts int64

	nnals int32
	nals  []*x264c.X264NalT

	picIn x264c.Picture

	// ticks per frame
	tpf int64
}

// NewEncoder returns new x264 encoder.
func NewEncoder(w io.Writer, opts *Options) (e *Encoder, err error) {
	e = &Encoder{}

	e.w = w
	e.pts = 0
	e.opts = opts

	e.csp = x264c.X264CspI420

	e.nals = make([]*x264c.X264NalT, 3)
	e.img = color.NewYCbCr(image.Rect(0, 0, e.opts.Width, e.opts.Height))

	param := x264c.X264ParamT{}

	if e.opts.Preset != "" && e.opts.Profile != "" {
		ret := x264c.ParamDefaultPreset(&param, e.opts.Preset, e.opts.Tune)
		if ret < 0 {
			err = fmt.Errorf("x264: invalid preset/tune name")
			return
		}
	} else {
		x264c.ParamDefault(&param)
	}

	//param.IThreads = 1
	param.IBitdepth = 8
	param.ICsp = e.csp
	param.IWidth = int32(e.opts.Width)
	param.IHeight = int32(e.opts.Height)
	param.BVfrInput = 0
	param.BRepeatHeaders = 1
	param.BAnnexb = 1
	param.ILogLevel = e.opts.LogLevel
	param.IKeyintMax = 60
	param.BIntraRefresh = 1
	param.IFpsNum = 60
	param.IFpsDen = 1

	param.Rc.IRcMethod = x264c.X264RcCrf
	param.Rc.FRfConstant = 28

	//param.BVfrInput = 1
	param.ITimebaseNum = 1
	param.ITimebaseDen = 1000

	e.tpf = int64(param.ITimebaseDen * param.IFpsDen / param.ITimebaseNum / param.IFpsNum)

	if e.opts.Profile != "" {
		ret := x264c.ParamApplyProfile(&param, e.opts.Profile)
		if ret < 0 {
			err = fmt.Errorf("x264: invalid profile name")
			return
		}
	}

	// Allocate on create instead while encoding
	var picIn x264c.Picture
	x264c.PictureInit(&picIn)
	//ret := x264c.PictureAlloc(&picIn, e.csp, int32(e.opts.Width), int32(e.opts.Height))
	//if ret < 0 {
	//err = fmt.Errorf("x264: cannot allocate picture")
	//return
	//}

	e.picIn = picIn
	//defer func() {
	//// Cleanup if intialization fail
	//if err != nil {
	//x264c.PictureClean(&picIn)
	//}
	//}()

	e.e = x264c.EncoderOpen(&param)
	if e.e == nil {
		err = fmt.Errorf("x264: cannot open the encoder")
		return
	}

	ret := x264c.EncoderHeaders(e.e, e.nals, &e.nnals)
	if ret < 0 {
		err = fmt.Errorf("x264: cannot encode headers")
		return
	}

	if ret > 0 {
		b := C.GoBytes(e.nals[0].PPayload, C.int(ret))
		n, er := e.w.Write(b)
		if er != nil {
			err = er
			return
		}

		if int(ret) != n {
			err = fmt.Errorf("x264: error writing headers, size=%d, n=%d", ret, n)
		}
	}

	return
}

// Encode encodes image.
func (e *Encoder) Encode(im image.Image) (err error) {
	var picOut x264c.Picture

	e.img.ToYCbCr(im)

	picIn := e.picIn

	picIn.Img.ICsp = e.csp

	picIn.Img.IPlane = 3
	picIn.Img.IStride[0] = int32(e.opts.Width)
	picIn.Img.IStride[1] = int32(e.opts.Width) / 2
	picIn.Img.IStride[2] = int32(e.opts.Width) / 2

	y, cb, cr := C.CBytes(e.img.Y), C.CBytes(e.img.Cb), C.CBytes(e.img.Cr)
	picIn.Img.Plane[0] = y
	picIn.Img.Plane[1] = cb
	picIn.Img.Plane[2] = cr

	//e.img.CopyToCPointer(picIn.Img.Plane[0], picIn.Img.Plane[1], picIn.Img.Plane[2])

	picIn.IPts = e.pts
	e.pts += e.tpf

	log.Printf("pts: %v", e.pts)

	ret := x264c.EncoderEncode(e.e, e.nals, &e.nnals, &picIn, &picOut)
	C.free(y)
	C.free(cb)
	C.free(cr)
	if ret < 0 {
		err = fmt.Errorf("x264: cannot encode picture")
		return
	}

	if ret > 0 {
		b := C.GoBytes(e.nals[0].PPayload, C.int(ret))

		n, er := e.w.Write(b)
		if er != nil {
			err = er
			return
		}

		if int(ret) != n {
			err = fmt.Errorf("x264: error writing payload, size=%d, n=%d", ret, n)
		}
	}

	return
}

// Flush flushes encoder.
func (e *Encoder) Flush() (err error) {
	var picOut x264c.Picture

	for x264c.EncoderDelayedFrames(e.e) > 0 {
		ret := x264c.EncoderEncode(e.e, e.nals, &e.nnals, nil, &picOut)
		if ret < 0 {
			err = fmt.Errorf("x264: cannot encode picture")
			return
		}

		if ret > 0 {
			b := C.GoBytes(e.nals[0].PPayload, C.int(ret))

			n, er := e.w.Write(b)
			if er != nil {
				err = er
				return
			}

			if int(ret) != n {
				err = fmt.Errorf("x264: error writing payload, size=%d, n=%d", ret, n)
			}
		}
	}

	return
}

// Close closes encoder.
func (e *Encoder) Close() error {
	picIn := e.picIn
	x264c.PictureClean(&picIn)
	x264c.EncoderClose(e.e)
	return nil
}
