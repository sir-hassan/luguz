package captcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"sync"
)

// Generator interface represents the an interface that generates captchas :).
type Generator interface {
	// Generate generates a new captcha
	Generate() (*Captcha, error)

	// Release func helps implementations to release captcha resources or
	// buffers to reuse by the underlying logic to reduce heap allocations.
	Release(*Captcha) error
}

// Solution struct represent the captcha solution.
type Solution struct {
	X int
	Y int
	W int
	H int
}

// Captcha Struct represents a rendered instance of a captcha.
type Captcha struct {
	// represents the solution.
	Solution Solution `json:"solution"`
	// represents the image data.
	Data []byte `json:"data"`

	// holds final rendered json data:
	// {"solution":{"X":350,"Y":50,"W":142,"H":142},"data":"BASE_ENCODED_PNG_IMAGE"}
	renderedData []byte

	// b1 and b2 used to reuse the buffers.
	b1, b2 *bytes.Buffer
}

// Bytes fuc returns the final rendered json bytes.
func (c *Captcha) Bytes() []byte {
	return c.renderedData
}

// GeneratorConfig struct holds all the generator configuration.
type GeneratorConfig struct {
	// Height of the png image to generated.
	Height int
	// Width of the png image to generated.
	Width int
	// Circles is how many circles to draw in the png image to generate.
	Circles int
}

type basicGenerator struct {
	config GeneratorConfig
	pool   *sync.Pool
}

var _ Generator = &basicGenerator{}

// NewBasicGenerator func creates a new basicGenerator instance.
func NewBasicGenerator(config GeneratorConfig) Generator {
	return &basicGenerator{
		config: config,
		pool: &sync.Pool{
			New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 20*1024)) },
		},
	}
}

// Release func reuse b1 and b2 buffers to the buffers pool.
func (s *basicGenerator) Release(c *Captcha) error {
	s.pool.Put(c.b1)
	s.pool.Put(c.b2)
	return nil
}

func (s *basicGenerator) Generate() (*Captcha, error) {
	img := image.NewGray(image.Rect(0, 0, s.config.Width, s.config.Height))
	solX, solY, solW := drawCaptcha(img, s.config.Width, s.config.Height, s.config.Circles)

	buf := s.pool.Get().(*bytes.Buffer)
	buf.Reset()

	enc := NewPngEncoder()
	if err := enc.Encode(buf, img); err != nil {
		return nil, fmt.Errorf("encoding png failed, err: %s", err)
	}

	cpt := &Captcha{Data: buf.Bytes(), Solution: Solution{X: solX, Y: solY, W: solW, H: solW}}

	buf2 := s.pool.Get().(*bytes.Buffer)
	buf2.Reset()

	cpt.b1 = buf
	cpt.b2 = buf2

	jsonEnc := json.NewEncoder(buf2)

	if err := jsonEnc.Encode(&cpt); err != nil {
		return nil, fmt.Errorf("encoding json failed, err: %s", err)
	}
	cpt.renderedData = buf2.Bytes()

	return cpt, nil
}
