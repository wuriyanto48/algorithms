package main

import (
	"io"
	"strings"
)

// Cartesian represent cartesian plane
type Cartesian struct {
	plane  [][]string
	width  int
	height int
	s      string
}

// NewCartesianPlane the Cartesian constructor
func NewCartesianPlane(s string, width, height int) *Cartesian {
	plane := make([][]string, height+1)
	for y := 0; y < height; y++ {
		plane[y] = make([]string, width)
	}

	return &Cartesian{
		plane:  plane,
		width:  width,
		height: height,
		s:      s,
	}
}

// Draw will draw the plane
func (c *Cartesian) Draw() {
	for y := c.height; y >= 0; y-- {
		var xCoord []string
		for x := 0; x <= c.width; x++ {
			xCoord = append(xCoord, c.s)
		}
		c.plane[y] = xCoord
	}

}

// Plot will draw w to the point p
func (c *Cartesian) Plot(w string, p *Point) {
	c.plane[int(p.Y)][int(p.X)] = w
}

// WriteTo will write cartesian plane to the dst
func (c *Cartesian) WriteTo(dst io.Writer) (int64, error) {
	var written int64 = 0
	for i := len(c.plane) - 1; i >= 0; i-- {
		ss := strings.Join(c.plane[i], " ")

		w, err := dst.Write([]byte(ss))
		if err != nil {
			return written, err
		}

		written = written + int64(w)
		w, err = dst.Write([]byte{0xD, 0xA})
		if err != nil {
			return written, err
		}

		written = written + int64(w)
	}

	return written, nil
}

// Plane will return plane
func (c *Cartesian) Plane() [][]string {
	return c.plane
}
