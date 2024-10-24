package serviceinterfaces

type IKafkaProducer interface {
	ProduceMessage(topic string, message []byte) error
	ProduceMessagesInBatch(topic string, messages [][]byte) error
	StartDeliveryReportsHandler()
	Close()
	Flush(timeInSeconds uint)
}
