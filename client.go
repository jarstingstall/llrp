package llrp

import (
	"net"
)

// Client used to send/parse messages
type Client struct {
	Host string
	Port string
	conn net.Conn
}

// NewClient returns a new initialized Client or an
// error if unable to make a connection
func NewClient(host, port string) *Client {
	return &Client{
		Host: host,
		Port: port,
	}
}

// Connect establishes a tcp connection with the LLRP server
// and sets it on the Client object for future reads/writes
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.Host+":"+c.Port)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// Close will close the underlying tcp connection and
// return any errors
func (c *Client) Close() error {
	return c.conn.Close()
}
