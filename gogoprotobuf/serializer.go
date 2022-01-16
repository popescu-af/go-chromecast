package gogoprotobuf

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"
	chromecast "github.com/popescu-af/go-chromecast"
	"github.com/popescu-af/go-chromecast/gogoprotobuf/pb"
)

type Serializer struct {
	Conn   io.ReadWriter
	Logger chromecast.Logger
	rMu    sync.Mutex
	sMu    sync.Mutex
}

// Receive receives a message
func (s *Serializer) Receive() (env chromecast.Envelope, pay []byte, err error) {
	s.rMu.Lock()
	defer s.rMu.Unlock()

	var length uint32
	err = binary.Read(s.Conn, binary.BigEndian, &length)
	if err != nil {
		return env, pay, fmt.Errorf("failed to read packet length: %s", err)
	}
	if length == 0 {
		return env, pay, fmt.Errorf("empty packet")
	}

	packet := make([]byte, length)
	_, err = io.ReadFull(s.Conn, packet)
	if err != nil {
		return env, pay, fmt.Errorf("failed to read full packet: %s", err)
	}

	cmessage := &pb.CastMessage{}
	err = proto.Unmarshal(packet, cmessage)
	if err != nil {
		return env, pay, fmt.Errorf("failed to unmarshal packet: %s", err)
	}

	env = chromecast.Envelope{
		Source:      *cmessage.SourceId,
		Destination: *cmessage.DestinationId,
		Namespace:   *cmessage.Namespace,
	}

	s.Logger.Log(
		"msg", env.Destination+" ⇐ "+env.Source,
		// "destination", env.Destination,
		// "source", env.Source,
		"namespace", env.Namespace,
		"payload", strings.Replace(*cmessage.PayloadUtf8, `"`, `'`, -1),
	)

	return env, []byte(*cmessage.PayloadUtf8), nil
}

// Send sends a payload
func (s *Serializer) Send(env chromecast.Envelope, pay []byte) error {
	payloadString := string(pay)
	message := &pb.CastMessage{
		ProtocolVersion: pb.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &env.Source,
		DestinationId:   &env.Destination,
		Namespace:       &env.Namespace,
		PayloadType:     pb.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadString,
	}

	proto.SetDefaults(message)

	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %s", err)
	}

	s.Logger.Log(
		"msg", env.Source+" ⇒ "+env.Destination,
		// "source", env.Source,
		// "destination", env.Destination,
		"namespace", env.Namespace,
		"payload", strings.Replace(*message.PayloadUtf8, `"`, `'`, -1),
	)

	s.sMu.Lock()
	defer s.sMu.Unlock()

	err = binary.Write(s.Conn, binary.BigEndian, uint32(len(data)))
	if err != nil {
		return fmt.Errorf("failed to write length: %s", err)
	}
	_, err = s.Conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data: %s", err)
	}

	return nil
}
