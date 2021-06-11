package captcha

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
)

type Solution struct {
	X int
	Y int
	W int
	H int
}

type Captcha struct {
	Solution Solution `json:"solution"`
	Data     []byte   `json:"data"`

	renderedData []byte
}

func (c *Captcha) Bytes() []byte {
	return c.renderedData
}

type Generator interface {
	Generate() (*Captcha, error)
	Release(*Captcha) error
}

type staticGenerator struct {
}

var _ Generator = &staticGenerator{}

func NewStaticGenerator() *staticGenerator {
	return &staticGenerator{}
}

func (s *staticGenerator) Release(c *Captcha) error {
	return nil
}

func (s *staticGenerator) Generate() (*Captcha, error) {
	const width, height, circlesCount = 500, 200, 10

	img := image.NewGray(image.Rect(0, 0, width, height))

	solX, solY, solW := drawCaptcha(img, width, height, circlesCount)

	buf := bytes.NewBuffer(make([]byte, 0, 4*1024))

	enc := NewPngEncoder()

	if err := enc.Encode(buf, img); err != nil {
		return nil, err
	}

	buf2 := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(buf2, buf.Bytes())

	cpt := &Captcha{Data: buf2, Solution: Solution{X: solX, Y: solY, W: solW, H: solW}}

	data, err := json.Marshal(&cpt)
	if err != nil {
		return nil, err
	}
	cpt.renderedData = data

	return cpt, nil
}
