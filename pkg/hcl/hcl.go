package hcl

import (
	"github.com/chrisfenner/bytecolor/pkg/cylinder"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	// Chosen by experimentation: Gives a palette that seems balanced.
	hueShift = float64(30)
	// Chosen by experimentation: Just barely not too saturated.
	chroma = float64(0.065)
	// Chosen by experimentation: Avoids too many "nearly white" shades.
	lightness = float64(0.9 / 8)
)

func New() (*cylinder.Palette, error) {
	return cylinder.NewPalette(
		hueShift,
		chroma,
		lightness,
		colorful.Hcl,
		func(c1, c2 colorful.Color) float64 {
			return c1.DistanceLab(c2)
		},
	)
}
