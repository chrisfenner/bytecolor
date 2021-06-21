package cylinder

import (
	"fmt"

	"github.com/chrisfenner/bytecolor/pkg/polar"
	"github.com/lucasb-eyer/go-colorful"
)

type ColorModel = func(angle, radius, height float64) colorful.Color

type Palette struct {
	bitcolors [8]colorful.Color
	model     ColorModel
}

const (
	angleShift = float64(0)
	// this baseRadius baseHeight is calculated such that there is no clamping even for 0x0f
	//baseRadius = float64(0.384)
	baseRadius = float64(0.5)
	baseHeight = float64(1.0 / 8)
)

func NewPalette(angleShift, baseRadius, baseHeight float64, model ColorModel) (*Palette, error) {
	if angleShift > 45.0 || angleShift < 0.0 {
		return nil, fmt.Errorf("angleShift must be between 0 and 45 degrees")
	}
	if baseRadius > 1.0 || baseRadius < 0.0 {
		return nil, fmt.Errorf("baseRadius must be between 0 and 1.0")
	}
	if baseHeight > 1.0 || baseHeight < 0.0 {
		return nil, fmt.Errorf("baseHeight must be between 0 and 1.0")
	}
	var bitcolors [8]colorful.Color
	// Divide the 8 bits of the byte into 8 evenly spaced hues with given baseHeight and given baseRadius.
	for i := 0; i < 8; i++ {
		bitcolors[i] = model(angleShift+float64(i)*360.0/8.0, baseRadius, baseHeight)
	}
	return &Palette{
		bitcolors: bitcolors,
		model:     model,
	}, nil
}

func (p *Palette) Select(val byte) [3]byte {
	var mixPolars []polar.Coord
	mixValue := float64(0)
	for i := 0; i < 8; i++ {
		if val&(1<<i) != 0 {
			h, s, v := p.bitcolors[i].Hsv()
			mixValue += v
			mixPolars = append(mixPolars, polar.Coord{
				Degrees: h,
				Radius:  s,
			})
		}
	}
	mix := polar.Add(mixPolars...)
	result := p.model(mix.Degrees, mix.Radius, mixValue).Clamped()
	r, g, b, _ := result.RGBA()
	return [3]byte{byte(r / 256), byte(g / 256), byte(b / 256)}
}
