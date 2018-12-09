package main

import (
	"encoding/json"
	"fmt"
	"github.com/rdavidnota/tiered-synchronous/source/rpc/mail"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://" + mail.RabbitmqConfig["user"] + ":" + mail.RabbitmqConfig["password"] + "@" +
		mail.RabbitmqConfig["host"] + ":" + mail.RabbitmqConfig["port"] + "/")
	mail.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	mail.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		mail.RabbitmqConfig["queue_response"], // name
		true,                                  // durable
		false,                                 // delete when usused
		false,                                 // exclusive
		false,                                 // no-wait
		nil,                                   // arguments
	)
	mail.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	mail.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		mail.RabbitmqConfig["queue_request"], // queue
		"",                                   // consumer
		false,                                // auto-ack
		false,                                // exclusive
		false,                                // no-local
		false,                                // no-wait
		nil,                                  // args
	)
	mail.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// n, err := strconv.Atoi(string(d.Body))
			// failOnError(err, "Failed to convert body to integer")
			// log.Printf(" [.] fib(%d)", n)
			// response := fib(n)
			result := mail.Analyze(d.Body)
			response, _ := json.Marshal(result)

			fmt.Println(d.ReplyTo)
			fmt.Println(d.CorrelationId)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          response,
				})
			mail.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
