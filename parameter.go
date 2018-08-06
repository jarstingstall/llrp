package llrp

const ConnectionAttemptEventType = 256

// ParameterHeader provides information about a parameter
type ParameterHeader struct {
	Type   uint16
	Length uint16
}

type ReaderEventNotificationData struct {
	ParameterHeader
	Timestamp              Timestamp
	ConnectionAttemptEvent ConnectionAttemptEvent
}

type Timestamp struct {
	ParameterHeader
	Microseconds uint64
}

type ConnectionAttemptEvent struct {
	ParameterHeader
	Status uint16
}
