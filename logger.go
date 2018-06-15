package chromecast

// Logger is for structured logging in services
// (like in https://github.com/go-kit/kit/tree/master/log)
type Logger interface {
	Log(keyvals ...interface{}) error
}
