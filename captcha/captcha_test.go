package captcha

import (
	"testing"
)

func Benchmark(b *testing.B) {
	generator := NewStaticGenerator()
	for i := 0; i < b.N; i++ {
		if _, err := generator.Generate(); err != nil {
			b.Errorf("error in generating")
		}
	}
}
