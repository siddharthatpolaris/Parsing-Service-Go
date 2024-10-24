package serviceinterfaces

type IKafkaConsumerFactory interface {
	CreateConsumer(consumerGroupID string) (IKafkaConsumer, error)
}