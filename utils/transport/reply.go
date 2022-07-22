package transport

import (
	"fmt"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// The SSH reply, executed command and the prompt
type SSHReply struct{ Node, Source, Prompt, Command, Response string }

// The Slog will include the entire SSHReply
//   Each field will be prefixed by a character.
//   # - command sent
//   | - result received
//   ? - prompt part of the result
func (r *SSHReply) Slog() string {

	s := fmt.Sprintf("%s execute '%s'", r.Node, r.Command)
	if r.Source != "" {
		s += fmt.Sprintf(" (from template %s)", r.Source)
	}
	const lf0 = "\n ┆ "
	const lf = "\n | "
	s += lf0
	s += strings.Join(strings.Split(r.Prompt, "\n"), lf0)
	s += " " + r.Command + lf
	s += strings.Join(strings.Split(r.Response, "\n"), lf)

	if DebugCount > 3 { // add bytestring
		s += fmt.Sprintf("%s| %v%s ? %v", lf, []byte(r.Response), lf, []byte(r.Prompt))
	}

	return s
}

func (r *SSHReply) Log(level ...log.Level) *SSHReply {
	if r.Response == "" {
		return r
	}
	if len(level) == 0 {
		log.Info(r.Slog())
		return r
	}
	switch level[0] {
	case log.WarnLevel:
		log.Warn(r.Slog())
	case log.ErrorLevel:
		log.Error(r.Slog())
	default:
		log.Info(r.Slog())
	}
	return r
}

func (r *SSHReply) Debug(message string) {
	msg := message
	// if len(t) > 0 {
	// 	msg = t[0].(string)
	// }
	_, fn, line, _ := runtime.Caller(1)
	msg += fmt.Sprintf("(%s line %d)", fn, line)
	msg += r.Slog()
	log.Debug(msg)
}
