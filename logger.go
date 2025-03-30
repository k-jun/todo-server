package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type RwWrapper struct {
	rw   http.ResponseWriter
	mw   io.Writer
	code int
}

func NewRwWrapper(rw http.ResponseWriter, buf io.Writer) *RwWrapper {
	return &RwWrapper{
		rw: rw,
		mw: io.MultiWriter(rw, buf),
	}
}

func (r *RwWrapper) Header() http.Header {
	return r.rw.Header()
}

func (r *RwWrapper) Write(i []byte) (int, error) {
	return r.mw.Write(i)
}

func (r *RwWrapper) WriteHeader(statusCode int) {
	r.code = statusCode
	r.rw.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(b))
		logger.InfoContext(r.Context(), "request body", slog.Any("requestBody", string(b)), slog.Any("requestHeader", r.Header))
		buf := &bytes.Buffer{}
		rw := NewRwWrapper(w, buf)
		next.ServeHTTP(rw, r)

		var responseBody interface{}
		_ = json.Unmarshal(buf.Bytes(), &responseBody)
		logger.InfoContext(r.Context(), "response body", slog.Any("responseBody", responseBody))

	})
}
