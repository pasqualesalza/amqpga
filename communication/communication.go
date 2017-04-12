package communication

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/pasqualesalza/amqpga/config"
	"github.com/pasqualesalza/amqpga/util"
)

// Creates the request queue.
func CreateRequestQueue(channel *amqp.Channel) *amqp.Queue {
	requestQueue, err := channel.QueueDeclare(
		config.RequestQueueName, // name
		true,  // durable
		true,  // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	util.FailOnError(err, "Failed to declare the request queue")
	log.WithFields(log.Fields{
		"queue": config.RequestQueueName,
	}).Info("Request queue created")

	return &requestQueue
}

// Creates the response queue.
func CreateResponseQueue(channel *amqp.Channel) *amqp.Queue {
	responseQueue, err := channel.QueueDeclare(
		config.ResponseQueueName, // name
		true,  // durable
		true,  // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	util.FailOnError(err, "Failed to declare the response queue")
	log.WithFields(log.Fields{
		"queue": config.ResponseQueueName,
	}).Info("Request queue created")

	return &responseQueue
}

// Connects to the server.
func Connect(host string) *amqp.Connection {
	connection, err := amqp.Dial(host)

	util.FailOnError(err, "Failed to connect to RabbitMQ.")
	log.WithFields(log.Fields{
		"rabbitMQHost": host,
	}).Info("Connected to RabbitMQ")

	return connection
}

// Opens a channel.
func OpenChannel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()

	util.FailOnError(err, "Failed to open a channel")
	log.Info("Channel opened")

	return channel
}

// Sets the dispatcher to be fair.
func SetFairDispatch(channel *amqp.Channel) {
	err := channel.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	)

	util.FailOnError(err, "Failed to set the fair dispatch")
	log.Info("Fair dispatch set")
}

// Consumes a queue.
func ConsumeQueue(channel *amqp.Channel, queue *amqp.Queue) <-chan amqp.Delivery {
	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // args
	)

	util.FailOnError(err, fmt.Sprintf("Failed to register as consumer to the %v queue", queue.Name))
	log.WithFields(log.Fields{
		"queue": queue.Name,
	}).Infof("Registered as consumer to the %v queue", queue.Name)

	return messages
}

// Sends a message to the queue.
func PublishMessage(data []byte, channel *amqp.Channel, queue *amqp.Queue) {
	// Publishes the message.
	err := channel.Publish(
		"",         // exchange
		queue.Name, // key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "binary/gob",
			Body:        data})
	util.FailOnError(err, fmt.Sprintf("Failed to publish a message on the %v queue", queue.Name))
}
