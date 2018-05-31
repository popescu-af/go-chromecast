package zeroconf_test

import (
	"github.com/oliverpool/go-chromecast/discovery"
	"github.com/oliverpool/go-chromecast/discovery/zeroconf"
)

// Ensure interface is satisfied
var _ discovery.Scanner = zeroconf.Scanner{}
