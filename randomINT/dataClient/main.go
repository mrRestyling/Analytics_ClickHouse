package main

import (
	rabbitmq "calcServ/dataClient/RMQ"
	"log"
)

func main() {

	conn, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/")
	fail(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := rabbitmq.OpenChannel(conn)
	fail(err, "Failed to open a channel")
	defer ch.Close()

	q, err := rabbitmq.DeclareQueue(ch, "RandomINT")
	fail(err, "Failed to declare a queue")

	msgs, err := rabbitmq.Consume(ch, q)
	fail(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	var forever chan struct{}
	<-forever
}

func fail(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
