package hsv

import (
	"github.com/chrisfenner/bytecolor/pkg/cylinder"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	hueShift   = float64(0)
	saturation = float64(0.5)
	value      = float64(1.0 / 8)
)

func New() (*cylinder.Palette, error) {
	return cylinder.NewPalette(
		hueShift,
		saturation,
		value,
		Hsv,
		func(c1, c2 colorful.Color) float64 {
			return c1.DistanceRgb(c2)
		},
	)
}
