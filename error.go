package chromecast

// ErrorString represents some constant errors
type ErrorString string

// Error returns the underlying error string
func (e ErrorString) Error() string {
	return string(e)
}

// ErrAppNotFound is returned when a given app was not found
const ErrAppNotFound = ErrorString("app not found")
