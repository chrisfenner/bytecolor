package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	igif "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"
	"strings"

	"github.com/chrisfenner/bytecolor/pkg/gif"
	"github.com/chrisfenner/bytecolor/pkg/hcl"
	"github.com/chrisfenner/bytecolor/pkg/hsl"
	"github.com/chrisfenner/bytecolor/pkg/hsv"
	"github.com/chrisfenner/bytecolor/pkg/luv"
	"github.com/chrisfenner/bytecolor/pkg/tester"
)

var (
	palette = flag.String("palette", "hsv", "which color palette to use")
	in      = flag.String("in", "", "the path of the input file(s) (comma-separated)")
)

func main() {
	retval := 0
	err := mainWithError()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		retval = -1
	}
	os.Exit(retval)
}

func mainWithError() error {
	flag.Parse()
	if *in == "" {
		return fmt.Errorf("please provide input file or files (comma-separated)")
	}
	infiles := strings.Split(*in, ",")
	if len(infiles) == 0 {
		return fmt.Errorf("please provide at least one input file (comma-separated")
	}

	var pal tester.Palette
	var err error
	switch strings.ToLower(*palette) {
	case "hsl":
		pal, err = hsl.New()
		if err != nil {
			return err
		}
	case "hsv":
		pal, err = hsv.New()
		if err != nil {
			return err
		}
	case "hcl":
		pal, err = hcl.New()
		if err != nil {
			return err
		}
	case "luv":
		pal, err = luv.New()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported palette '%s'", *palette)
	}

	gifs := make([]*image.Paletted, len(infiles))
	var bounds *image.Rectangle
	var sb strings.Builder
	for i := range infiles {
		// Read the image
		sb.WriteString(strings.Split(path.Base(infiles[i]), ".")[0] + "-")
		imageFile, err := os.Open(infiles[i])
		if err != nil {
			return err
		}
		defer imageFile.Close()
		m, _, err := image.Decode(imageFile)
		if err != nil {
			return err
		}

		// Encode it as a gif to a buffer using library dithering
		var buf bytes.Buffer
		if err := gif.Encode(&buf, pal, m); err != nil {
			return err
		}
		// Decode the gif to get the paletted image
		g, err := igif.DecodeAll(&buf)
		if err != nil {
			return err
		}
		// Make sure it has the same bounds as all the others
		if bounds == nil {
			bounds = &g.Image[0].Rect
		} else {
			if g.Image[0].Rect != *bounds {
				return fmt.Errorf("when passing multiple images, please make sure they are all the same size")
			}
		}
		gifs[i] = g.Image[0]
	}
	sb.WriteString(strings.ToLower(*palette))
	sb.WriteString(".gif")

	// XOR all the images together if there are more than 1
	for i := 1; i < len(gifs); i++ {
		for j := 0; j < len(gifs[0].Pix); j++ {
			gifs[0].Pix[j] ^= gifs[i].Pix[j]
		}
	}

	outfile := sb.String()
	w, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer w.Close()

	result := igif.GIF{
		Image: gifs[:1],
		Delay: make([]int, 1),
	}
	if err := igif.EncodeAll(w, &result); err != nil {
		return err
	}

	fmt.Printf("converted/XORed %d images to %s-256 palette as a GIF in %s.\n", len(infiles), strings.ToLower(*palette), outfile)

	return nil
}
