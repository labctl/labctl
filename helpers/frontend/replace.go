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

func (r ReplaceText) Read(p []byte) (int, error) {
	c, err := r.File.Read(p)
	if err != nil {
		return c, err
	}
	s0 := []byte(`"/labctl`)
	s1 := []byte(fmt.Sprintf("%8s", `"`+r.Url)[:8])
	for {
		li := bytes.Index(p, s0)
		if li < 0 {
			break
		}
		for i := 0; i < len(s0); i++ {
			p[li+i] = s1[i]
		}
	}

	return c, err
}
