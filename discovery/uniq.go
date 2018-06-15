package discovery

import chromecast "github.com/oliverpool/go-chromecast"

// Uniq forward all client deduplicated
func Uniq(in <-chan *chromecast.Device, out chan<- *chromecast.Device) {
	seen := make(map[string]struct{})
	for c := range in {
		if c == nil {
			continue
		}
		if _, ok := seen[c.ID()]; ok {
			continue
		}
		out <- c
		seen[c.ID()] = struct{}{}
	}
	close(out)
}
