package config

const (
	RequestQueueName               = "amqpga_request"
	ResponseQueueName              = "amqpga_response"
	EtcdExperimentConfigurationKey = "/services/amqpga/experiment"
	EtcdRabbitMQConfigurationKey   = "/services/rabbitmq"
	EtcdMongoDBConfigurationKey    = "/services/mongodb"
)
