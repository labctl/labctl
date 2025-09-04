package colorize

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"

	"github.com/charmbracelet/log"
)

func Spawn(colors []*Colorize, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	var oldState *term.State

	defer func() {
		// Cleanup signals when done.
		signal.Stop(ch)
		close(ch)
		// Make sure to close the pty at the end.
		// err = ptmx.Close()
		// Restore old state
		if oldState != nil {
			_ = term.Restore(int(os.Stdin.Fd()), oldState)
		}
	}()

	// Set stdin in raw mode.
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	// Use a channel to copy std to pty
	stdInCh := make(chan string)
	go func() {
		buf := make([]byte, 50)
		for {
			nr, err := os.Stdin.Read(buf)
			if err != nil {
				close(stdInCh)
				log.Errorf("Error in read string: %s", err)
			}
			stdInCh <- string(buf[:nr])
		}
	}()
	go func() {
		ecnt := 10
		var err error
		for in := range stdInCh {
			_, err = ptmx.WriteString(in)
			if err != nil && ecnt > 0 {
				ecnt -= 1
				log.Errorf("%d: %s", ecnt, err)
			}
		}
	}()

	// Copy the pty to Stdout and in the process add some color
	_, err = CopyReplace(os.Stdout, ptmx, colors)
	if err != nil {
		log.Debugf("done. %s\n", err)
	}

	close(stdInCh)

	return nil
}

// errInvalidWrite means that a write returned an impossible count.
var errInvalidWrite = errors.New("invalid write result")

func CopyReplace(dst io.Writer, src io.Reader, regexes []*Colorize) (written int64, err error) {
	var buf []byte
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			buf2 := buf[0:nr]
			for _, rex := range regexes {
				buf2 = rex.Regexp.ReplaceAll(buf2, rex.Replace)
			}
			nw, ew := dst.Write(buf2)
			if nw > nr {
				nw = nr
			}
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
