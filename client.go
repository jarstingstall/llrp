package llrp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"unicode/utf8"
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

	msg, err := c.readReaderEventNotifaction()
	if err != nil {
		return err
	}

	if msg.ReaderEventNotificationData.ConnectionAttemptEvent.Status != 0 {
		switch msg.ReaderEventNotificationData.ConnectionAttemptEvent.Status {
		case 2:
			err = errors.New("connection failed, a Client initiated connection already exists")
		case 3:
			err = errors.New("connection failed")
		case 4:
			err = errors.New("connection failed, another connection attempted")
		}
		return err
	}

	fmt.Printf("%+v\n", msg)
	return nil
}

// Close will close the underlying tcp connection and
// return any errors
func (c *Client) Close() error {
	b := []byte{4, CloseConnectionType, 0, 0, 0, 10, 0, 0, 0, 1}
	_, err := c.conn.Write(b)
	if err != nil {
		return err
	}
	res, err := c.readCloseConnectionResponse()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", res)
	return c.conn.Close()
}

func (c *Client) readMessageHeader() (MessageHeader, error) {
	b := make([]byte, 10)
	_, err := c.conn.Read(b)
	if err != nil {
		return MessageHeader{}, err
	}
	return MessageHeader{
		Type:   b[1],
		Length: binary.BigEndian.Uint32(b[2:6]),
		ID:     binary.BigEndian.Uint32(b[6:]),
	}, nil
}

func (c *Client) readParameterHeader() (ParameterHeader, error) {
	b := make([]byte, 4)
	_, err := c.conn.Read(b)
	if err != nil {
		return ParameterHeader{}, err
	}
	return ParameterHeader{
		Type:   binary.BigEndian.Uint16(b[:2]),
		Length: binary.BigEndian.Uint16(b[2:]),
	}, nil
}

func (c *Client) readReaderEventNotifaction() (ReaderEventNotification, error) {
	mh, err := c.readMessageHeader()
	if err != nil {
		return ReaderEventNotification{}, err
	}
	if mh.Type != ReaderEventNotificationType {
		return ReaderEventNotification{}, errors.New("expected ReaderEventNotification message")
	}

	rend, err := c.readReaderEventNotifactionData()
	if err != nil {
		return ReaderEventNotification{}, err
	}

	return ReaderEventNotification{
		MessageHeader:               mh,
		ReaderEventNotificationData: rend,
	}, nil
}

func (c *Client) readReaderEventNotifactionData() (ReaderEventNotificationData, error) {
	ph, err := c.readParameterHeader()
	if err != nil {
		return ReaderEventNotificationData{}, err
	}

	timestamp, err := c.readTimestamp()
	if err != nil {
		return ReaderEventNotificationData{}, err
	}

	rend := ReaderEventNotificationData{
		ParameterHeader: ph,
		Timestamp:       timestamp,
	}
	optionalParamsLength := ph.Length - 16

	for optionalParamsLength > 0 {
		ph, err = c.readParameterHeader()
		if err != nil {
			return rend, err
		}

		switch ph.Type {
		case ConnectionAttemptEventType:
			b := make([]byte, 2)
			_, err = c.conn.Read(b)
			if err != nil {
				return rend, err
			}
			rend.ConnectionAttemptEvent = ConnectionAttemptEvent{
				ParameterHeader: ph,
				Status:          binary.BigEndian.Uint16(b),
			}
		}
		optionalParamsLength -= ph.Length
	}

	return rend, nil
}

func (c *Client) readTimestamp() (Timestamp, error) {
	ph, err := c.readParameterHeader()
	if err != nil {
		return Timestamp{}, err
	}

	b := make([]byte, 8)
	_, err = c.conn.Read(b)
	if err != nil {
		return Timestamp{}, err
	}

	return Timestamp{
		ParameterHeader: ph,
		Microseconds:    binary.BigEndian.Uint64(b),
	}, nil
}

func (c *Client) readCloseConnectionResponse() (CloseConnectionResponse, error) {
	mh, err := c.readMessageHeader()
	if err != nil {
		return CloseConnectionResponse{}, err
	}
	if mh.Type != CloseConnectionResponseType {
		return CloseConnectionResponse{}, errors.New("expected CloseConnectionResponse message")
	}

	status, err := c.readLLRPStatus()
	if err != nil {
		return CloseConnectionResponse{}, err
	}

	return CloseConnectionResponse{
		MessageHeader: mh,
		LLRPStatus:    status,
	}, nil
}

func (c *Client) readLLRPStatus() (LLRPStatus, error) {
	ph, err := c.readParameterHeader()
	if err != nil {
		return LLRPStatus{}, err
	}

	b := make([]byte, 4)
	_, err = c.conn.Read(b)
	if err != nil {
		return LLRPStatus{}, err
	}

	statusCode := binary.BigEndian.Uint16(b[:2])
	edbc := binary.BigEndian.Uint16(b[2:])

	status := LLRPStatus{
		ParameterHeader:           ph,
		StatusCode:                statusCode,
		ErrorDescriptionByteCount: edbc,
	}

	// Success status code
	if statusCode == 0 {
		return status, nil
	}

	if edbc > 0 {
		b := make([]byte, edbc)
		_, err = c.conn.Read(b)
		if err != nil {
			return status, err
		}

		var errDesc []rune
		for len(b) > 0 {
			r, size := utf8.DecodeRune(b)
			errDesc = append(errDesc, r)
			b = b[size:]
		}
		status.ErrorDescription = string(errDesc)
	}

	return status, nil
}
