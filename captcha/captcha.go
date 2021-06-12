package captcha

import (
	"bytes"
	"encoding/json"
	"image"
	"sync"
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

	b1, b2 *bytes.Buffer
}

func (c *Captcha) Bytes() []byte {
	return c.renderedData
}

type Generator interface {
	Generate() (*Captcha, error)
	Release(*Captcha) error
}

type staticGenerator struct {
	pool *sync.Pool
}

var _ Generator = &staticGenerator{}

func NewStaticGenerator() *staticGenerator {
	return &staticGenerator{
		pool: &sync.Pool{
			New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 20*1024)) },
		},
	}
}

func (s *staticGenerator) Release(c *Captcha) error {
	s.pool.Put(c.b1)
	s.pool.Put(c.b2)
	return nil
}

func (s *staticGenerator) Generate() (*Captcha, error) {
	const width, height, circlesCount = 500, 200, 10

	img := image.NewGray(image.Rect(0, 0, width, height))

	solX, solY, solW := drawCaptcha(img, width, height, circlesCount)

	buf := s.pool.Get().(*bytes.Buffer)
	buf.Reset()

	enc := NewPngEncoder()

	if err := enc.Encode(buf, img); err != nil {
		return nil, err
	}

	cpt := &Captcha{Data: buf.Bytes(), Solution: Solution{X: solX, Y: solY, W: solW, H: solW}}

	buf2 := s.pool.Get().(*bytes.Buffer)
	buf2.Reset()

	cpt.b1 = buf
	cpt.b2 = buf2

	jsonEnc := json.NewEncoder(buf2)

	if err := jsonEnc.Encode(&cpt); err != nil {
		return nil, err
	}
	cpt.renderedData = buf2.Bytes()

	return cpt, nil
}
