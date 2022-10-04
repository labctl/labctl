package webpty

import (
	"strings"

	"github.com/google/shlex"
	"github.com/labctl/labctl/utils"
	"github.com/pkg/errors"
)

func allow(cmd string) (string, []string, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return "", nil, errors.New("empty command")
	}
	args, err := shlex.Split(cmd)
	if err != nil {
		return "", nil, errors.Wrapf(err, "could not split command")
		// helpers.WsErrorf(ws, "%s\n\nwhile trying to split:\n\t%s", err, cmd)
		// return
	}
	if len(args) == 0 {
		return "", nil, errors.New("no cmd specified.")
	}
	c := args[0]
	if c == "labctl" && strings.Contains(cmd, "color") && !strings.HasPrefix(cmd, "labctl color ssh") {
		return "", nil, errors.New("only 'labctl color ssh <host>' allowed!")
	}

	if !utils.Contains([]string{"ping", "clab", "labctl"}, c) {
		return "", nil, errors.Errorf("command not allowed: %s", cmd)
	}
	return args[0], args[1:], nil
}
