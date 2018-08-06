package llrp

const CloseConnectionResponseType = 4
const CloseConnectionType = 14
const ReaderEventNotificationType = 63

// MessageHeader provides information about a message
type MessageHeader struct {
	Type   uint8
	Length uint32
	ID     uint32
}

type ReaderEventNotification struct {
	MessageHeader
	ReaderEventNotificationData ReaderEventNotificationData
}

type CloseConnection []byte

func newCloseConnection() []byte {
	return []byte{4, CloseConnectionType, 0, 0, 0, 10, 0, 0, 0, 1}
}

type CloseConnectionResponse struct {
	MessageHeader
	LLRPStatus LLRPStatus
}
