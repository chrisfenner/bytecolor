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
	)
}
