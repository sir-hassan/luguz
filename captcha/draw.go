package captcha

import (
	"image"
	"image/color"
	"image/draw"
	"math/bits"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func drawCircles(img draw.Image, height int, width int) {
	for i := 0; i < 10; i++ {
		x := rand.Intn(255) + 50
		y := rand.Intn(100) + 50
		r := rand.Intn(20) + 10
		drawCircleInt(img, x, y, r, color.Gray{Y: uint8(255)}, i == 0)
	}
}

func drawCircle(img draw.Image, x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(x0+x, y0+y, c)
		img.Set(x0+y, y0+x, c)
		img.Set(x0-y, y0+x, c)
		img.Set(x0-x, y0+y, c)
		img.Set(x0-x, y0-y, c)
		img.Set(x0-y, y0-x, c)
		img.Set(x0+y, y0-x, c)
		img.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}

func drawCircleInt(img draw.Image, x0, y0, r int, c color.Color, open bool) {
	x := 0
	y := r
	xsq := 0
	rsq := r * r
	ysq := rsq
	// Loop x from 0 to the line x==y. Start y at r and each time
	// around the loop either keep it the same or decrement it.

	i := rand.Intn(8-1+1) + 1

	for x <= y {
		if !open || i != 1 {
			img.Set(x+x0, y+y0, c)

		}
		if !open || i != 2 {
			img.Set(y+x0, x+y0, c)

		}
		if !open || i != 3 {
			img.Set(-x+x0, y+y0, c)

		}
		if !open || i != 4 {
			img.Set(-y+x0, x+y0, c)

		}
		if !open || i != 5 {
			img.Set(x+x0, -y+y0, c)

		}
		if !open || i != 6 {
			img.Set(y+x0, -x+y0, c)

		}
		if !open || i != 7 {
			img.Set(-x+x0, -y+y0, c)

		}
		if !open || i != 8 {
			img.Set(-y+x0, -x+y0, c)
		}

		// New x^2 = (x+1)^2 = x^2 + 2x + 1
		xsq = xsq + 2*x + 1
		x++
		// Potential new y^2 = (y-1)^2 = y^2 - 2y + 1
		y1sq := ysq - 2*y + 1
		// Choose y or y-1, whichever gives smallest error
		a := xsq + ysq
		b := xsq + y1sq
		if a-rsq >= rsq-b {
			y--
			ysq = y1sq
		}
	}
}

func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("image: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}

func mul3NonNeg(x int, y int, z int) int {
	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
}
