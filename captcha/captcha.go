package captcha

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
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

type BPool struct {
	pool *sync.Pool
}

func (pool *BPool) Get() *png.EncoderBuffer {
	return pool.pool.Get().(*png.EncoderBuffer)
}

func (pool *BPool) Put(buffer *png.EncoderBuffer) {
	pool.pool.Put(buffer)
}

var _ png.EncoderBufferPool = &BPool{}

func (s *staticGenerator) Generate() (*Captcha, error) {
	const width, height = 500, 200

	img := image.NewGray(image.Rect(0, 0, width, height))

	drawCircles(img, width, height)

	buf := bytes.NewBuffer(make([]byte, 0, 4*1024))

	enc := &png.Encoder{
		CompressionLevel: png.DefaultCompression,
		BufferPool: &BPool{&sync.Pool{
			New: func() interface{} { return new(png.EncoderBuffer) },
		}},
	}

	if err := enc.Encode(buf, img); err != nil {
		return nil, err
	}

	buf2 := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(buf2, buf.Bytes())

	cpt := &Captcha{Data: buf2, Solution: Solution{X: 20, Y: 50, W: 3, H: 4}}

	data, err := json.Marshal(&cpt)
	if err != nil {
		return nil, err
	}
	cpt.renderedData = data

	return cpt, nil
}
