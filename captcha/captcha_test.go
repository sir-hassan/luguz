package captcha

import "testing"

func BenchmarkCaptchaGenerator(b *testing.B) {
	cfg := GeneratorConfig{
		Width:   100,
		Height:  100,
		Circles: 10,
	}
	generator := NewBasicGenerator(cfg)

	for i := 0; i < b.N; i++ {
		generator.Generate()
	}
}
