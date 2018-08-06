package llrp

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
