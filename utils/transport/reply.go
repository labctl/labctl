package transport

import (
	"fmt"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// The SSH reply, executed command and the prompt
type SSHReply struct{ Result, Prompt, Command string }

// The LogString will include the entire SSHReply
//   Each field will be prefixed by a character.
//   # - command sent
//   | - result received
//   ? - prompt part of the result
func (r *SSHReply) LogString(node string, linefeed, debug bool) string {
	ind := 12 + len(node)
	prefix := "\n" + strings.Repeat(" ", ind)
	s := ""
	if linefeed {
		s = "\n" + strings.Repeat(" ", 11)
	}
	s += node + " # " + r.Command
	s += prefix + "| "
	s += strings.Join(strings.Split(r.Result, "\n"), prefix+"| ")
	if debug { // Add the prompt & more
		s = "" + strings.Repeat(" ", ind) + s
		s += prefix + "? "
		s += strings.Join(strings.Split(r.Prompt, "\n"), prefix+"? ")
		if DebugCount > 3 { // add bytestring
			s += fmt.Sprintf("%s| %v%s ? %v", prefix, []byte(r.Result), prefix, []byte(r.Prompt))
		}

	}
	return s
}

func (r *SSHReply) Log(node string, level ...log.Level) *SSHReply {
	if r.Result == "" {
		return r
	}
	if len(level) == 0 {
		log.Info(r.LogString(node, false, false))
		return r
	}
	switch level[0] {
	case log.WarnLevel:
		log.Warn(r.LogString(node, false, false))
	case log.ErrorLevel:
		log.Warn(r.LogString(node, false, false))
	default:
		log.Info(r.LogString(node, false, false))
	}
	return r
}

func (r *SSHReply) Debug(node, message string) {
	msg := message
	// if len(t) > 0 {
	// 	msg = t[0].(string)
	// }
	_, fn, line, _ := runtime.Caller(1)
	msg += fmt.Sprintf("(%s line %d)", fn, line)
	msg += r.LogString(node, true, true)
	log.Debug(msg)
}
