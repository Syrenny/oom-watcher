package app

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

const (
	panelIconHeight = 22
	panelIconPadX   = 3
	panelIconPadY   = 3
	segmentWidth    = 6
	segmentHeight   = 12
	segmentThick    = 2
	charGap         = 2
	percentWidth    = 8
)

var segmentMap = map[rune][7]bool{
	'0': {true, true, true, true, true, true, false},
	'1': {false, true, true, false, false, false, false},
	'2': {true, true, false, true, true, false, true},
	'3': {true, true, true, true, false, false, true},
	'4': {false, true, true, false, false, true, true},
	'5': {true, false, true, true, false, true, true},
	'6': {true, false, true, true, true, true, true},
	'7': {true, true, true, false, false, false, false},
	'8': {true, true, true, true, true, true, true},
	'9': {true, true, true, true, false, true, true},
	'-': {false, false, false, false, false, false, true},
}

type panelIcons struct {
	Visible []byte
	Blank   []byte
}

func renderPanelIcons(text string) (panelIcons, error) {
	visible, err := renderPanelIcon(text, true)
	if err != nil {
		return panelIcons{}, err
	}

	blank, err := renderPanelIcon(text, false)
	if err != nil {
		return panelIcons{}, err
	}

	return panelIcons{
		Visible: visible,
		Blank:   blank,
	}, nil
}

func renderPanelIcon(text string, visible bool) ([]byte, error) {
	width := panelIconWidth(text)
	img := image.NewNRGBA(image.Rect(0, 0, width, panelIconHeight))
	if visible {
		drawPanelText(img, text, color.NRGBA{R: 245, G: 247, B: 250, A: 255})
	}
	return encodePNG(img)
}

func panelIconWidth(text string) int {
	width := panelIconPadX * 2
	for _, ch := range text {
		width += glyphWidth(ch)
	}
	if len(text) > 1 {
		width += (len(text) - 1) * charGap
	}
	return width
}

func drawPanelText(img *image.NRGBA, text string, c color.NRGBA) {
	x := panelIconPadX
	for _, ch := range text {
		drawGlyph(img, x, panelIconPadY, ch, c)
		x += glyphWidth(ch) + charGap
	}
}

func drawGlyph(img *image.NRGBA, x, y int, ch rune, c color.NRGBA) {
	switch ch {
	case '%':
		drawPercent(img, x, y, c)
	default:
		drawSevenSegment(img, x, y, ch, c)
	}
}

func drawSevenSegment(img *image.NRGBA, x, y int, ch rune, c color.NRGBA) {
	segments, ok := segmentMap[ch]
	if !ok {
		segments = segmentMap['-']
	}

	if segments[0] {
		fillRect(img, x+1, y, x+1+segmentWidth, y+segmentThick, c)
	}
	if segments[1] {
		fillRect(img, x+segmentWidth, y+1, x+segmentWidth+segmentThick, y+1+segmentHeight/2, c)
	}
	if segments[2] {
		fillRect(img, x+segmentWidth, y+segmentHeight/2+1, x+segmentWidth+segmentThick, y+segmentHeight+1, c)
	}
	if segments[3] {
		fillRect(img, x+1, y+segmentHeight, x+1+segmentWidth, y+segmentHeight+segmentThick, c)
	}
	if segments[4] {
		fillRect(img, x, y+segmentHeight/2+1, x+segmentThick, y+segmentHeight+1, c)
	}
	if segments[5] {
		fillRect(img, x, y+1, x+segmentThick, y+1+segmentHeight/2, c)
	}
	if segments[6] {
		fillRect(img, x+1, y+segmentHeight/2, x+1+segmentWidth, y+segmentHeight/2+segmentThick, c)
	}
}

func drawPercent(img *image.NRGBA, x, y int, c color.NRGBA) {
	fillRect(img, x, y+1, x+2, y+3, c)
	fillRect(img, x+percentWidth-2, y+segmentHeight-1, x+percentWidth, y+segmentHeight+1, c)
	for offset := 0; offset < 6; offset++ {
		fillRect(img, x+2+offset, y+segmentHeight-1-offset, x+3+offset, y+segmentHeight-offset, c)
	}
}

func glyphWidth(ch rune) int {
	if ch == '%' {
		return percentWidth
	}

	return segmentWidth + segmentThick
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
