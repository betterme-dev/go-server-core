package mq

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net"
	"os"
	"strconv"
	"time"
)

// Client encapsulates a pointer to an amqp.Connection
type Client struct {
	conn *amqp.Connection
}

func NewConnection() (*amqp.Connection, error) {
	rabbitUser := os.Getenv("MQ_USERNAME")
	rabbitPass := os.Getenv("MQ_PASSWORD")
	rabbitHost := os.Getenv("MQ_HOST")
	rabbitPort := os.Getenv("MQ_PORT")
	rabbitVhost := os.Getenv("MQ_VHOST")

	amqpUri := "amqp://" + rabbitUser + ":" + rabbitPass + "@" + rabbitHost + ":" + rabbitPort
	amqpConfig := amqp.Config{
		Vhost: "/" + rabbitVhost,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 2*time.Second) // FIXME move timeout to config
		},
	}
	amqpConn, err := amqp.DialConfig(amqpUri, amqpConfig)
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

func NewQueueForChannel(name string, ch *amqp.Channel) (q amqp.Queue, err error) {
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

	go consumeLoop(msgs, handlerFunc, done)

	return nil
}

func (c *Client) Publish(queueName string, body []byte) error {
	if c.conn == nil {
		return fmt.Errorf("tried to send message before connection was initialized")
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

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery, c chan struct{}), done <-chan bool) {
	maxGoroutines := make(chan struct{}, 1024) // prevent too much of concurrency
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
