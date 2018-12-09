package main

import (
	"encoding/json"
	"fmt"
	"github.com/rdavidnota/tiered-synchronous/source/rpc/documents"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://" + documents.RabbitmqConfig["user"] + ":" + documents.RabbitmqConfig["password"] + "@" +
		documents.RabbitmqConfig["host"] + ":" + documents.RabbitmqConfig["port"] + "/")
	documents.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	documents.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		documents.RabbitmqConfig["queue_response"], // name
		true,                             // durable
		false,                            // remove when usused
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	documents.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	documents.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		documents.RabbitmqConfig["queue_request"], // queue
		"",                              // consumer
		false,                           // auto-ack
		false,                           // exclusive
		false,                           // no-local
		false,                           // no-wait
		nil,                             // args
	)
	documents.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			result := documents.Analyze(d.Body)
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
			documents.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
