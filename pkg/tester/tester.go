package tester

import (
	"encoding/hex"
	"fmt"
	"sort"

	terminal "github.com/wayneashleyberry/terminal-dimensions"
	"github.com/wayneashleyberry/truecolor/pkg/color"
)

var (
	tooCloseSquared = float64(8000)
)

type rgb = [3]byte

type Palette interface {
	// Select returns 8bpc R,G,B values for a given byte value
	Select(b byte) rgb
}

func invert(c rgb) rgb {
	res := [3]byte{255 - c[0], 255 - c[1], 255 - c[2]}
	// For certain medium grays, the inverted color is not distinct
	// Print white in these cases
	distance := distanceSquared(c, res)
	if distance < tooCloseSquared {
		res = [3]byte{255, 255, 255}
	}
	return res
}

func distanceSquared(a, b rgb) float64 {
	res := float64(0)

	for i := 0; i < 3; i++ {
		res += (float64(a[i]) - float64(b[i])) * (float64(a[i]) - float64(b[i]))
	}

	return res
}

func grayCode(row, column uint) byte {
	result := byte(0)

	if i := row % 4; i >= 1 && i <= 2 {
		result |= 0x01
	}
	if i := row % 8; i >= 2 && i <= 5 {
		result |= 0x02
	}
	if i := row % 16; i >= 4 && i <= 11 {
		result |= 0x04
	}
	if i := row % 16; i >= 8 && i <= 15 {
		result |= 0x08
	}
	if i := column % 4; i >= 1 && i <= 2 {
		result |= 0x10
	}
	if i := column % 8; i >= 2 && i <= 5 {
		result |= 0x20
	}
	if i := column % 16; i >= 4 && i <= 11 {
		result |= 0x40
	}
	if i := column % 16; i >= 8 && i <= 15 {
		result |= 0x80
	}

	return result
}

func Test(p Palette) error {
	if err := grayCodeFill(p); err != nil {
		return err
	}
	if err := numericOrder(p); err != nil {
		return err
	}
	if err := channelOrder(p); err != nil {
		return err
	}
	return nil
}

func grayCodeFill(p Palette) error {
	// Get the current console width and height for tiling.
	// Each cell will be 2 characters wide, to hold hex values.
	x, err := terminal.Width()
	if err != nil {
		return err
	}
	x /= 2
	y, err := terminal.Height()
	if err != nil {
		return err
	}
	if x < 16 || y < 16 {
		return fmt.Errorf("detected terminal size (%d,%d) not big enough for test", x, y)
	}
	if y > 20 {
		y = 20
	}

	// Center the 16x16 on the screen
	lPad := (x - 16) / 2
	uPad := (y - 16) / 2

	for i := uint(0); i < y; i++ {
		fmt.Printf("\n")
		for j := uint(0); j < x; j++ {
			// Shift the code values left/up by the padding
			code := grayCode((i+16-uPad)%16, (j + 16 - lPad))
			msg := "  "
			// If we are filling in the center 16x16 square, print a value
			if i >= uPad && i < (uPad+16) && j >= lPad && j < (lPad+16) {
				msg = hex.EncodeToString([]byte{code})
			}
			bg := p.Select(code)
			fg := invert(bg)
			color.Color(fg[0], fg[1], fg[2]).Background(bg[0], bg[1], bg[2]).Print(msg)
		}
	}
	return nil
}

func numericOrder(p Palette) error {
	for i := 0; i < 256; i++ {
		bg := p.Select(byte(i))
		color.Background(bg[0], bg[1], bg[2]).Print(" ")
	}
	fmt.Printf("\n")
	return nil
}

func channelOrder(p Palette) error {
	colors := make([][3]byte, 256)
	for i := range colors {
		colors[i] = p.Select(byte(i))
	}
	for channel := 0; channel < 3; channel++ {
		sort.Slice(colors, func(i, j int) bool {
			ic := colors[i]
			jc := colors[j]
			if ic[channel] != jc[channel] {
				return ic[channel] < jc[channel]
			}
			if ic[(channel+1)%3] != jc[(channel+1)%3] {
				return ic[(channel+1)%3] < jc[(channel+1)%3]
			}
			return ic[(channel+2)%3] < jc[(channel+2)%3]
		})
		for _, c := range colors {
			color.Background(c[0], c[1], c[2]).Print(" ")
		}
		fmt.Printf("\n")
	}
	return nil
}
