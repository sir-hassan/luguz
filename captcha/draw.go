package captcha

import (
	"image/color"
	"image/draw"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func drawCaptcha(img draw.Image, width int, height int, circlesNum int) (solX, solY, solW int) {
	for i := 0; i < circlesNum; i++ {
		minDimension := height
		if height > width {
			minDimension = width
		}
		maxRadius := (minDimension)/2 - 1
		minRadius := minDimension / 20

		r := randBetween(minRadius, maxRadius)
		x := randBetween(r, width-r-1)
		y := randBetween(r, height-r-1)
		drawCircle(img, x, y, r, color.Gray{Y: uint8(255)}, i == 0)
		if i == 0 {
			solX = x - r
			solY = y - r
			solW = 2 * r
			drawSquare(img, x-r, y-r, 2*r, color.Gray{Y: uint8(255)})
		}
	}
	return solX, solY, solW
}

func drawCircle(img draw.Image, x0, y0, r int, c color.Color, open bool) {
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
func drawSquare(img draw.Image, x, y, w int, c color.Color) {

	for i := x; i <= x+w; i++ {
		img.Set(i, y, c)
		img.Set(i, y+w, c)
	}
	for i := y; i <= y+w; i++ {
		img.Set(x, i, c)
		img.Set(x+w, i, c)
	}
}

func randBetween(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
