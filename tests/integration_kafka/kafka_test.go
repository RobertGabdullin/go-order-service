//go:build integration

package integration_kafka

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/event_broker"
)

type KafkaIntegrationSuite struct {
	suite.Suite
	client *event_broker.KafkaClient
}

func (s *KafkaIntegrationSuite) SetupSuite() {
	client, err := event_broker.NewKafkaClient([]string{"127.0.0.1:9093"}, nil)
	s.Require().NoError(err)
	s.client = client
}

func (s *KafkaIntegrationSuite) TearDownSuite() {
	s.client.CloseProducer()
}

func (s *KafkaIntegrationSuite) TestProduceAndConsume() {
	topic := "integration_test_topic"
	message := "integration_test_message"

	err := s.client.ProduceEvent(topic, message)
	s.NoError(err)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9093"}, config)
	s.Require().NoError(err)
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	s.Require().NoError(err)
	defer partitionConsumer.Close()

	select {
	case msg := <-partitionConsumer.Messages():
		s.Equal(message, string(msg.Value))
	case <-time.After(10 * time.Second):
		s.Fail("Did not receive message in time")
	}
}

func TestKafkaIntegrationSuite(t *testing.T) {
	suite.Run(t, new(KafkaIntegrationSuite))
}
