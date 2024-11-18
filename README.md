# Analytics_ClickHouse


## randomINT
Сбор информации с помощью RabbitMQ и ClickHouse

## Технологии
- RabbitMQ
- ClickHouse
- chi
- Docker

## Описание

Сервер - портал, который дает нам информацию по запросу 

Клиент - делает запрос к серверу и пересылает информацию в RabbitMQ

Консьюмер - обрабатывает сообщения из RabbitMQ и кладет её в ClickHouse


