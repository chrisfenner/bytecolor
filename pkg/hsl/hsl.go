package hsl

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

type HSVPalette struct {
	bitcolors [8][3]byte
}

func New(hueShift, saturation, luminanceOfWhite float64) (*HSVPalette, error) {
	if hueShift > 45.0 || hueShift < 0.0 {
		return nil, fmt.Errorf("hueShift must be between 0 and 45 degrees")
	}
	if saturation > 1.0 || saturation < 0.0 {
		return nil, fmt.Errorf("saturation must be between 0 and 1.0")
	}
	if luminanceOfWhite > 1.0 || luminanceOfWhite < 0.0 {
		return nil, fmt.Errorf("luminanceOfWhite must be between 0 and 1.0")
	}
	var bitcolors [8][3]byte
	// Divide the 8 bits of the byte into 8 evenly spaced hues with 1/8 luminance and given saturation.
	for i := 0; i < 8; i++ {
		c := colorful.Hsl(hueShift+float64(i)*360.0/8.0, saturation, luminanceOfWhite/8.0)
		fmt.Printf("%v\n", c)
		r, g, b, _ := c.RGBA()
		bitcolors[i] = [3]byte{byte(r / 256), byte(g / 256), byte(b / 256)}
	}
	return &HSVPalette{
		bitcolors: bitcolors,
	}, nil
}

func (p *HSVPalette) Select(b byte) [3]byte {
	result := [3]byte{0, 0, 0}
	for i := 0; i < 8; i++ {
		if b&(1<<i) != 0 {
			result[0] += p.bitcolors[i][0]
			result[1] += p.bitcolors[i][1]
			result[2] += p.bitcolors[i][2]
		}
	}
	return result
}
