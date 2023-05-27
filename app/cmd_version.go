package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils/colorize"
	log "github.com/sirupsen/logrus"
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

	v, err := latestVersion(10)
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

var versionCh chan string

// Start the process to fetch the latest tag
// Retrieve it with latestVersion(timeout)
func fetchLatestVersion() {
	versionCh = make(chan string, 1)
	go _fetchLatestVersion()
}

func _fetchLatestVersion() {
	client := &http.Client{}
	res, err := client.Get("https://api.github.com/repos/labctl/labctl/releases/latest")
	if err != nil || res.StatusCode != 200 {
		versionCh <- fmt.Sprintf("error occurred during latest version fetch: %v", err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		versionCh <- fmt.Sprintf("error reading version response: %s", err)
		return
	}

	r := ghapi{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		versionCh <- fmt.Sprintf("error decoding GitHub API response: %s", err)
		return
	}
	ver := strings.TrimLeft(r.TagName, "v")
	if version == "0.0.0" {
		version = ver + "-next"
	}
	if Ctx.DebugCount > 3 {
		log.Debugf("Current version: %s, available version: %s", version, ver)
	}
	versionCh <- ver
}

// Log the latest version if already available (and not already done).
// Dont wait more than 1 second
func logLatestVersion(delay_s time.Duration) {
	v, err := latestVersion(delay_s) // only wait for 1 second usually
	if err != nil {
		log.Errorf("%s", err)
	}
	if v != "" && v != version {
		log.Infof("🎉 New labctl version %s available! (current version: %s)", v, version)
	}
}

// Get the latestVersion. Should follow fetchLatestVersion
// Can potentially block up to delay_s seconds
func latestVersion(delay_s time.Duration) (string, error) {
	if versionCh == nil {
		return "", nil
	}
	defer func() {
		close(versionCh)
		versionCh = nil
	}()

	if Ctx.DebugCount > 1 && delay_s > 1 {
		log.Debugf("Waiting up to %v for the labctl version check to github.com", delay_s*time.Second)
	}

	select {
	case ver, ok := <-versionCh:
		if ok {
			if strings.HasPrefix(ver, "error") {
				return "", fmt.Errorf(ver)
			}
			return ver, nil
		}
	case <-time.After(delay_s * time.Second):
		return "", fmt.Errorf("timeout fetching the new version")
	}
	return "", fmt.Errorf("could not check the version")
}
