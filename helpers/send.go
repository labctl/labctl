package helpers

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/labctl/labctl/utils/transport"
	"github.com/scrapli/scrapligo/driver/opoptions"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/util"

	log "github.com/sirupsen/logrus"
	"github.com/srl-labs/containerlab/clab/config"
	"github.com/srl-labs/containerlab/nodes"
)

const (
	ActionSend    = "send"
	ActionCompare = "compare"
	ActionCommit  = "commit"
)

// validate that we have only allowed actions
func validateSendAction(action string) error {
	switch action {
	case ActionSend, ActionCommit, ActionCompare:
		return nil
	}
	return fmt.Errorf("unsupported action: %s", action)
}

func parseAddrPort(host string) (string, int, error) {
	ap, err := netip.ParseAddrPort(host)
	if err == nil {
		return ap.Addr().String(), int(ap.Port()), nil
	}
	prt := strings.Split(host, ":")
	if len(prt) == 2 {
		port, err := strconv.Atoi(prt[1])
		if err == nil {
			return prt[0], port, nil
		}
	}
	return "", 0, fmt.Errorf("invalid address:port %s", err)
}

const (
	vkSshAddr   = "ssh"          // Override host:port (default=LongName:21 from containerlab)
	vkSshDriver = "ssh_driver"   // Enables scrapli, with the specified driver
	vkSshUser   = "ssh_username" // Username
	vkSshPass   = "ssh_password" // Password, if not the default
)

// Splits configuration into lines that shoul dbe sent to the device
func configSplitLines(c *config.NodeConfig) ([][]string, []string) {
	var res [][]string
	var lres []string
	var infores []string

	for idx, data := range c.Data {
		lres = make([]string, 0)

		for _, l := range strings.Split(data, "\n") {
			l = strings.TrimSpace(l)
			if l == "" || strings.HasPrefix(l, "#") {
				continue
			}
			lres = append(lres, l)
		}
		if len(lres) > 0 {
			res = append(res, lres)
			infores = append(infores, c.Info[idx])
		}
	}
	return res, infores
}

type KindMap struct {
	driver          string
	commit, compare []string
}

var KMap = map[string]KindMap{
	"vr-sros": {
		driver:  "nokia_sros",
		commit:  []string{"validate", "commit", "discard"},
		compare: []string{"validate", "compare /", "discard"},
	},
	"srl": {
		driver:  "nokia_sros",
		commit:  []string{"commit now", "discard stay"},
		compare: []string{"compare", "discard"},
	},
}

func ConfigSend(c *config.NodeConfig, action string) ([]*transport.SSHReply, error) {

	var err error
	var result []*transport.SSHReply

	// break config into lines (but still keep templates separate)
	to_send, to_send_info := configSplitLines(c)
	if len(to_send) == 0 {
		return result, nil
	}

	// SSH usernames & passwords
	ssh_user, ssh_pass := "admin", "admin"
	if v, ok := nodes.GetDefaultCredentialsForKind(c.TargetNode.Kind); ok == nil {
		ssh_user, ssh_pass = v[0], v[1]
	}
	if v, ok := c.TargetNode.Config.Vars[vkSshUser]; ok {
		ssh_user = fmt.Sprintf("%v", v)
	}
	if v, ok := c.TargetNode.Config.Vars[vkSshPass]; ok {
		ssh_pass = fmt.Sprintf("%v", v)
	}

	// SSH override host / port, more generic approach than just containerlab
	ssh_host, ssh_port := c.TargetNode.LongName, 22
	if v, ok := c.TargetNode.Config.Vars[vkSshAddr]; ok {
		ssh_host, ssh_port, err = parseAddrPort(fmt.Sprintf("%v", v))
		if err != nil {
			return nil, err
		}
	}

	// Get the ssh_driver
	ssh_driver, ok := c.TargetNode.Config.Vars[vkSshDriver]
	if !ok {
		if r, ok := KMap[c.TargetNode.Kind]; ok {
			ssh_driver = r.driver
		}
	}

	if ssh_driver == "original" {
		tx, err := transport.NewSSHTransport(
			c.TargetNode,
			transport.WithUserNamePassword(ssh_user, ssh_pass),
			transport.HostKeyCallback(),
		)
		if err != nil {
			return nil, err
		}
		tx.Port = ssh_port
		tx.Target = c.TargetNode.ShortName

		err = tx.Connect(ssh_host)
		if err != nil {
			return nil, err
		}
		defer tx.Close()

		transaction := action != ActionSend // commit & compare requires transactions

		for _, d1 := range to_send {
			//time.Sleep(10 * time.Millisecond)

			err := tx.K.Start(tx, transaction)
			if err != nil {
				return nil, err
			}

			sshr, err := tx.Write(d1)
			if err != nil {
				return nil, fmt.Errorf("could not write config: %s\n%s", err, d1)
			}

			for _, r := range sshr {
				if r.Response != "" {
					result = append(result, r)
				}
			}

			if transaction {
				var reply *transport.SSHReply
				//var msg string

				switch action {
				case ActionCommit:
					reply, err = tx.K.Commit(tx)
					// msg = fmt.Sprintf("%s: COMMIT - %d lines", c.TargetNode.ShortName, len(d1))
				case ActionCompare:
					reply, err = tx.K.Compare(tx)
					if err == nil {
						_, err = tx.K.Discard(tx)
					}
				}
				if reply.Response != "" {
					result = append(result, reply)
					//msg += reply.Slog()
				}
				// 				if err != nil {
				// 					log.Error(msg)
				// 					return nil, err
				// 				}
				// 				if len(msg) > 0 {
				// //					log.Info(msg)
				// 					result = append(result, reply)

				// 				}
			}
		}
	} else {

		opt := []util.Option{
			options.WithAuthNoStrictKey(),
			options.WithAuthUsername(ssh_user),
			options.WithAuthPassword(ssh_pass),
			options.WithPort(ssh_port),
		}

		if DebugCount > 3 {
			opt = append(opt, options.WithDefaultLogger())
		}

		p, err := platform.NewPlatform(ssh_driver, ssh_host, opt...)
		if err != nil {
			log.Errorf("failed to create platform for %s; error: %s", c.TargetNode.Kind, err)
		}

		d, err := p.GetNetworkDriver()
		if err != nil {
			log.Errorf("failed to fetch network driver from the platform for %s; error: %s", c.TargetNode.Kind, err)
		}

		err = d.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open driver for %s; error: %+v", c.TargetNode.Kind, err)
		}
		defer d.Close()

		for idx, data := range to_send {

			if action == ActionSend {
				return nil, fmt.Errorf("send not implemented")
			}

			count, istart := len(data), len(result)

			if action == ActionCompare {
				if r, ok := KMap[c.TargetNode.Kind]; ok {
					for _, c := range r.compare {
						data = append(data, c)
					}
				}
			}
			if action == ActionCommit {
				if r, ok := KMap[c.TargetNode.Kind]; ok {
					for _, c := range r.commit {
						data = append(data, c)
					}
				}

			}

			res, err := d.SendConfigs(data, opoptions.WithNoStripPrompt())
			if err != nil {
				return nil, fmt.Errorf("failed to send configs; error: %+v", err)
			}

			for _, res := range res.Responses {
				response := d.Channel.PromptPattern.ReplaceAllString(res.Result, "")
				if len(response) > 0 {
					p := res.Result[len(response):]
					result = append(result, &transport.SSHReply{
						Node:     c.TargetNode.LongName,
						Source:   to_send_info[idx],
						Prompt:   p,
						Command:  res.Input,
						Response: response,
					})
				}
			}

			changes := len(result) - istart

			if action == ActionCommit {
				msg := fmt.Sprintf("%s comitted %d lines", c.TargetNode.LongName, count)
				if changes > 0 {
					msg += fmt.Sprintf("%d failed commands", changes)
					log.Warnln(msg)
				} else {
					log.Infoln(msg)
				}
			}
			if action == ActionCompare {
				if changes == 0 {
					log.Infof("%s no changes", c.TargetNode.LongName)
				}
			}

		}
	}

	for _, r := range result {
		if r.Command == "compare /" {
			r.Log(log.InfoLevel)
		} else {
			r.Log(log.ErrorLevel)
		}

	}

	return result, nil

}
