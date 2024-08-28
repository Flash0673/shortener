package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type want struct {
	expectedCode int
	expectedLocation string
}

type test struct {
	method string
	path string
	body string
	db map[string]string
	want want
}
//https://practicum.yandex.ru/ POST
//LoYUInWAFzfCLCV7RkmeRSAMNLU3er87xAu2OLw GET -> https://practicum.yandex.ru/ 
func TestHandler(t *testing.T) {
	tests := []test{
		{
			method: http.MethodPost,
			path: "/",
			body: "https://practicum.yandex.ru/",
			db: make(map[string]string),
			want: want{
				expectedCode: http.StatusCreated,
				expectedLocation: "",
			},
		},
		{
			method: http.MethodGet,
			path: "/LoYUInWAFzfCLCV7RkmeRSAMNLU3er87xAu2OLw",
			db: map[string]string{"LoYUInWAFzfCLCV7RkmeRSAMNLU3er87xAu2OLw": "https://practicum.yandex.ru/"},
			want: want{
				expectedCode: http.StatusTemporaryRedirect,
				expectedLocation: "https://practicum.yandex.ru/",
			},
		},
		{
			method: http.MethodPut,
			path: "/",
			db: map[string]string{},
			want: want{
				expectedCode: http.StatusBadRequest,
				expectedLocation: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h := Handler(tt.db)
			h(w, req)

			resp := w.Result()

			assert.Equal(t, tt.want.expectedCode, resp.StatusCode)
			assert.Equal(t, tt.want.expectedLocation, resp.Header.Get("Location"))
		})
	}
}