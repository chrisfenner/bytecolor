package hsl

import (
	"github.com/chrisfenner/bytecolor/pkg/cylinder"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	hueShift = float64(0)
	// this saturation value is calculated such that there is no clamping even for 0x0f
	saturation = float64(0.384)
	lightness  = float64(1.0 / 12)
)

func New() (*cylinder.Palette, error) {
	return cylinder.NewPalette(
		hueShift,
		saturation,
		lightness,
		colorful.Hsl,
		func(c1, c2 colorful.Color) float64 {
			return c1.DistanceRgb(c2)
		},
		map[byte][3]byte{
			255: [3]byte{255, 255, 255},
		},
	)
}
