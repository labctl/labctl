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
	fs http.FileSystem
}

// Open file with fallback to index.html
// https://gist.github.com/hauxe/f2ea1901216177ccf9550a1b8bd59178#file-http_static_correct-go
func (spafs SPAFileSystem) Open(path string) (http.File, error) {
	log.Infof("serve %s", path)
	var fi fs.FileInfo
	f, err := spafs.fs.Open(path)
	if err == nil {
		fi, err = f.Stat()
	}
	if err != nil || fi.IsDir() {
		log.Errorf("err1 %s", err)
		f, err = spafs.fs.Open("/index.html")
		if err != nil {
			log.Fatal(err)
		}
		return f, nil
	}

	log.Infof("all ok %s", path)
	return spafs.fs.Open(path)
}

// Embed static web app. Ensure we have index.html
//go:embed html/index.html
//go:embed html/*
var labctlHtml embed.FS

func LabctlFileServer() http.Handler {
	return SinglePageAppFileServer(labctlHtml, "html")
}

func SinglePageAppFileServer(hfs fs.FS, root string) http.Handler {
	spaFS, err := fs.Sub(hfs, strings.TrimRight(root, "/"))
	if err != nil {
		log.Fatal(err)
	}
	// Ensure we have index.html
	_, err = fs.Stat(spaFS, "index.html")
	if err != nil {
		log.Fatal(err)
	}

	return http.FileServer(SPAFileSystem{http.FS(spaFS)})
}
