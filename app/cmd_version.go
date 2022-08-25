package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	if r.Upgrade {
		if v == version {
			fmt.Printf("\nYou already have the latest version\n")
			return nil
		}
		fmt.Printf(`%s

 to upgrade, run:
  bash -c "$(curl -sL https://labctl.net/get.sh)"

`, v)
		return nil
	}

	if err != nil {
		fmt.Printf("could not fetch latest version: %s\n", err)
	}
	if v != "" {
		fmt.Println(v)
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
	body, err := ioutil.ReadAll(res.Body)
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
	newV := strings.TrimLeft(r.TagName, "v")

	if Ctx.DebugCount > 3 {
		log.Debugf("Current version: %s, available version: %s", version, newV)
	}
	if newV != version {
		versionCh <- newV
	} else {
		versionCh <- ""
	}
}

// Log the latest version if already available (and not already done).
// Dont wait more than 1 second
func logLatestVersion() {
	v, err := latestVersion(1) // only wait for 1 second usually
	if err != nil {
		log.Errorf("%s", err)
	}
	if v != "" {
		log.Infof(v)
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

	select {
	case ver, ok := <-versionCh:
		if ok {
			if strings.HasPrefix(ver, "error") {
				return "", fmt.Errorf(ver)
			}
			return fmt.Sprintf("🎉 New labctl version %s is available!", ver), nil
		}
	case <-time.After(delay_s * time.Second):
		return "", fmt.Errorf("timeout fetching the new version")
	}
	return "", fmt.Errorf("could not check the version")
}
