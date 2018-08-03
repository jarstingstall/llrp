package llrp

import (
	"encoding/binary"
	"net"
)

const ReaderEventNotificationMsg = 63

// MessageHeader provides information about a message
type MessageHeader struct {
	Type   uint8
	Length uint32
	ID     uint32
}

func readMessageHeader(conn net.Conn) (MessageHeader, error) {
	b := make([]byte, 10)
	_, err := conn.Read(b)
	if err != nil {
		return MessageHeader{}, err
	}
	return MessageHeader{
		Type:   b[1],
		Length: binary.BigEndian.Uint32(b[2:6]),
		ID:     binary.BigEndian.Uint32(b[6:]),
	}, nil
}
