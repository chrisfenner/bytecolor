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
	"github.com/chrisfenner/bytecolor/pkg/windows"
)

var (
	palette = flag.String("palette", "hsv", "which color palette to use")
	in      = flag.String("in", "", "the path of the input file(s) (comma-separated)")
	animate = flag.String("animate", "", "for 2-image merges, whether to animate (vertical or horizontal)")
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
	case "win":
		pal, err = windows.New()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported palette '%s'", *palette)
	}

	gifs := make([]*image.Paletted, len(infiles))
	var bounds *image.Rectangle
	var config *image.Config
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
		// Also, save the config for re-use
		if bounds == nil {
			bounds = &g.Image[0].Rect
			config = &g.Config
		} else {
			if g.Image[0].Rect != *bounds {
				return fmt.Errorf("when passing multiple images, please make sure they are all the same size")
			}
		}
		gifs[i] = g.Image[0]
	}
	sb.WriteString(strings.ToLower(*palette))
	sb.WriteString(".gif")
	gifTemplate := gifs[0]
	const delay = 5 // 20fps
	result := igif.GIF{
		Config: *config,
	}

	if *animate != "" {
		if len(gifs) != 2 {
			return fmt.Errorf("'animate' option requires 2 images")
		}
		offset := 0
		numFrames := 0
		switch *animate {
		case "vertical":
			offset = gifTemplate.Stride * 4
			numFrames = gifTemplate.Rect.Dy() / 4
		case "horizontal":
			offset = 4
			numFrames = gifTemplate.Rect.Dx() / 4
		default:
			return fmt.Errorf("unrecognized animation option '%s', only 'vertical' or 'horizontal' are supported", *animate)
		}
		for i := 0; i < numFrames; i++ {
			outData := xorWithOffset(gifs[0].Pix, gifs[1].Pix, offset*i)
			result.Image = append(result.Image, frame(gifTemplate, outData))
			result.Delay = append(result.Delay, delay)
		}
	} else {
		inDatas := make([][]byte, len(gifs))
		for i := range inDatas {
			inDatas[i] = gifs[i].Pix
		}
		outData := xorAll(inDatas)
		result.Image = append(result.Image, frame(gifTemplate, outData))
		result.Delay = append(result.Delay, delay)
	}

	outfile := sb.String()
	w, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := igif.EncodeAll(w, &result); err != nil {
		return err
	}

	fmt.Printf("converted/XORed %d images to %s-256 palette as a GIF in %s.\n", len(infiles), strings.ToLower(*palette), outfile)

	return nil
}

func frame(template *image.Paletted, data []byte) *image.Paletted {
	result := *template
	result.Pix = data
	return &result
}

func xorWithOffset(a, b []byte, offset int) []byte {
	result := make([]byte, len(a))
	for i := range result {
		//result[i] = a[i] ^ b[i]
		result[i] = a[i] ^ b[(i+offset)%len(a)]
	}
	return result
}

func xorAll(datas [][]byte) []byte {
	result := datas[0]
	for i := 1; i < len(datas); i++ {
		result = xorWithOffset(result, datas[i], 0)
	}
	return result
}
