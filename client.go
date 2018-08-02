package llrp

// Client used to connect and send/parse messages
type Client struct {
	Host string
	Port string
}

// NewClient return a new initialized LLRP client
func NewClient(host, port string) *Client {
	return &Client{
		Host: host,
		Port: port,
	}
}
