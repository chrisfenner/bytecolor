package polar

import "math"

type Coord struct {
	Degrees float64
	Radius  float64
}

func Add(coords ...Coord) Coord {
	resultX := float64(0)
	resultY := float64(0)

	for _, coord := range coords {
		rads := coord.Degrees * math.Pi / 180.0
		resultX += (coord.Radius * math.Cos(rads))
		resultY += (coord.Radius * math.Sin(rads))
	}

	rads := math.Atan2(resultY, resultX)
	resultR := math.Sqrt(resultY*resultY + resultX*resultX)
	// Normalize the angle to [0, 360)
	degrees := rads * 180.0 / math.Pi
	for degrees < 0.0 {
		degrees += 360.0
	}
	for degrees >= 360.0 {
		degrees -= 360.0
	}
	return Coord{
		Degrees: degrees,
		Radius:  resultR,
	}
}
