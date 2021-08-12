package transport

import (
	"fmt"
	"strings"
)

// implements SShKind
type SSHKindSRL struct{}

func (sk *SSHKindSRL) Start(s *SSHTransport, transaction bool) error {
	s.PromptChar = "#" // ensure it's '#'
	r0 := s.Run("enter candidate private", 5)
	r1 := s.Run("discard stay", 2)
	if !strings.Contains(r1.result, "Nothing to discard") {
		r0.result += "; " + r1.result
		r0.command += "; " + r1.command
	}
	r0.Log(s.Target)
	return nil
}

func (sk *SSHKindSRL) Commit(s *SSHTransport) (*SSHReply, error) {
	r := s.Run("commit now", 10)
	if strings.Contains(r.result, "All changes have been committed") {
		r.result = ""
	} else {
		return r, fmt.Errorf("could not commit '%s'", r.result)
	}
	return r, nil
}

func (sk *SSHKindSRL) Compare(s *SSHTransport) (*SSHReply, error) {

	return nil, nil
}

func (sk *SSHKindSRL) Discard(s *SSHTransport) (*SSHReply, error) {
	r1 := s.Run("discard stay", 2)
	return r1, nil
}

func (sk *SSHKindSRL) PromptParse(s *SSHTransport, in *string) *SSHReply {
	return promptParseNoSpaces(in, s.PromptChar, 2)
}
