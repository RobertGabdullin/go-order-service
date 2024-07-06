//go:build unit

package kafka

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestProduceEvent(t *testing.T) {
	mockProducer := new(MockSyncProducer)
	producer = mockProducer

	topic := "test_topic"
	message := "test_message"

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	mockProducer.On("SendMessage", msg).Return(int32(0), int64(0), nil)

	err := ProduceEvent(topic, message)
	assert.NoError(t, err)
	mockProducer.AssertCalled(t, "SendMessage", msg)
	mockProducer.AssertExpectations(t)
}

func TestCloseProducer(t *testing.T) {
	mockProducer := new(MockSyncProducer)
	producer = mockProducer

	mockProducer.On("Close").Return(nil)

	CloseProducer()
	mockProducer.AssertCalled(t, "Close")
	mockProducer.AssertExpectations(t)
}
