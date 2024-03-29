package tx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/scrapli/scrapligo/driver/opoptions"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/response"
	"github.com/scrapli/scrapligo/util"
	"github.com/srl-labs/containerlab/clab/config"
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

	// Usernames & passwords
	ssh_user, ssh_pass := "admin", "admin"
	// if v, ok := nodes.GetDefaultCredentialsForKind(st.TargetNode.Kind); ok == nil {
	// 	ssh_user, ssh_pass = v[0], v[1]
	// }
	ssh_user = get(st.Vars, vkSshUser, ssh_user)
	ssh_pass = get(st.Vars, vkSshPass, ssh_pass)

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
			ssh_platform = r.Platform
		} else {
			return nil, fmt.Errorf("no platform definition for node kind %s", st.TargetNode.Kind)
		}
	}

	opt := []util.Option{
		options.WithTransportType("standard"),
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(ssh_user),
		options.WithAuthPassword(ssh_pass),
		options.WithPort(ssh_port),
		options.WithSSHConfigFileSystem(),
		options.WithSSHKnownHostsFileSystem(),
		// allows adding additional algorithms
		// Host clab-*
		// .  HostkeyAlgorithms +ssh-rsa
		// .  PubkeyAcceptedAlgorithms +ssh-rsa
	}

	opt = append(opt, options.WithLogger(scrapliLoggerInstance()))

	p, err := platform.NewPlatform(ssh_platform, ssh_host, opt...)
	if err != nil {
		log.Errorf("failed to create platform for %s; error: %s", st.TargetNode.Kind, err)
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
		log.Errorf("failed to fetch network driver from the platform for %s/%s; error: %s", st.TargetNode.Kind, ssh_platform, err)
	}

	err = d.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open %s; error: %s (Can you ssh to the node?)", ssh_host, err)
	}
	defer d.Close()

	var result []*Response

	for _, lines := range st.Config {
		count, istart := len(lines.Commands), len(result)

		// Transaction commands
		var actionCmds []string
		if km, ok := KindMap[st.TargetNode.Kind]; ok {
			if action == ACompare {
				actionCmds = km.Compare
			} else if action == ACommit {
				actionCmds = km.Commit
			}
		}
		lines.Commands = append(lines.Commands, actionCmds...)

		// Scarpligo's response
		var mres *response.MultiResponse

		default_level := log.DebugLevel
		if action == ASend {
			mres, err = d.SendCommands(lines.Commands, opoptions.WithNoStripPrompt())
		} else {
			mres, err = d.SendConfigs(lines.Commands, opoptions.WithNoStripPrompt())
			default_level = log.InfoLevel
		}
		if err != nil {
			return nil, fmt.Errorf("failed to send configs/commands; error: %+v", err)
		}

		// convert
		for _, res := range mres.Responses {
			resNoPrompt := d.Channel.PromptPattern.ReplaceAllString(res.Result, "")
			if len(resNoPrompt) > 0 {
				res := &Response{
					Node:     st.TargetNode.ShortName,
					Source:   lines.Source,
					Prompt:   res.Result[len(resNoPrompt):],
					Command:  res.Input,
					Response: resNoPrompt,
					Level:    default_level,
				}

				// Mark action commands
				for _, c := range actionCmds {
					if res.Command == c {
						res.Level = log.InfoLevel
						res.Source += " " + action.Char()
						if strings.HasSuffix(c, " ") {
							// action commands ending with a space considered "silent"
							istart += 1
							res.Level = log.DebugLevel
						}
					}
				}

				result = append(result, res)
			}
		}

		changes := len(result) - istart

		if action == ACommit {
			r := &Response{
				Node:     st.TargetNode.ShortName,
				Source:   lines.Source + " " + action.Char(),
				Response: fmt.Sprintf("committed %d lines", count),
				Level:    log.InfoLevel,
			}
			if changes > 0 {
				r.Response += fmt.Sprintf(", %d failed commands", changes)
				r.Level = log.WarnLevel
			}
			result = append(result, r)
		}
		if action == ACompare {
			if changes == 0 {
				r := &Response{
					Node:     st.TargetNode.ShortName,
					Source:   lines.Source + " " + action.Char(),
					Response: "No changes to the configuration",
					Level:    log.InfoLevel,
				}
				result = append(result, r)
			}
		}
	}
	if log.DebugLevel > 1 {
		log.Debugf("send done, closing...")
	}
	return result, nil
}

func scrapliLoggerInstance() *logging.Instance {
	lvl := logging.Critical
	fun := log.Error
	if DebugCount > 0 {
		lvl = logging.Info
		fun = log.Info
	} else if DebugCount > 3 {
		lvl = logging.Debug
		fun = log.Debug
	}

	li, err := logging.NewInstance(
		logging.WithLevel(lvl), logging.WithLogger(fun),
	)
	if err != nil {
		log.Error(err)
	}
	return li
}
