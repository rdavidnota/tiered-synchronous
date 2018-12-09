package mail

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://" + RabbitmqConfig["user"] + ":" + RabbitmqConfig["password"] + "@" +
		RabbitmqConfig["host"] + ":" + RabbitmqConfig["port"] + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		RabbitmqConfig["queue_response"], // name
		true,                             // durable
		false,                            // delete when usused
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		RabbitmqConfig["queue_request"], // queue
		"",                              // consumer
		false,                           // auto-ack
		false,                           // exclusive
		false,                           // no-local
		false,                           // no-wait
		nil,                             // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// n, err := strconv.Atoi(string(d.Body))
			// failOnError(err, "Failed to convert body to integer")
			// log.Printf(" [.] fib(%d)", n)
			// response := fib(n)
			response := "mail"
			fmt.Println(d.ReplyTo)
			fmt.Println(d.CorrelationId)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(response),
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
