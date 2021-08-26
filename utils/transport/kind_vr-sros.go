package transport

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// implements SShKind
type SSHKindVRSROS struct{}

func (sk *SSHKindVRSROS) Start(s *SSHTransport, transaction bool) error {
	s.PromptChar = "#" // ensure it's '#'
	//s.debug = true
	r := s.Run("/environment more false", 5)
	if r.Result != "" {
		log.Warnf("%s Are you in MD-Mode?%s", s.Target, r.LogString(s.Target, true, false))
	}

	if transaction {
		r = s.Run("/configure global", 5)
		r.Result = strings.Replace(r.Result, "#2054", "2054", 1)
		r.Result = strings.Replace(r.Result, "INFO: CLI 2054: Entering global configuration mode", "", 1)
		//r.Result = strings.Replace(r.Result, "Entering global configuration mode", "", 1)
		r.Log(s.Target)
		s.Run("discard /", 1).Log(s.Target)
	}
	return nil
}

func (sk *SSHKindVRSROS) Commit(s *SSHTransport) (*SSHReply, error) {
	res := s.Run("commit", 30)
	return res, nil
}

func (sk *SSHKindVRSROS) Compare(s *SSHTransport) (*SSHReply, error) {
	var res *SSHReply
	s.Run("validate", 5).Log("", log.WarnLevel)
	res = s.Run("compare /", 5)
	return res, nil
}

func (sk *SSHKindVRSROS) Discard(s *SSHTransport) (*SSHReply, error) {
	res := s.Run("discard /", 5)
	return res, nil
}

func (sk *SSHKindVRSROS) PromptParse(s *SSHTransport, in *string) *SSHReply {
	// SROS MD-CLI \r...prompt
	r := strings.LastIndex(*in, "\r\n\r\n")
	if r > 0 {
		return &SSHReply{
			Result: (*in)[:r],
			Prompt: (*in)[r+4:] + s.PromptChar,
		}
	}
	return nil
}
