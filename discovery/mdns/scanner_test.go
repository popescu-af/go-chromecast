package mdns_test

import (
	"github.com/oliverpool/go-chromecast/discovery"
	"github.com/oliverpool/go-chromecast/discovery/mdns"
)

// Ensure interface is satisfied
var _ discovery.Scanner = mdns.Scanner{}
