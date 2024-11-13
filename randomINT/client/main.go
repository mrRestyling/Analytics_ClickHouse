package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Data struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

func main() {

	url := "http://localhost:8080/"
	data := Data{Num1: 665, Num2: 1}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"RandomINT", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Qos(2, 0, false)
	if err != nil {
		log.Fatal("err")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		for {
			if err := Post(ctx, url, data, ch, q.Name); err != nil {
				log.Println("Ошибка отправки запроса:", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	select {}

}

func Post(ctx context.Context, url string, data Data, ch *amqp.Channel, queueName string) error {

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Ошибка конвертации в json", err)
	}

	// Создание запроса на сервер
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Ошибка создания запроса", err)
	}

	// Установка заголовка Content-Type в application/json
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	for {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Сервер недоступен")
			time.Sleep(5 * time.Second)
			continue
		}

		// Чтение ответа от сервера
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ошибка чтения ответа: %v", err)
		}

		// Decode the response from JSON
		var response struct {
			Sum int `json:"sum"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("ошибка декодирования ответа: %v", err)
		}

		// Send the response to the RabbitMQ queue
		err = ch.PublishWithContext(ctx,
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(fmt.Sprintf("%d", response.Sum)),
			},
		)
		if err != nil {
			return fmt.Errorf("ошибка отправки ответа в RabbitMQ: %v", err)
		}

		fmt.Println("Ответ от сервера:", response.Sum)
		return nil
	}
}
