package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Data struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

func main() {

	url := "http://localhost:8080/"
	data := Data{Num1: 665, Num2: 1}

	for {

		Post(url, data)
		time.Sleep(1 * time.Second)

	}

}

func Post(url string, data Data) {

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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Ошибка отправки запроса", err)
	}

	// Чтение ответа от сервера
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Декодирование ответа из JSON
	var response struct {
		Sum int `json:"sum"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Вывод ответа от сервера
	fmt.Println("Ответ от сервера:", response.Sum)

}
