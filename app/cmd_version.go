package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/chigopher/pathlib"
	"github.com/fatih/color"
	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils/colorize"
)

var (
	version = "0.0.0"
	commit  = "<commit>"
	date    = "<date>"
)

type CmdVersion struct {
	Upgrade bool `short:"u" help:"Upgrade labctl" default:"false"`
}

func (r *CmdVersion) Run(ctx *helpers.Context) error {
	fmt.Print("\n⚙🛠  ", colorize.ColorStr("higreen,bold", "Labctl"), " - ", colorize.ColorStr("italic", "https://labctl.net"), "\n\n")
	fmt.Printf("   version: %s\n", color.HiWhiteString(version))
	fmt.Printf("    commit: %s\n", commit)
	fmt.Printf("      date: %s\n", date)
	fmt.Print("\n")

	v, err := LatestVersion.Get(10)
	if err != nil {
		fmt.Printf("could not fetch latest version: %s\n", err)
	}
	if v != "" && v != version {
		fmt.Printf("🎉 New labctl version %s is available!\n", v)
	}
	if v == version {
		fmt.Printf("🎉 You have the latest version!\n")
		return nil
	}

	if r.Upgrade {
		fmt.Print(`
 to upgrade, run:
  bash -c "$(curl -sL https://labctl.net/get.sh)"
`)
	}
	return nil
}

type ghapi struct {
	TagName string `json:"tag_name"`
}

type Version struct {
	ch chan string
}

var LatestVersion Version

// Fetch the latest version in the background
func (v *Version) Fetch() {
	v.ch = make(chan string, 1)
	go v.fetch()
}

func (v *Version) fetch() {
	fVer := fileVersion()
	if fVer != "" {
		v.ch <- fVer
		return
	}

	client := &http.Client{}
	res, err := client.Get("https://api.github.com/repos/labctl/labctl/releases/latest")
	if err != nil || res.StatusCode != 200 {
		v.ch <- fmt.Sprintf("error occurred during latest version fetch: %v [%d]", err, res.StatusCode)
		return
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Debug("Error closing request", "err", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		v.ch <- fmt.Sprintf("error reading version response: %s", err)
		return
	}

	r := ghapi{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		v.ch <- fmt.Sprintf("error decoding GitHub API response: %s", err)
		return
	}
	ver := strings.TrimLeft(r.TagName, "v")
	if version == "0.0.0" {
		version = ver + "-next"
	}
	if Ctx != nil && Ctx.DebugCount > 3 {
		log.Debugf("Current version: %s, available version: %s", version, ver)
	}
	v.ch <- ver

	fileVersion(ver)
}

// Log the latest version. Delay up to delay_s
func (v *Version) Log(delay_s time.Duration) {
	res, err := v.Get(delay_s) // only wait for 1 second usually
	if err != nil {
		log.Error(err)
	}
	if res != "" && res != version {
		log.Infof("🎉 New labctl version %s available! (current version: %s)", res, version)
	}
}

// Get the latest version from the channel
func (v *Version) Get(delay_s time.Duration) (string, error) {
	if v.ch == nil {
		return "", nil
	}
	defer func() {
		close(v.ch)
		v.ch = nil
	}()

	if Ctx != nil && Ctx.DebugCount > 1 && delay_s > 1 {
		log.Debugf("Waiting up to %v for the labctl version check to github.com", delay_s*time.Second)
	}

	select {
	case ver, ok := <-v.ch:
		if ok {
			if strings.HasPrefix(ver, "error") {
				return "", fmt.Errorf("%s", ver)
			}
			return ver, nil
		}
	case <-time.After(delay_s * time.Second):
		return "", fmt.Errorf("timeout fetching the new version")
	}
	return "", fmt.Errorf("could not check the version")
}

// Set/Get the version from a file (mtime < 1hour ago)
func fileVersion(setV ...string) string {
	p := pathlib.NewPath(os.TempDir()).Join(".labctl.version")
	if len(setV) > 0 {
		err := p.WriteFile([]byte(strings.Join(setV, ".")))
		if err != nil {
			log.Warn("could not write version file", "file", p.String(), "err", err)
		}
		return ""
	}

	info, err := p.Stat()
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}

	if err != nil || info.ModTime().Add(time.Hour).Before(time.Now()) {
		err := p.Remove()
		if err != nil {
			log.Warn("could not access version file", "file", p.String(), "err", err)
		}
		return ""
	}
	val, err := p.ReadFile()
	if err != nil {
		log.Warn("could not read version file", "file", p.String(), "err", err)
	}
	return string(val)
}
