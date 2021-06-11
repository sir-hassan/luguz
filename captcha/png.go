package captcha

import (
	"image/png"
	"sync"
)

type buffersPool struct {
	sync.Pool
}

func (p *buffersPool) Get() *png.EncoderBuffer {
	return p.Pool.Get().(*png.EncoderBuffer)
}
func (p *buffersPool) Put(buffer *png.EncoderBuffer) {
	p.Pool.Put(buffer)
}

var _ png.EncoderBufferPool = &buffersPool{}

func NewPngEncoder() *png.Encoder {
	return &png.Encoder{
		CompressionLevel: png.DefaultCompression,
		BufferPool: &buffersPool{sync.Pool{
			New: func() interface{} { return new(png.EncoderBuffer) },
		}},
	}
}
