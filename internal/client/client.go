package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Client contains a clients name and net.Conn information
type Client struct {
	Name string
	Mark string `validate:"oneof:X Y"`
	Conn net.Conn
}

func (c Client) String() string {
	return fmt.Sprintf("%s-%s-%s", c.Conn.RemoteAddr(), c.Name, c.Mark)
}

// GetSetName asks the user for a name, and sets it to the client
func (c *Client) GetSetName() error {
	// Ask for a user name for this client
	if err := c.Msg("Enter a username:"); err != nil {
		return errors.Wrap(err, "failed to write 'name' output to client")
	}
	// Read client name input
	name, err := bufio.NewReader(c.Conn).ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "failed to read client input for 'name'")
	}
	c.Name = name[:len(name)-1]
	return nil
}

// Close a client connection
func (c *Client) Close(msg string) error {
	if err := c.Msg(msg); err != nil {
		return errors.Wrap(err, "failed to send close msg to client. failed to close client")
	}
	if err := c.Msg("\\close"); err != nil {
		return errors.Wrap(err, "failed to send close signal to client. failed to client client")
	}
	if err := c.Conn.Close(); err != nil {
		return errors.Wrap(err, "failed to close the connection to the client")
	}
	return nil
}

// Msg writes to the client
func (c *Client) Msg(m string) error {
	for i, msg := range strings.Split(m, "\n") {
		send := []byte(msg)
		send = append(send, []byte("\n")...)
		_, err := c.Conn.Write(send)
		if err != nil {
			return errors.Wrapf(err, "failed to write msg %s to client %s",
				string(msg),
			)
		}
		log.Info().Int("i", i).Str("msg", msg).Str("client", fmt.Sprintf("%s", c)).Bytes("msgB", send).Msg("msg sent")
		time.Sleep(time.Millisecond * 10)
	}
	return nil
}
