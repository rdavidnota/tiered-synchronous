package documents

import (
	"encoding/json"
	"fmt"
	"github.com/rdavidnota/tiered-synchronous/source/commands/documents"
	"github.com/rdavidnota/tiered-synchronous/source/commands/utils"
	"github.com/rdavidnota/tiered-synchronous/source/domain"
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

			analize(d.Body)

			response := "documents"
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

func analize(request []byte) domain.Result {

	message := domain.RequestBase{Action: ""}
	err := json.Unmarshal(request, &message)

	utils.Check(err)

	if message.Action == "create" {

	} else if message.Action == "delete" {
		return delete(request)
	} else if message.Action == "get" {
		response = GetFile(param)
	} else if message.Action == "download" {
		response
		GetFileById(param)
	} else if message.Action == "list" {
		return list()
	}
}

func create(request []byte){
	requestCreate := domain.RequestCreateDocument{}
	err := json.Unmarshal(request, &requestCreate)
	utils.Check(err)


	documents.CreatedFile(requestCreate.Name, )

	return list()
}

func list() domain.Result {
	listDocument := documents.ListFiles()
	jsonResult, err := json.Marshal(listDocument)

	utils.Check(err)

	return domain.Result{
		Code:    0,
		Message: "OK",
		Json:    string(jsonResult),
	}
}

func delete(request []byte) domain.Result {
	requestDelete := domain.RequestDeleteDocument{}
	err := json.Unmarshal(request, &requestDelete)
	utils.Check(err)
	documents.DeleteFile(requestDelete.ID)

	return list()
}
