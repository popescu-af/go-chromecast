package command

import (
	"fmt"

	chromecast "github.com/oliverpool/go-chromecast"
)

var ErrAppNotFound = fmt.Errorf("app not found in status")

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
	return "", ErrAppNotFound
}
