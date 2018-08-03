package llrp

import (
	"encoding/binary"
	"net"
)

const ReaderEventNotificationDataParam = 246

// ParameterHeader provides information about a parameter
type ParameterHeader struct {
	Type   uint16
	Length uint16
}

func readParameterHeader(conn net.Conn) (ParameterHeader, error) {
	b := make([]byte, 4)
	_, err := conn.Read(b)
	if err != nil {
		return ParameterHeader{}, err
	}
	return ParameterHeader{
		Type:   binary.BigEndian.Uint16(b[:2]),
		Length: binary.BigEndian.Uint16(b[2:]),
	}, nil
}
