// Copyright 2018 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package driver

import (
	"github.com/streadway/amqp"
)

// RabbitMQ type
type RabbitMQ struct {
	Connection *amqp.Connection
	Host       string
}

// NewRabbitMQ creates an instance of RabbitMQ
func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{}
}

// WithHost define the host
func (r *RabbitMQ) WithHost(host string) *RabbitMQ {
	r.Host = host

	return r
}

// Dial connects to RabbitMQ
func (r *RabbitMQ) Dial() error {
	var err error

	r.Connection, err = amqp.Dial(r.Host)

	if err != nil {
		return err
	}

	return nil
}

// Send sends a message
func (r *RabbitMQ) Send(queue, routingKey, message string) error {
	ch, err := r.Connection.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		queue,    // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	err = ch.Publish(
		queue,      // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// Consume consumes a queue
func (r *RabbitMQ) Consume(queue, routingKey string, callback func(msg string)) error {
	ch, err := r.Connection.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		queue,    // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		queue,      // exchange
		false,
		nil,
	)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		return err
	}

	for d := range msgs {
		callback(string(d.Body))
	}

	return nil
}

// Close closes the connection
func (r *RabbitMQ) Close() {
	r.Connection.Close()
}
