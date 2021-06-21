package hcl

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
		colorful.Hcl,
	)
}
