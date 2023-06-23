package helpers

import (
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func WatchFS(ctx *Context, cb func(string)) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("Could not watch file changes: %s", err)
	}
	topoDir, err := filepath.Abs(filepath.Dir(ctx.TopoFile))
	if err != nil {
		log.Errorf("Could not start watcher on topo %s: %s", ctx.TopoFile, err)
		return watcher
	}

	// Start listening for events.
	go func() {
		last_name := ""
		last_now := time.Now()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op.Has(fsnotify.Write) {
					// Don't send more than 1 change per second
					if last_name == event.Name && time.Since(last_now).Seconds() < 1 {
						continue
					}
					last_name = event.Name
					last_now = time.Now()

					res, err := filepath.Rel(topoDir, event.Name)
					if err != nil {
						log.Warnf("Rel %s", event.Name)
						cb(event.Name)
					} else {
						cb(res)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add the topofile path
	err = watcher.Add(topoDir)
	if err != nil {
		log.Errorf("Could not watch path %s: %s", filepath.Dir(ctx.TopoFile), err)
	}
	// Add template file paths
	for p := ctx.TemplatePaths.Oldest(); p != nil; p = p.Next() {
		pth, err := filepath.Abs(p.Value)
		if err != nil {
			log.Errorf("Could not get abs of path %s: %s", pth, err)
		}
		err = watcher.Add(pth)
		if err != nil {
			log.Errorf("Could not watch path %s: %s", p.Value, err)
		}
	}
	return watcher
}
