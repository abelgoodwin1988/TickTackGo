package client

import (
	"bufio"
	"net"

	"github.com/pkg/errors"
)

// Client contains a clients name and net.Conn information
type Client struct {
	Name string
	Conn net.Conn
}

// GetSetName asks the user for a name, and sets it to the client
func (c *Client) GetSetName() error {
	// Ask for a user name for this client
	_, err := c.Conn.Write([]byte("Enter a username:\n"))
	if err != nil {
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
func (c *Client) Close() {
	c.Conn.Write([]byte("connection closing from server\n"))
	c.Conn.Write([]byte("\\close"))
	c.Conn.Close()
}
