package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/sir-hassan/luguz/captcha"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	tt := []struct {
		name   string
		height int
		width  int
	}{
		{
			name:   "size 10x10",
			height: 10,
			width:  10,
		},
		{
			name:   "size 100x200",
			height: 10,
			width:  10,
		},
		{
			name:   "size 1000x2000",
			height: 10,
			width:  10,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/captcha", nil)
			responseRecorder := httptest.NewRecorder()

			generator := captcha.NewBasicGenerator(captcha.GeneratorConfig{Height: tc.height, Width: tc.width, Circles: 10})
			logger := log.NewNopLogger()
			handler := createHandler(logger, generator)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != http.StatusOK {
				t.Errorf("want status '%d', got '%d'", http.StatusOK, responseRecorder.Code)
			}
			cpt := captcha.Captcha{}
			if err := json.Unmarshal(responseRecorder.Body.Bytes(), &cpt); err != nil {
				t.Errorf("decoding response body err: %s", err)
			}
			if cpt.Solution.X == 0 && cpt.Solution.Y == 0 && cpt.Solution.H == 0 {
				t.Errorf("decoding response body to zero")
			}
			if cpt.Solution.X > tc.width || cpt.Solution.Y > tc.height {
				t.Errorf("unexpected solution")
			}
			if string(cpt.Data[1:4]) != "PNG" {
				t.Errorf("want status '%s', got '%s'", cpt.Data[1:4], "dddd")
			}
		})
	}
}
