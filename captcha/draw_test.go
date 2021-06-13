package captcha

import (
	"image"
	"image/color"
	"testing"
)

func TestDrawCircle(t *testing.T) {
	tt := []struct {
		name   string
		width  int
		height int
		x0     int
		y0     int
		r      int
	}{
		{
			name:   "size 120x120",
			width:  120,
			height: 120,
			x0:     60,
			y0:     60,
			r:      50,
		},
		{
			name:   "size 100x50",
			width:  100,
			height: 50,
			x0:     50,
			y0:     25,
			r:      15,
		},
		{
			name:   "size 100x100",
			width:  100,
			height: 100,
			x0:     40,
			y0:     40,
			r:      10,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			width, height := tc.width, tc.height
			img := image.NewGray(image.Rect(0, 0, width, height))

			x0, y0, r := tc.x0, tc.y0, tc.r
			drawCircle(img, x0, y0, r, color.Gray{Y: uint8(255)}, false)

			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					val := img.Pix[img.PixOffset(x, y)]
					if val == 255 {
						rr := (x0-x)*(x0-x) + (y0-y)*(y0-y)
						if !isBetween(rr, r*(r-1), r*(r+1)) {
							t.Errorf("unexpected point distance from origin r2:%d", rr)
						}
					}
				}
			}
		})
	}
}

func isBetween(x, a, b int) bool {
	return x >= a && x <= b
}
