package chromecast

import (
	"fmt"
	"strings"
)

type Launcher interface {
	Launch(appID string) (Status, error)
	Stop() (Status, error)
}

type StatusResponse struct {
	Status *Status `json:"status"`
}

type Status struct {
	Applications []*ApplicationSession `json:"applications"`
	Volume       *Volume               `json:"volume,omitempty"`
}

func (s Status) String() string {
	var str strings.Builder

	if s.Applications != nil {
		if len(s.Applications) == 0 {
			str.WriteString("No application running\n")
		} else {
			str.WriteString(fmt.Sprintf("Running applications: %d\n", len(s.Applications)))
			for _, app := range s.Applications {
				str.WriteString(fmt.Sprintf(" - [%s] %s\n", *app.DisplayName, *app.StatusText))
			}
		}
	}
	if s.Volume != nil {
		str.WriteString(fmt.Sprintf("Volume: %.2f", *s.Volume.Level))
		if *s.Volume.Muted {
			str.WriteString(" (muted)")
		}
	}

	return str.String()
}

type ApplicationSession struct {
	AppID       *string      `json:"appId,omitempty"`
	DisplayName *string      `json:"displayName,omitempty"`
	Namespaces  []*Namespace `json:"namespaces"`
	SessionID   *string      `json:"sessionId,omitempty"`
	StatusText  *string      `json:"statusText,omitempty"`
	TransportId *string      `json:"transportId,omitempty"`
}

type Namespace struct {
	Name string `json:"name"`
}

type Volume struct {
	Level *float64 `json:"level,omitempty"`
	Muted *bool    `json:"muted,omitempty"`
}
