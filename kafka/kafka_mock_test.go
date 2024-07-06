package kafka

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/mock"
)

type MockSyncProducer struct {
	mock.Mock
}

func (m *MockSyncProducer) AbortTxn() error {
	panic("unimplemented")
}

func (m *MockSyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	panic("unimplemented")
}

func (m *MockSyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	panic("unimplemented")
}

func (m *MockSyncProducer) BeginTxn() error {
	panic("unimplemented")
}

func (m *MockSyncProducer) CommitTxn() error {
	panic("unimplemented")
}

func (m *MockSyncProducer) IsTransactional() bool {
	panic("unimplemented")
}

func (m *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	panic("unimplemented")
}

func (m *MockSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	panic("unimplemented")
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	args := m.Called(msg)
	return args.Get(0).(int32), args.Get(1).(int64), args.Error(2)
}

func (m *MockSyncProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}
