package chromecast

// Sender sends a payload (without expecting a reply)
type Sender interface {
	Send(env Envelope, payload interface{}) error
}

// Requester sends a payload and expects one payload as reply
type Requester interface {
	Request(env Envelope, payload IdentifiablePayload) (<-chan []byte, error)
}

// Listener allows to listen to specific messages and forward them (non-blocking) on ch
type Listener interface {
	Listen(env Envelope, responseType string, ch chan<- []byte)
}

// Client interface is too weak
type Client interface {
	Listener
	Sender
	Requester
	Dispatch() error
	Close() error
}
