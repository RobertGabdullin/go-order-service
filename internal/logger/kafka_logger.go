package logger

import (
	"encoding/json"
	"fmt"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/event_broker"
)

type KafkaLogger struct {
	OutputMode  string
	KafkaTopic  string
	KafkaClient event_broker.EventProducer
}

func (l KafkaLogger) Log(e APIEvent) error {
	eventData, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}

	if l.OutputMode == "kafka" {
		if err := l.KafkaClient.ProduceEvent(l.KafkaTopic, string(eventData)); err != nil {
			return fmt.Errorf("failed to produce event to Kafka: %v", err)
		}
		return nil
	}

	fmt.Printf("API Event: %s\n", eventData)
	return nil
}
