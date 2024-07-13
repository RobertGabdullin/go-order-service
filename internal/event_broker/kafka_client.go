package event_broker

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type KafkaClient struct {
	producer sarama.SyncProducer
}

func NewKafkaClient(brokers []string, config *sarama.Config) (*KafkaClient, error) {
	if config == nil {
		config = sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Retry.Max = 5
		config.Producer.Return.Successes = true
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaClient{producer: producer}, nil
}

func (kc *KafkaClient) ProduceEvent(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := kc.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}
	return nil
}

func (kc *KafkaClient) CloseProducer() {
	if err := kc.producer.Close(); err != nil {
		log.Printf("failed to close Kafka producer: %v", err)
	}
}

func StartConsumer(brokers []string, topic string) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("failed to start consumer for partition: %v", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		fmt.Printf("Consumed message: %s\n", string(message.Value))
	}
}
