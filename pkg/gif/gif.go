package gif

import (
	"image"
	"image/color"
	"image/gif"
	"io"
)

type rgb = [3]byte

type Palette interface {
	// Select returns 8bpc R,G,B values for a given byte value
	Select(b byte) rgb
	// Nearest returns the byte value corresponding to the approximate color
	Nearest(c color.Color) byte
}

func nearestColor(p Palette, c color.Color) rgb {
	return p.Select(p.Nearest(c))
}

type paletteQuantizer struct {
	palette Palette
}

func (p *paletteQuantizer) Quantize(_ color.Palette, _ image.Image) color.Palette {
	// This is a special type of Quantizer that doesn't care about the image.
	// It dutifully reports the 256 colors of the underlying gif.Palette
	// (which is in turn some implementation under bytecolor).
	colors := color.Palette(make([]color.Color, 256))
	for i := 0; i < 256; i++ {
		rgb := p.palette.Select(byte(i))
		colors[i] = color.RGBA{rgb[0], rgb[1], rgb[2], 255}
	}
	return colors
}

func Encode(w io.Writer, p Palette, m image.Image) error {
	q := &paletteQuantizer{p}
	opts := &gif.Options{
		NumColors: 256,
		Quantizer: q,
	}
	return gif.Encode(w, m, opts)
}
