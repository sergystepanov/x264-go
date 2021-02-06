package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/sergystepanov/x264-go/v2"
)

func main() {
	file, err := os.Create("screen.264")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	bounds := screenshot.GetDisplayBounds(0)

	opts := &x264.Options{
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		FrameRate: 60,
		Tune:      "zerolatency",
		Preset:    "ultrafast",
		Profile:   "baseline",
		LogLevel:  x264.LogDebug,
	}

	enc, err := x264.NewEncoder(file, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	defer enc.Close()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Second / time.Duration(120))

	for range ticker.C {
		select {
		case <-s:
			enc.Flush()

			err = file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			os.Exit(0)
		default:
		}
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			continue
		}

		err = enc.Encode(img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		//enc.Flush()
	}
}
