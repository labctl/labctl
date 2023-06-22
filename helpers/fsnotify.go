package helpers

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func WatchFS(ctx *Context, cb func(string)) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("Could not watch file changes: %s", err)
	}

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					// log.Println("modified file:", event.Name)
					cb(event.Name)
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
	err = watcher.Add(filepath.Dir(ctx.TopoFile))
	if err != nil {
		log.Errorf("Could not watch path %s: %s", filepath.Dir(ctx.TopoFile), err)
	}
	// Add template file paths
	for p := ctx.TemplatePaths.Oldest(); p != nil; p = p.Next() {
		err = watcher.Add(p.Value)
		if err != nil {
			log.Errorf("Could not watch path %s: %s", p.Value, err)
		}
	}
	return watcher
}
