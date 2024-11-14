package main

import (
	chDB "calcServ/dataClient/CH"
	rabbitmq "calcServ/dataClient/RMQ"
	"context"
	"log"
)

func main() {

	connDB, err := chDB.ConnectDB()
	fail(err, "Failed to connect to ClickHouse")

	db := chDB.New(connDB)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connRMQ, err := rabbitmq.Dial("amqp://guest:guest@localhost:5672/")
	fail(err, "Failed to connect to RabbitMQ")
	defer connRMQ.Close()

	ch, err := rabbitmq.OpenChannel(connRMQ)
	fail(err, "Failed to open a channel")
	defer ch.Close()

	q, err := rabbitmq.DeclareQueue(ch, "RandomINT")
	fail(err, "Failed to declare a queue")

	msgs, err := rabbitmq.Consume(ch, q)
	fail(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err = db.Push(ctx, d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	var forever chan struct{}
	<-forever

	err = db.Db.Close()
	fail(err, "Failed to close connection to ClickHouse")
}

func fail(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
