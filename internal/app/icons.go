package app

import (
	"bytes"
	"image"
	"image/png"
)

type trayIcons struct {
	Blank []byte
}

func newTrayIcons() (trayIcons, error) {
	blank, err := renderBlankIcon()
	if err != nil {
		return trayIcons{}, err
	}

	return trayIcons{
		Blank: blank,
	}, nil
}

func renderBlankIcon() ([]byte, error) {
	img := image.NewNRGBA(image.Rect(0, 0, 24, 24))
	return encodePNG(img)
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
