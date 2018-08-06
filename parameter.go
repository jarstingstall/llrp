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

type LLRPStatus struct {
	ParameterHeader
	StatusCode                uint16
	ErrorDescriptionByteCount uint16
	ErrorDescription          string
	FieldError                FieldError
	ParameterError            ParameterError
}

type FieldError struct {
	ParameterHeader
	FieldNum  uint16
	ErrorCode uint16
}

type ParameterError struct {
	ParameterHeader
	ParmeterType uint16
	ErrorCode    uint16
}
