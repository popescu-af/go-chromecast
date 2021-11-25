package zeroconf_test

import (
	"github.com/popescu-af/go-chromecast/discovery"
	"github.com/popescu-af/go-chromecast/discovery/zeroconf"
)

// Ensure interface is satisfied
var _ discovery.Scanner = zeroconf.Scanner{}
