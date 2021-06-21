package hsv

import (
	"github.com/chrisfenner/bytecolor/pkg/cylinder"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	hueShift = float64(0)
	// this saturation value is calculated such that there is no clamping even for 0x0f
	saturation = float64(0.384)
	value      = float64(1.0 / 8)
)

func New() (*cylinder.Palette, error) {
	return cylinder.NewPalette(
		hueShift,
		saturation,
		value,
		colorful.Hsv,
	)
}
