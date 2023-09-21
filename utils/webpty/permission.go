package webpty

import (
	"strings"

	"github.com/google/shlex"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

// Check if the command is allowed to be executed
func allow(cmd string) (string, []string, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return "", nil, errors.New("empty command")
	}
	blocked_chars := []string{"&&", ";", "|", ">", "<"}
	for _, b := range blocked_chars {
		if strings.Contains(cmd, b) {
			return "", nil, errors.Errorf("your command may not contain: %s", b)
		}
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

	if c == "labctl" && !strings.HasPrefix(cmd, "labctl color ssh ") {
		return "", nil, errors.New("only 'labctl color ssh <host>' allowed!")
	}

	is_script := strings.HasPrefix(c, "./scripts/")
	is_allowed := slices.Contains([]string{"ping", "clab", "labctl"}, c)

	if !(is_allowed || is_script) {
		return "", nil, errors.Errorf("command not allowed: %s", cmd)
	}
	return args[0], args[1:], nil
}
