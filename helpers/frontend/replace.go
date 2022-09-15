package frontend

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
)

type ReplaceText struct {
	File http.File
	Url  string
}

func (r ReplaceText) Close() error {
	return r.File.Close()
}

func (r ReplaceText) Seek(offset int64, whence int) (int64, error) {
	return r.File.Seek(offset, whence)
}

func (r ReplaceText) Readdir(count int) ([]fs.FileInfo, error) {
	return r.File.Readdir(count)
}

func (r ReplaceText) Stat() (fs.FileInfo, error) {
	return r.File.Stat()
}

// Replace in-place. Only up to length of s0
func replaceAll(p []byte, s0 []byte, s1 []byte) {
	for {
		li := bytes.Index(p, s0)
		if li < 0 {
			break
		}
		for i := 0; i < len(s0); i++ {
			p[li+i] = s1[i]
		}
	}
}

func (r ReplaceText) Read(p []byte) (int, error) {
	c, err := r.File.Read(p)
	if err != nil || r.Url == "/labctl" {
		return c, err
	}
	// html and js
	replaceAll(p, []byte(`"/labctl`), []byte(fmt.Sprintf("%8s", `"`+r.Url)))
	// css fonts
	replaceAll(p, []byte(`url(/labctl`), []byte(fmt.Sprintf("url(%7s", r.Url)))
	return c, err
}
