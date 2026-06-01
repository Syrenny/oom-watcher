package app

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

type trayIcons struct {
	Normal []byte
	Alert  []byte
	Blank  []byte
}

func newTrayIcons() (trayIcons, error) {
	normal, err := renderChipIcon(color.NRGBA{R: 74, G: 222, B: 128, A: 255}, false)
	if err != nil {
		return trayIcons{}, err
	}

	alert, err := renderChipIcon(color.NRGBA{R: 248, G: 113, B: 113, A: 255}, true)
	if err != nil {
		return trayIcons{}, err
	}

	blank, err := renderBlankIcon()
	if err != nil {
		return trayIcons{}, err
	}

	return trayIcons{
		Normal: normal,
		Alert:  alert,
		Blank:  blank,
	}, nil
}

func renderBlankIcon() ([]byte, error) {
	img := image.NewNRGBA(image.Rect(0, 0, 24, 24))
	return encodePNG(img)
}

func renderChipIcon(accent color.NRGBA, warning bool) ([]byte, error) {
	const size = 24

	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	body := color.NRGBA{R: 225, G: 231, B: 239, A: 255}
	outline := color.NRGBA{R: 15, G: 23, B: 42, A: 255}
	pin := color.NRGBA{R: 148, G: 163, B: 184, A: 255}

	fillRect(img, 7, 5, 17, 19, body)
	fillRect(img, 8, 6, 16, 18, color.NRGBA{R: 30, G: 41, B: 59, A: 255})

	fillRect(img, 7, 5, 17, 6, outline)
	fillRect(img, 7, 18, 17, 19, outline)
	fillRect(img, 7, 5, 8, 19, outline)
	fillRect(img, 16, 5, 17, 19, outline)

	for y := 8; y <= 14; y += 3 {
		fillRect(img, 10, y, 14, y+2, accent)
	}

	for x := 4; x <= 18; x += 4 {
		fillRect(img, x, 3, x+1, 5, pin)
		fillRect(img, x, 19, x+1, 21, pin)
	}

	for y := 8; y <= 14; y += 3 {
		fillRect(img, 5, y, 7, y+1, pin)
		fillRect(img, 17, y, 19, y+1, pin)
	}

	if warning {
		fillRect(img, 14, 7, 17, 15, accent)
		fillRect(img, 14, 16, 17, 18, accent)
	}

	return encodePNG(img)
}

func fillRect(img *image.NRGBA, x0, y0, x1, y1 int, c color.NRGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.SetNRGBA(x, y, c)
		}
	}
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
