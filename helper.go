package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/eknkc/basex"
)

var (
	charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	encoder *basex.Encoding
)

func newEncoder() (*basex.Encoding, error) {
	if encoder != nil {
		return encoder, nil
	}
	var err error
	encoder, err = basex.NewEncoding(charset)
	if err != nil {
		return nil, err
	}
	return encoder, nil
}

// EncodeURL generate a base62 string with a counter to reduce the
// possibilities of collision
func EncodeURL(s string, count int64) (string, error) {
	encoder, err := newEncoder()
	if err != nil {
		return "", err
	}
	s = fmt.Sprintf("%s%d", s, count)
	encStr := encoder.Encode([]byte(s))
	return encStr[len(encStr)-6:], nil
}

func allowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(rw, req)
	})
}

func prefixHTTP(url string) string {
	if url[:4] != "http" {
		purl := fmt.Sprintf("https://%s", url)
		log.Println("Prefixed:", purl)
		return purl
	}
	return url
}

// IsUP verify that the remote site is up
// if not it replace the url with a
// WayBackMachine snapshot url of the website
func IsUP(url string) (string, bool) {
	_, err := http.Head(url)
	if err != nil {
		wbURL, err := fetchWBMlink(url)
		if err != nil {
			return url, false
		}
		return wbURL, false
	}
	return "", true
}

func fetchWBMlink(url string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://archive.org/wayback/available?url=%s", url))
	if err != nil {
		return "", err
	}
	m := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return "", err
	}
	return m["url"], nil
}
