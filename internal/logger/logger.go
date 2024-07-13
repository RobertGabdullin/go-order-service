package logger

import (
	"time"
)

type APIEvent struct {
	Timestamp  time.Time `json:"timestamp"`
	MethodName string    `json:"method_name"`
	RawRequest string    `json:"raw_request"`
}

type Logger interface {
	Log(APIEvent) error
}
