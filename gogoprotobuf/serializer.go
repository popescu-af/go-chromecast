package gogoprotobuf

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/gogo/protobuf/proto"
	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/gogoprotobuf/pb"
	"github.com/oliverpool/go-chromecast/log"
)

type Serializer struct {
	Conn io.ReadWriter
	rMu  sync.Mutex
	sMu  sync.Mutex
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

	log.Printf("%s ⇐ %s [%s]: %+v",
		env.Destination, env.Source, env.Namespace, *cmessage.PayloadUtf8)

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

	log.Printf("%s ⇒ %s [%s]: %s", env.Source, env.Destination, env.Namespace, *message.PayloadUtf8)

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
