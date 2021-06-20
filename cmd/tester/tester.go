package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chrisfenner/bytecolor/pkg/hcl"
	"github.com/chrisfenner/bytecolor/pkg/hsl"
	"github.com/chrisfenner/bytecolor/pkg/hsv"
	"github.com/chrisfenner/bytecolor/pkg/tester"
)

var (
	palette = flag.String("palette", "hsv", "which color palette to test")
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
	var pal tester.Palette
	var err error
	switch strings.ToLower(*palette) {
	case "hsl":
		pal, err = hsl.New(30, 1.0, 0.98)
		if err != nil {
			return err
		}
	case "hsv":
		pal, err = hsv.New(30, 1.0, 0.98)
		if err != nil {
			return err
		}
	case "hcl":
		pal, err = hcl.New(30, 1.0, 0.95)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported palette '%s'", *palette)
	}

	if err := tester.Test(pal); err != nil {
		return err
	}

	return nil
}
