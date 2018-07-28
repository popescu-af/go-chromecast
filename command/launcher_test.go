package command_test

import (
	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

var _ chromecast.AmpController = command.Launcher{}.AmpController()
