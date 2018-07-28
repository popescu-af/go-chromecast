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
		var n int
		for {
			b := make([]byte, 10)
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
	go forwardKeyPress(buf, out)
	return readStdin(buf, done)
}

func forwardKeyPress(in <-chan []byte, out chan<- KeyPress) {
	for b := range in {
		kp := bytesToKeyPress(b)
		if kp.Type != Unsupported {
			out <- kp
		} else {
			for _, by := range b {
				out <- KeyPress{
					Type: Unsupported,
					Key:  by,
				}
			}
		}
	}
	close(out)
}

func bytesToKeyPress(b []byte) KeyPress {
	// 1 byte
	if len(b) == 1 {
		by := b[0]
		switch {
		case 'a' <= by && by <= 'z':
			return KeyPress{
				Type: LowerCaseLetter,
				Key:  by,
			}
		case 'A' <= by && by <= 'Z':
			return KeyPress{
				Type: UpperCaseLetter,
				Key:  by,
			}
		case by == ' ':
			return KeyPress{
				Type: SpaceBar,
				Key:  by,
			}
		case by == 27:
			// escape
			return KeyPress{
				Type: Escape,
				Key:  by,
			}
		}
	}

	// 3 bytes
	if len(b) == 3 && b[0] == 27 && b[1] == 91 {
		// arrow (see const Up Left...)
		return KeyPress{
			Type: Arrow,
			Key:  b[2],
		}
	}
	return KeyPress{
		Type: Unsupported,
	}
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

// Arrows
const (
	Up byte = 65 + iota
	Down
	Right
	Left
)
