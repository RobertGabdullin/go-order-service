package api

import (
	"fmt"
	"log"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/logger"
)

func (s *server) logEvent(methodName string, request interface{}) {
	if s.logger == nil {
		return
	}
	event := logger.APIEvent{
		Timestamp:  time.Now(),
		MethodName: methodName,
		RawRequest: fmt.Sprintf("%v", request),
	}

	err := s.logger.Log(event)
	if err != nil {
		log.Printf("Failed to log event: %v", err)
	}
}
