package discovery

import chromecast "github.com/oliverpool/go-chromecast"

// DeviceMatcher allows to specicy which device should be accepted
type DeviceMatcher func(*chromecast.Device) bool

// WithName matches a device by its name
func WithName(name string) DeviceMatcher {
	return func(device *chromecast.Device) bool {
		return device != nil && device.Name() == name
	}
}

// WithID matches a device by its name
func WithID(id string) DeviceMatcher {
	return func(device *chromecast.Device) bool {
		return device != nil && device.ID() == id
	}
}

func matchAll(matchers ...DeviceMatcher) DeviceMatcher {
	return func(device *chromecast.Device) bool {
		for _, m := range matchers {
			if !m(device) {
				return false
			}
		}
		return true
	}
}
