package transport

import (
	"strings"
)

// An interface to implement kind specific methods for transactions and prompt checking
// each kind should implement this interface
type SSHKind interface {
	// Start a config transaction
	Start(s *SSHTransport, transaction bool) error
	// Commit a config transaction
	Commit(s *SSHTransport) (*SSHReply, error)
	// Commit a config transaction
	Compare(s *SSHTransport) (*SSHReply, error) // Prompt parsing function
	Discard(s *SSHTransport) (*SSHReply, error) // Prompt parsing function
	//
	// This function receives string, split by the delimiter and should ensure this is a valid prompt
	// Valid prompt, strip the prompt from the result and add it to the prompt in SSHReply
	//
	// A default implementation is promptParseNoSpaces, which simply ensures there are
	// no spaces between the start of the line and the #
	PromptParse(s *SSHTransport, in *string) *SSHReply
}

// This is a helper function to parse the prompt, and can be used by SSHKind's ParsePrompt
// Used in SRL today
func promptParseNoSpaces(in *string, promptChar string, lines int) *SSHReply {
	n := strings.LastIndex(*in, "\n")
	if n < 0 {
		return nil
	}
	if strings.Contains((*in)[n:], " ") {
		return nil
	}
	if lines > 1 {
		// Add another line to the prompt
		res := (*in)[:n]
		n = strings.LastIndex(res, "\n")
	}
	if n < 0 {
		n = 0
	}
	return &SSHReply{
		Result: (*in)[:n],
		Prompt: (*in)[n:] + promptChar,
	}
}
