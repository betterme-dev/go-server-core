package mq

import (
	"fmt"
	"os"
	"strconv"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

const (
	defaultMaxGoroutinesCount int = 1024
	defaultConnectTimeout     int = 2000 // in milliseconds
)

// Client encapsulates a pointer to an rabbitmq.Connection
type Client struct {
	conn *rabbitmq.Connection
}

func NewConnection() (*rabbitmq.Connection, error) {
	user := viper.GetString("MQ_USERNAME")
	pass := viper.GetString("MQ_PASSWORD")
	host := viper.GetString("MQ_HOST")
	port := viper.GetString("MQ_PORT")
	vhost := viper.GetString("MQ_VHOST")
	connectTimeout := viper.GetString("MQ_CONNECT_TIMEOUT")
	if connectTimeout == "" {
		connectTimeout = string(rune(defaultConnectTimeout))
	}

	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s//%s?connection_timeout=%s",
		user,
		pass,
		host,
		port,
		vhost,
		connectTimeout,
	)
	amqpConn, err := rabbitmq.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	return amqpConn, nil
}

func NewClient() (*Client, error) {
	amqpConn, err := NewConnection()
	if err != nil {
		return nil, err
	}

	return &Client{conn: amqpConn}, nil
}

func NewQueueForChannel(name string, ch *rabbitmq.Channel) (q amqp.Queue, err error) {
	args := make(amqp.Table)
	args["x-max-priority"] = int64(10)

	return ch.QueueDeclare(
		name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)
}

func (c *Client) Subscribe(queueName string, handlerFunc func(amqp.Delivery, chan struct{}), done chan bool) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}

	log.Infof("Declaring Queue (%s)", queueName)
	queue, err := NewQueueForChannel(queueName, ch)
	if err != nil {
		return errors.Wrap(err, "failed to register a queue")
	}

	// Set prefetch count if configured
	prefetchCount, err := strconv.Atoi(os.Getenv("MQ_PREFETCH"))
	if err == nil {
		err = ch.Qos(prefetchCount, 0, false)
		if err != nil {
			return errors.Wrap(err, "failed to set channel QoS")
		}
	} else {
		log.Warnf("failed to read QoS prefetch value: %s", err)
	}

	// Set max goroutines count if configured
	maxGoroutinesCount := viper.GetInt("MQ_MAX_GOROUTINES_COUNT")
	if maxGoroutinesCount == 0 {
		maxGoroutinesCount = defaultMaxGoroutinesCount
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to register a consumer")
	}

	go consumeLoop(msgs, handlerFunc, done, maxGoroutinesCount)

	return nil
}

func (c *Client) Publish(queueName string, body []byte) error {
	if c.conn == nil {
		return errors.New("tried to send message before connection was initialized")
	}
	ch, err := c.conn.Channel() // Get a channel from the connection
	if err != nil {
		return errors.Wrap(err, "failed to get channel")
	}
	defer func() {
		_ = ch.Close()
	}()

	queue, err := NewQueueForChannel(queueName, ch)
	if err != nil {
		return errors.Wrap(err, "failed to get channel queue")
	}

	// TODO try create custom exchange

	// Publishes a message onto the queue.
	err = ch.Publish(
		"",         // use the default exchange
		queue.Name, // routing key, e.g. our queue name
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body, // Our JSON body as []byte
		})
	log.Infof("A message was sent to queue %s: %s", queueName, body)

	return err
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery, c chan struct{}), done <-chan bool, maxGoroutinesCount int) {
	maxGoroutines := make(chan struct{}, maxGoroutinesCount) // prevent too much of concurrency
	for {
		select {
		case d, ok := <-deliveries:
			if ok {
				maxGoroutines <- struct{}{} // acquire lock
				// Invoke the handlerFunc func we passed as parameter.
				go handlerFunc(d, maxGoroutines)
			}
		case <-done:
			return
		}
	}
}

func (c *Client) IsConnectionClosed() bool {
	return c.conn.Connection.IsClosed()
}
