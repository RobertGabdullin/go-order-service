//go:build unit

package logger

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKafkaClient struct {
	mock.Mock
}

func (m *MockKafkaClient) ProduceEvent(topic, message string) error {
	args := m.Called(topic, message)
	return args.Error(0)
}

func TestKafkaLogger_LogToKafkaSuccess(t *testing.T) {
	t.Parallel()
	kafkaClient := new(MockKafkaClient)
	log := KafkaLogger{
		OutputMode:  "kafka",
		KafkaTopic:  "test_topic",
		KafkaClient: kafkaClient,
	}

	event := APIEvent{
		Timestamp:  time.Now(),
		MethodName: "TestMethod",
		RawRequest: "TestRequest",
	}

	eventData, _ := json.Marshal(event)

	kafkaClient.On("ProduceEvent", "test_topic", string(eventData)).Return(nil)

	err := log.Log(event)
	assert.NoError(t, err)

	kafkaClient.AssertCalled(t, "ProduceEvent", "test_topic", string(eventData))
	kafkaClient.AssertExpectations(t)
}

func TestKafkaLogger_LogToKafkaFailure(t *testing.T) {
	t.Parallel()
	kafkaClient := new(MockKafkaClient)
	log := KafkaLogger{
		OutputMode:  "kafka",
		KafkaTopic:  "test_topic",
		KafkaClient: kafkaClient,
	}

	event := APIEvent{
		Timestamp:  time.Now(),
		MethodName: "TestMethod",
		RawRequest: "TestRequest",
	}

	eventData, _ := json.Marshal(event)

	kafkaClient.On("ProduceEvent", "test_topic", string(eventData)).Return(errors.New("failed to produce event"))

	err := log.Log(event)
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to produce event to Kafka: failed to produce event")

	kafkaClient.AssertExpectations(t)
}

func TestKafkaLogger_LogToConsole(t *testing.T) {
	t.Parallel()
	log := KafkaLogger{
		OutputMode: "console",
	}

	event := APIEvent{
		Timestamp:  time.Now(),
		MethodName: "TestMethod",
		RawRequest: "TestRequest",
	}

	err := log.Log(event)
	assert.NoError(t, err)
}
