package command_test

import (
	chromecast "github.com/popescu-af/go-chromecast"
	"github.com/popescu-af/go-chromecast/command"
)

var _ chromecast.AmpController = command.Launcher{}.AmpController()
