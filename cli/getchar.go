package cli

import (
	"os"
	"os/exec"
)

func readStdin(out chan<- []byte, done <-chan struct{}) func() {
	//no buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	//no visible output
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	go func() {
		b := make([]byte, 10)
		var n int
		for {
			n, _ = os.Stdin.Read(b)
			select {
			case <-done:
				return
			default:
				out <- b[:n]
			}
		}
	}()

	return func() {
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		<-done
		close(out)
	}
}

func ReadStdinKeys(out chan<- KeyPress, done <-chan struct{}) func() {
	buf := make(chan []byte, 5)
	go func() {
		defer close(out)
		for {
			b, ok := <-buf
			if !ok {
				return
			}
			if len(b) == 1 {
				by := b[0]
				switch {
				case 'a' <= by && by <= 'z':
					out <- KeyPress{
						Type: LowerCaseLetter,
						Key:  by,
					}
				case 'A' <= by && by <= 'Z':
					out <- KeyPress{
						Type: UpperCaseLetter,
						Key:  by,
					}
				case by == ' ':
					out <- KeyPress{
						Type: SpaceBar,
						Key:  by,
					}
				case by == 27:
					// escape
					out <- KeyPress{
						Type: Escape,
						Key:  by,
					}
				default:
					out <- KeyPress{
						Type: Unsupported,
						Key:  by,
					}
				}
				continue
			}

			// multiple bytes
			if len(b) == 3 && b[0] == 27 && b[1] == 91 {
				// arrow
				out <- KeyPress{
					Type: Arrow,
					Key:  b[2],
				}
				continue
			}
			for _, by := range b {
				out <- KeyPress{
					Type: Unsupported,
					Key:  by,
				}
			}
		}

	}()
	return readStdin(buf, done)
}

type KeyPress struct {
	Type KeyType
	Key  byte
}

type KeyType int

const (
	LowerCaseLetter KeyType = iota
	UpperCaseLetter
	Arrow
	SpaceBar
	Escape
	Unsupported
)

const (
	ArrowUp byte = 65 + iota
	ArrowDown
	ArrowRight
	ArrowLeft
)
