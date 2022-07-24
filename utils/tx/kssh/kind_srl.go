package kssh

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
	if !strings.Contains(r1.Response, "Nothing to discard") {
		r0.Response += "; " + r1.Response
		r0.Command += "; " + r1.Command
	}
	r0.Log()
	return nil
}

func (sk *SSHKindSRL) Commit(s *SSHTransport) (*SSHReply, error) {
	r := s.Run("commit now", 10)
	if strings.Contains(r.Response, "All changes have been committed") {
		r.Response = ""
	} else {
		return r, fmt.Errorf("could not commit '%s'", r.Response)
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
