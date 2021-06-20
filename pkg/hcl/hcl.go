package hcl

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

type HCLPalette struct {
	hueShift         float64
	chroma           float64
	luminanceOfWhite float64
}

func New(hueShift, chroma, luminanceOfWhite float64) (*HCLPalette, error) {
	if hueShift > 45.0 || hueShift < 0.0 {
		return nil, fmt.Errorf("hueShift must be between 0 and 45 degrees")
	}
	if chroma > 1.0 || chroma < 0.0 {
		return nil, fmt.Errorf("chroma must be between 0 and 1.0")
	}
	if luminanceOfWhite > 1.0 || luminanceOfWhite < 0.0 {
		return nil, fmt.Errorf("luminanceOfWhite must be between 0 and 1.0")
	}
	return &HCLPalette{
		hueShift:         hueShift,
		chroma:           chroma,
		luminanceOfWhite: luminanceOfWhite,
	}, nil
}

func (p *HCLPalette) Select(val byte) [3]byte {
	toBlend := p.blendList(val)
	if len(toBlend) == 0 {
		return [3]byte{0, 0, 0}
	}
	c := toBlend[0]
	toBlend = toBlend[1:]
	for len(toBlend) > 0 {
		c = c.BlendHcl(toBlend[0], 0.5)
		toBlend = toBlend[1:]
	}
	r, g, b, _ := c.Clamped().RGBA()
	return [3]byte{byte(r / 256), byte(g / 256), byte(b / 256)}
}

func (p *HCLPalette) blendList(b byte) []colorful.Color {
	var result []colorful.Color
	for i := 0; i < 8; i++ {
		if b&(1<<i) != 0 {
			result = append(result, colorful.Hcl(
				p.hueShift+float64(i)*360.0/8.0, p.chroma, p.luminanceOfWhite/8.0))
		}
	}
	return result
}
