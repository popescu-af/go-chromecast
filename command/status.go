package command

import (
	"fmt"

	chromecast "github.com/oliverpool/go-chromecast"
)

func TransportForNamespace(st chromecast.Status, namespace string) (transport string, err error) {
	for _, app := range st.Applications {
		if app == nil || app.TransportId == nil {
			continue
		}
		for _, ns := range app.Namespaces {
			if ns == nil || ns.Name != namespace {
				continue
			}
			return *app.TransportId, nil
		}
	}
	return "", fmt.Errorf("no app with namespace '%s' could be found in status", namespace)
}
