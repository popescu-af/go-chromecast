package cast

import (
	"errors"
	"net"
	"time"

	"golang.org/x/net/context"

	"github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/log"
	castnet "github.com/barnybug/go-cast/net"
)

type Client struct {
	conn       *castnet.Connection
	ctx        context.Context
	cancel     context.CancelFunc
	heartbeat  *controllers.HeartbeatController
	connection *controllers.ConnectionController
	receiver   *controllers.ReceiverController
	media      *controllers.MediaController
}

const DefaultSender = "sender-0"
const DefaultReceiver = "receiver-0"

func NewClient() *Client {
	return &Client{ctx: context.Background()}
}

func (c *Client) Connect(host net.IP, port int) error {
	c.conn = castnet.NewConnection()
	err := c.conn.Connect(host, port)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(c.ctx)
	c.cancel = cancel

	// connect channel
	c.connection = controllers.NewConnectionController(c.conn, DefaultSender, DefaultReceiver)
	c.connection.Connect()

	// start heartbeat
	c.heartbeat = controllers.NewHeartbeatController(c.conn, DefaultSender, DefaultReceiver)
	c.heartbeat.Start(ctx)

	// start receiver
	c.receiver = controllers.NewReceiverController(c.conn, DefaultSender, DefaultReceiver)

	return nil
}

func (c *Client) NewChannel(sourceId, destinationId, namespace string) *castnet.Channel {
	return c.conn.NewChannel(sourceId, destinationId, namespace)
}

func (c *Client) Close() {
	c.cancel()
	c.conn.Close()
}

func (c *Client) Receiver() *controllers.ReceiverController {
	return c.receiver
}

func (c *Client) launchMediaApp() (string, error) {
	// get transport id
	status, err := c.receiver.GetStatus(5 * time.Second)
	if err != nil {
		return "", err
	}
	app := status.GetSessionByAppId(AppMedia)
	if app == nil {
		// needs launching
		status, err = c.receiver.LaunchApp(AppMedia, 5*time.Second)
		if err != nil {
			return "", err
		}
		app = status.GetSessionByAppId(AppMedia)
	}

	if app == nil {
		return "", errors.New("Failed to get media transport")
	}
	return *app.TransportId, nil
}

func (c *Client) IsPlaying() bool {
	status, err := c.receiver.GetStatus(5 * time.Second)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	app := status.GetSessionByAppId(AppMedia)
	if app == nil {
		return false
	}
	if *app.StatusText == "Ready To Cast" {
		return false
	}
	return true
}

func (c *Client) Media() (*controllers.MediaController, error) {
	if c.media == nil {
		transportId, err := c.launchMediaApp()
		if err != nil {
			return nil, err
		}
		conn := controllers.NewConnectionController(c.conn, DefaultSender, transportId)
		if err := conn.Connect(); err != nil {
			return nil, err
		}
		c.media = controllers.NewMediaController(c.conn, DefaultSender, transportId)
		if _, err := c.media.GetStatus(5 * time.Second); err != nil {
			return nil, err
		}
	}
	return c.media, nil
}
