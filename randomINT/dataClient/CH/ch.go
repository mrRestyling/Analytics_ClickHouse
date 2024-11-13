package ch

import (
	"context"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type StorageCH struct {
	Db clickhouse.Conn
}

func New(db clickhouse.Conn) *StorageCH {
	return &StorageCH{Db: db}
}

func ConnectDB() (clickhouse.Conn, error) {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:19000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "secret",
		},
		// TLS: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },
	})
	if err != nil {
		return nil, err
	}
	log.Println("соединение с ClickHouse успешно установлено")

	// Проверка соединения
	err = conn.Ping(context.Background())
	if err != nil {
		log.Println("ошибка соединения с ClickHouse", err)
		return nil, err
	}

	log.Println("соединение с ClickHouse проверено")

	err = conn.Exec(context.Background(), "INSERT INTO numbers (number) VALUES (666)")
	if err != nil {
		log.Println("ПРОВЕРКА данные не вставленны", err)
	}

	return conn, nil
}

func (s *StorageCH) Push(ctx context.Context, data []byte) error {

	strResult := string(data)

	query := "INSERT INTO numbers (number) VALUES (?)"
	err := s.Db.Exec(ctx, query, strResult)

	if err != nil {
		log.Fatal("ошибка вставки в CH", err)
		return err
	}
	return nil
}
