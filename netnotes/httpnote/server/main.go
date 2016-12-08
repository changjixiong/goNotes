package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Ser struct {
}

func (s *Ser) foo(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, strings.Repeat("foo", 100))

}

func (s *Ser) bar(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, strings.Repeat("bar", 100))
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func MakeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func main() {

	s := Ser{}

	http.HandleFunc("/foo", http.HandlerFunc(s.foo))
	http.HandleFunc("/bar", MakeGzipHandler(s.bar))

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
