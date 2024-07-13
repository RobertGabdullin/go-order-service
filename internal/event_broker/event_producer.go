package event_broker

type EventProducer interface {
	ProduceEvent(topic, message string) error
}
