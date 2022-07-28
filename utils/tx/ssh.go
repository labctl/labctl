package tx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/scrapli/scrapligo/driver/opoptions"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/util"
	"github.com/srl-labs/containerlab/clab/config"
	"github.com/srl-labs/containerlab/nodes"
	"github.com/srl-labs/containerlab/types"

	log "github.com/sirupsen/logrus"
)

// configLines that will be sent to the NE
type configLines struct {
	Commands []string // a list of commands, ready to send
	Source   string   // the source of these commands (i.e. template filename)
}

const (
	vkSshHost     = "ssh_host"     // (optional) Override host (default=LongName from containerlab)
	vkSshPort     = "ssh_port"     // (optional) Override port (default=22)
	vkSshPlatform = "ssh_platform" // (optional) Override the scrapli platform definition
	vkSshUser     = "ssh_username" // (optional) Username. Fallback to the Kind default
	vkSshPass     = "ssh_password" // (optional) Password. Fallback to the Kind default
)

type SSHTx struct { // implement the Tx interface
	TargetNode *types.NodeConfig
	Vars       map[string]interface{}
	Config     []configLines
}

// Get a value from map with a default
func get(vars map[string]interface{}, key string, defaultv string) string {
	if val, ok := vars[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return defaultv
}

// Prepare Config to be sent (from rendered templates in config.NodeConfig)
// Split a rendered template into lines, remove whitespace and comments.
// Returns an array of configLines, one entry per template
func (st *SSHTx) Prepare(c *config.NodeConfig) (int, error) {
	st.Config = make([]configLines, 0)
	st.TargetNode = c.TargetNode
	st.Vars = c.Vars

	for i, data := range c.Data {
		res := configLines{
			Commands: make([]string, 0),
			Source:   c.Info[i],
		}

		for _, l := range strings.Split(data, "\n") {
			l = strings.TrimSpace(l)
			if l == "" || strings.HasPrefix(l, "#") {
				continue
			}
			res.Commands = append(res.Commands, l)
		}
		if len(res.Commands) > 0 {
			st.Config = append(st.Config, res)
		}
	}
	return len(st.Config), nil
}

func (st *SSHTx) Send(action Action) ([]*Response, error) {

	if len(st.Config) == 0 {
		return nil, nil
	}

	var err error
	var result []*Response

	// Usernames & passwords
	ssh_user, ssh_pass := "admin", "admin"
	if v, ok := nodes.GetDefaultCredentialsForKind(st.TargetNode.Kind); ok == nil {
		ssh_user, ssh_pass = v[0], v[1]
	}
	ssh_user = get(st.Vars, vkSshUser, ssh_user)
	ssh_user = get(st.Vars, vkSshPass, ssh_pass)

	// Host & ports, allows an alternative to containerlab's host entries
	ssh_host := get(st.Vars, vkSshHost, st.TargetNode.LongName)
	ssh_port, err := strconv.Atoi(get(st.Vars, vkSshPort, "22"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %s", err)
	}

	// Get the scrapligo platform definition
	ssh_platform, ok := st.TargetNode.Config.Vars[vkSshPlatform]
	if !ok {
		if r, ok := KindMap[st.TargetNode.Kind]; ok {
			ssh_platform = r.platform
		} else {
			return nil, fmt.Errorf("No platform definition for node kind %s", st.TargetNode.Kind)
		}
	}

	opt := []util.Option{
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(ssh_user),
		options.WithAuthPassword(ssh_pass),
		options.WithPort(ssh_port),
	}

	if DebugCount > 3 {
		opt = append(opt, options.WithDefaultLogger())
	}

	p, err := platform.NewPlatform(ssh_platform, ssh_host, opt...)
	if err != nil {
		log.Errorf("failed to create platform for %s; error: %s", st.TargetNode.Kind, err)
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
		log.Errorf("failed to fetch network driver from the platform for %s; stror: %s", st.TargetNode.Kind, err)
	}

	err = d.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open driver for %s; error: %+v", st.TargetNode.Kind, err)
	}
	defer d.Close()

	for _, cline := range st.Config {

		if action == ASend {
			return nil, fmt.Errorf("send not implemented")
		}

		count, istart := len(cline.Commands), len(result)

		// Transaction commands
		var actionCmds []string
		if r, ok := KindMap[st.TargetNode.Kind]; ok {
			if action == ACompare {
				actionCmds = r.compare
			} else if action == ACommit {
				actionCmds = r.commit
			}
		}
		for _, c := range actionCmds {
			cline.Commands = append(cline.Commands, c)
		}

		res, err := d.SendConfigs(cline.Commands, opoptions.WithNoStripPrompt())
		if err != nil {
			return nil, fmt.Errorf("failed to send configs; error: %+v", err)
		}

		for _, res := range res.Responses {
			response := d.Channel.PromptPattern.ReplaceAllString(res.Result, "")
			if len(response) > 0 {
				res := &Response{
					Node:     st.TargetNode.LongName,
					Source:   cline.Source,
					Prompt:   res.Result[len(response):],
					Command:  res.Input,
					Response: response,
				}

				// check if these are action commands
				for _, c := range actionCmds {
					if res.Command == c {
						res.Source += " " + strings.ToUpper(action.String())
					}
				}

				result = append(result, res)
			}
		}

		changes := len(result) - istart

		if action == ACommit {
			msg := fmt.Sprintf("%s committed %d lines", st.TargetNode.LongName, count)
			if changes > 0 {
				msg += fmt.Sprintf("%d failed commands", changes)
				log.Warnln(msg)
			} else {
				log.Infoln(msg)
			}
		}
		if action == ACompare {
			if changes == 0 {
				log.Infof("%s no changes", st.TargetNode.LongName)
			}
		}

	}

	for _, r := range result {
		if strings.Contains(r.Source, strings.ToUpper(action.String())) {
			r.Log(log.InfoLevel)
		} else {
			r.Log(log.ErrorLevel)
		}
	}
	return result, nil
}
