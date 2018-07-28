package cli

import (
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
)

func ReadStdinKeyPresses(ctx context.Context, out chan<- KeyPress) {
	buf := make(chan []byte, 5)
	go forwardKeyPress(buf, out)
	forwardStdin(ctx, buf)
}

// forwardStdin
// guarantuees that on return (when ctx is Done):
// - out will have been closed
// - the stty will be clean
func forwardStdin(ctx context.Context, out chan<- []byte) {
	//no buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	//no visible output
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// reset on close
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	forwardBytes(ctx, os.Stdin, out)
}

func forwardBytes(ctx context.Context, r io.Reader, out chan<- []byte) {
	var mu sync.Mutex

	go func() {
		var n int
		for {
			b := make([]byte, 10)
			n, _ = r.Read(b)
			mu.Lock()
			if ctx.Err() != nil {
				mu.Unlock()
				return
			}
			// out is guarantueed not to be closed (otherwise ctx.Err != nil)
			select {
			case out <- b[:n]:
			case <-ctx.Done():
				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}()

	<-ctx.Done()
	mu.Lock()
	close(out)
	mu.Unlock()
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

// KeyPress represents a typed key
type KeyPress struct {
	Type KeyType
	Key  byte
}

// KeyType represent a class of keypress
type KeyType int

// KeyTypes
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
