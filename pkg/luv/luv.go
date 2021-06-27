package luv

import (
	"github.com/chrisfenner/bytecolor/pkg/cylinder"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	// Chosen experimentally. Gives an even balance of R/G/B.
	hueShift = float64(10)
	// Chosen experimentally. Just as saturated as I can stand.
	chroma = float64(0.065)
	// Chosen experimentally. Avoids too many white shades.
	lightness = float64(0.9 / 8)
)

func New() (*cylinder.Palette, error) {
	return cylinder.NewPalette(
		hueShift,
		chroma,
		lightness,
		func(h, c, l float64) colorful.Color {
			return colorful.LuvLCh(l, c, h)
		},
		func(c1, c2 colorful.Color) float64 {
			return c1.DistanceLuv(c2)
		},
	)
}
