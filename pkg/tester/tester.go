package tester

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"math"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
	tc "github.com/wayneashleyberry/truecolor/pkg/color"
)

var (
	tooCloseSquared = float64(8000)
)

type rgb = [3]byte

type Palette interface {
	// Select returns 8bpc R,G,B values for a given byte value
	Select(b byte) rgb
	// Nearest returns the byte value corresponding to the approximate color
	Nearest(c color.Color) byte
}

func nearestColor(p Palette, c color.Color) rgb {
	return p.Select(p.Nearest(c))
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
	if err := numericOrder(p); err != nil {
		return err
	}
	if err := hueOrder(p); err != nil {
		return err
	}
	if err := lightnessOrder(p); err != nil {
		return err
	}
	if err := grayCodeFill(p); err != nil {
		return err
	}
	if err := hslGamut(p); err != nil {
		return err
	}
	if err := ones(p); err != nil {
		return err
	}
	return nil
}

func ones(p Palette) error {
	x, err := terminal.Width()
	if err != nil {
		return err
	}
	vals := make([]byte, 256)
	for i := range vals {
		vals[i] = byte(i)
	}
	sort.Slice(vals, func(i, j int) bool {
		i1s := countOnes(vals[i])
		i2s := countOnes(vals[j])
		if i1s != i2s {
			return i1s < i2s
		}
		return vals[i] < vals[j]
	})
	rows := make([][]byte, 9)
	for _, val := range vals {
		ones := countOnes(val)
		rows[ones] = append(rows[ones], val)
	}
	for _, row := range rows {
		for i, val := range row {
			if i%int(x) == 0 {
				fmt.Printf("\n")
			}
			bg := p.Select(val)
			tc.Background(byte(bg[0]), byte(bg[1]), byte(bg[2])).Print(" ")
		}
	}
	fmt.Printf("\n")
	return nil
}

func countOnes(b byte) int {
	ones := 0
	for i := 1; i < 256; i <<= 1 {
		if b&byte(i) != 0 {
			ones++
		}
	}
	return ones
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
			tc.Color(fg[0], fg[1], fg[2]).Background(bg[0], bg[1], bg[2]).Print(msg)
		}
	}
	fmt.Printf("\n")
	return nil
}

func hslGamut(p Palette) error {
	// Get the current console width and height for tiling.
	x, err := terminal.Width()
	if err != nil {
		return err
	}
	y, err := terminal.Height()
	if err != nil {
		return err
	}
	if x < 16 || y < 16 {
		return fmt.Errorf("detected terminal size (%d,%d) not big enough for test", x, y)
	}
	if y > 16 {
		y = 16
	}

	// Grayscale across x
	fmt.Printf("\n")
	for j := uint(0); j < x; j++ {
		l := 1.0 / float64(x) * float64(j)
		c := colorful.Hsl(0.0, 0.0, l)
		bg := nearestColor(p, c)
		tc.Background(bg[0], bg[1], bg[2]).Print(" ")
	}
	// Draw an HSL rectangle
	for i := uint(0); i <= y; i++ {
		fmt.Printf("\n")
		for j := uint(0); j < x; j++ {
			// Every x is a step around the hue circle
			// Every y is a step in the lightness
			// Saturation = 1.00
			h := 360.0 / float64(x) * float64(j)
			l := 1.0 / float64(y) * float64(i)
			c := colorful.Hsl(h, 1.0, l)
			bg := nearestColor(p, c)
			tc.Background(bg[0], bg[1], bg[2]).Print(" ")
		}
	}
	fmt.Printf("\n")
	return nil
}

func numericOrder(p Palette) error {
	x, err := terminal.Width()
	if err != nil {
		return err
	}
	for i := 0; i < 256; i++ {
		if i%int(x) == 0 {
			fmt.Printf("\n")
		}
		bg := p.Select(byte(i))
		tc.Background(bg[0], bg[1], bg[2]).Print(" ")
	}
	fmt.Printf("\n")
	return nil
}

func hueOrder(p Palette) error {
	x, err := terminal.Width()
	if err != nil {
		return err
	}
	colors := make([]colorful.Color, 256)
	for i := range colors {
		rgb := p.Select(byte(i))
		colors[i], _ = colorful.MakeColor(color.RGBA{rgb[0], rgb[1], rgb[2], 255})
	}
	sort.Slice(colors, func(i, j int) bool {
		hi, ci, li := colors[i].Hcl()
		hj, cj, lj := colors[j].Hcl()
		const minC = 0.1
		// If one color is very un-colorful, put it first.
		if ci < minC && cj >= minC {
			return true
		}
		if ci >= minC && cj < minC {
			return false
		}
		// If both colors are very un-colorful, order by lightness.
		if ci < minC && cj < minC {
			return li < lj
		}
		// If neither color is un-colorful, order by hue.
		return hi < hj
	})
	for i, c := range colors {
		if i%int(x) == 0 {
			fmt.Printf("\n")
		}
		r, g, b, _ := c.RGBA()
		tc.Background(byte(r/256), byte(g/256), byte(b/256)).Print(" ")
	}
	fmt.Printf("\n")
	return nil
}

func lightnessOrder(p Palette) error {
	x, err := terminal.Width()
	if err != nil {
		return err
	}
	colors := make([]colorful.Color, 256)
	for i := range colors {
		rgb := p.Select(byte(i))
		colors[i], _ = colorful.MakeColor(color.RGBA{rgb[0], rgb[1], rgb[2], 255})
	}
	sort.Slice(colors, func(i, j int) bool {
		hi, ci, li := colors[i].Hcl()
		hj, cj, lj := colors[j].Hcl()
		const minC = 0.1
		// If one color is very un-colorful, put it first.
		if ci < minC && cj >= minC {
			return true
		}
		if ci >= minC && cj < minC {
			return false
		}
		// If both colors are very un-colorful, order by lightness (reversed).
		if ci < minC && cj < minC {
			return li > lj
		}
		const lThresh = 1.0 / 8
		// If both colors have different lightness, compare them by lightness.
		if math.Abs(li-lj) > lThresh {
			return li < lj
		}
		// If both colors have very close lightness, order by hue.
		return hi < hj
	})
	for i, c := range colors {
		if i%int(x) == 0 {
			fmt.Printf("\n")
		}
		r, g, b, _ := c.RGBA()
		tc.Background(byte(r), byte(g/256), byte(b/256)).Print(" ")
	}
	fmt.Printf("\n")
	return nil
}
