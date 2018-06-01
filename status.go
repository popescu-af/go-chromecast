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

func (st Status) String() string {
	var str strings.Builder

	if st.Applications != nil {
		if len(st.Applications) == 0 {
			str.WriteString("No application running\n")
		} else {
			str.WriteString(fmt.Sprintf("Running applications: %d\n", len(st.Applications)))
			for _, app := range st.Applications {
				str.WriteString(fmt.Sprintf(" - [%s] %s\n", *app.DisplayName, *app.StatusText))
			}
		}
	}
	if st.Volume != nil {
		str.WriteString(fmt.Sprintf("Volume: %.2f", *st.Volume.Level))
		if *st.Volume.Muted {
			str.WriteString(" (muted)")
		}
	}

	return str.String()
}

func (st Status) AppSupporting(namespace string) (apps []ApplicationSession) {
	for _, app := range st.Applications {
		if app == nil {
			continue
		}
		for _, ns := range app.Namespaces {
			if ns == nil || ns.Name != namespace {
				continue
			}
			apps = append(apps, *app)
		}
	}
	return apps
}

func (st Status) AppWithID(id string) *ApplicationSession {
	for _, app := range st.Applications {
		if app == nil {
			continue
		}
		if app.AppID != nil && *app.AppID == id {
			return app
		}
	}
	return nil
}

func (st Status) FirstDestinationSupporting(namespace string) (string, error) {
	apps := st.AppSupporting(namespace)
	for _, app := range apps {
		if app.TransportId != nil {
			return *app.TransportId, nil
		}
	}
	return "", ErrAppNotFound
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
