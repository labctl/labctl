package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labctl/labctl/utils/transport"
	scrapli "github.com/scrapli/scrapligo/driver/base"
	scraplicore "github.com/scrapli/scrapligo/driver/core"
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

func ConfigSend(c *config.NodeConfig, action string) ([]*transport.SSHReply, error) {

	var err error
	ct, ok := c.TargetNode.Labels["config.transport"]
	if !ok {
		ct = "ssh"
	}

	var result []*transport.SSHReply

	// Override host / port
	// more generic approach than just containerlabs
	host, port := c.TargetNode.LongName, 22
	if v, ok := c.TargetNode.Config.Vars["ssh_ip"]; ok {
		ct = "ssh"
		host = fmt.Sprintf("%v", v)
		prt := strings.Split(host, ":")
		if len(prt) > 1 {
			host = prt[0]
			port, err = strconv.Atoi(prt[1])
			if err != nil {
				return nil, fmt.Errorf("invalid ssh_ip %s", err)
			}
		}
	}

	switch ct {
	case "ssh":

		tx, err := transport.NewSSHTransport(
			c.TargetNode,
			transport.WithUserNamePassword(
				nodes.DefaultCredentials[c.TargetNode.Kind][0],
				nodes.DefaultCredentials[c.TargetNode.Kind][1]),
			transport.HostKeyCallback(),
		)
		if err != nil {
			return nil, err
		}
		tx.Port = port
		tx.Target = c.TargetNode.ShortName

		err = tx.Connect(host)
		if err != nil {
			return nil, err
		}
		defer tx.Close()

		transaction := action != ActionSend // commit & compare requires transactions

		for _, d1 := range c.Data {
			//time.Sleep(10 * time.Millisecond)

			err := tx.K.Start(tx, transaction)
			if err != nil {
				return nil, err
			}

			l, sshr, err := tx.Write(&d1)
			if err != nil {
				return nil, fmt.Errorf("could not write config: %s\n%s", err, d1)
			}

			for i := range sshr {
				result = append(result, sshr[i])

			}

			if transaction {
				var reply *transport.SSHReply
				var msg string

				switch action {
				case ActionCommit:
					reply, err = tx.K.Commit(tx)
					msg = fmt.Sprintf("%s: COMMIT - %d lines", c.TargetNode.ShortName, l)
				case ActionCompare:
					reply, err = tx.K.Compare(tx)
					if err == nil {
						_, err = tx.K.Discard(tx)
					}
				}
				if reply.Result != "" {
					msg += reply.LogString(c.TargetNode.ShortName, true, false)
				}
				if err != nil {
					log.Error(msg)
					return nil, err
				}
				if len(msg) > 0 {
					log.Info(msg)
				}
			}
		}

	case "scrapli":

		opt := []scrapli.Option{
			scrapli.WithAuthStrictKey(false),
			scrapli.WithAuthPassword(nodes.DefaultCredentials[c.TargetNode.Kind][1]),
			scrapli.WithAuthUsername(nodes.DefaultCredentials[c.TargetNode.Kind][0]),
		}

		opt = append(opt, scrapli.WithPort(port))

		d, err := scraplicore.NewCoreDriver(host, "nokia_sros", opt...)
		if err != nil {
			log.Errorf("%s: %s", c.TargetNode.Kind, err)
		}

		err = d.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open driver; error: %+v", err)
		}
		defer d.Close()

		// send some configs
		_, err = d.SendConfigs(c.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to send configs; error: %+v", err)
		}

	case "grpc":
		return nil, fmt.Errorf("transport grpc not implemented yet")

	default:
		return nil, fmt.Errorf("unknown transport: %s", ct)
	}
	return result, nil

}
