package cli

import (
	"fmt"
	"io"

	"github.com/oliverpool/go-chromecast"
)

func FprintStatus(out io.Writer, status cast.Status) {
	if status.Applications != nil {
		if len(status.Applications) == 0 {
			fmt.Fprintln(out, "No application running")
		} else {
			fmt.Fprintf(out, "Running applications: %d\n", len(status.Applications))
			for _, app := range status.Applications {
				fmt.Fprintf(out, " - [%s] %s\n", *app.DisplayName, *app.StatusText)
			}
		}
	}
	if status.Volume != nil {
		fmt.Fprintf(out, "Volume: %.2f", *status.Volume.Level)
		if *status.Volume.Muted {
			fmt.Fprint(out, " (muted)\n")
		} else {
			fmt.Fprint(out, "\n")
		}
	}
}
