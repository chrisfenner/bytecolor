package luv

import (
	"github.com/chrisfenner/bytecolor/pkg/cylinder"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	hueShift  = float64(30)
	chroma    = float64(0.05)
	lightness = float64(1.0 / 8)
)

func New() (*cylinder.Palette, error) {
	return cylinder.NewPalette(
		hueShift,
		chroma,
		lightness,
		func(h, c, l float64) colorful.Color {
			return colorful.LuvLCh(l, c, h)
		},
	)
}
