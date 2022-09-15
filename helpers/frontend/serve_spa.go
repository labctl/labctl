package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// FileSystem custom file system handler
type SPAFileSystem struct {
	fs      http.FileSystem
	urlRoot string
}

func (spafs SPAFileSystem) index(path, msg string) (http.File, error) {
	log.Warnf("%s --> index.html: %s", path, msg)
	f, err := spafs.fs.Open("/index.html")
	if err != nil {
		log.Fatal(err)
	}
	return ReplaceText{
		File: f,
		Url:  spafs.urlRoot,
	}, nil
}

// Open file with fallback to index.html
// https://gist.github.com/hauxe/f2ea1901216177ccf9550a1b8bd59178#file-http_static_correct-go
func (spafs SPAFileSystem) Open(path string) (http.File, error) {
	path = strings.TrimPrefix(path, spafs.urlRoot)
	f, err := spafs.fs.Open(path)
	if err != nil {
		return spafs.index(path, err.Error())
	}
	info, err := f.Stat()
	if err != nil {
		return spafs.index(path, err.Error())
	}
	if info.IsDir() {
		f.Close()
		return spafs.index(path, "isDir()")
	}
	log.Debugf("OK '%s'", path)
	return ReplaceText{
		File: f,
		Url:  spafs.urlRoot,
	}, nil
}

// Embed static web app. Ensure we have index.html
//go:embed html/index.html
//go:embed html/*
var labctlHtml embed.FS

func LabctlFileServer(url string) http.Handler {
	return SinglePageAppFileServer(labctlHtml, url, "html")
}

func SinglePageAppFileServer(hfs fs.FS, urlRoot, fsRoot string) http.Handler {
	spaFS, err := fs.Sub(hfs, strings.TrimRight(fsRoot, "/"))
	if err != nil {
		log.Fatal(err)
	}
	// Ensure we have index.html
	_, err = fs.Stat(spaFS, "index.html")
	if err != nil {
		log.Fatal(err)
	}

	return http.FileServer(SPAFileSystem{
		fs:      http.FS(spaFS),
		urlRoot: strings.TrimRight(urlRoot, "/"),
	})
}

type SlashFix struct {
	Mux http.Handler
}

func (h *SlashFix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.Replace(r.URL.Path, "//", "/", -1)
	r.URL.Path = strings.TrimRight(r.URL.Path, "/")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}
	h.Mux.ServeHTTP(w, r)
}
