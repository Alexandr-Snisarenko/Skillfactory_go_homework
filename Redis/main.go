package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	// импорт драйвера
	"github.com/go-redis/redis/v8"
)

// сущность для хранения в СУБД
type book struct {
	ID    int
	Title string
}

func main() {
	// подключение к СУБД
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // без пароля
		DB:       0,  // БД по умолчанию
	})
	books := []book{
		{ID: 1, Title: "1984"},
		{ID: 2, Title: "Clean Architecture"},
	}
	// выполнение запросов
	err := updateCache(redisClient, books)
	if err != nil {
		log.Fatal(err)
	}
	data, err := cache(redisClient, []int{1, 2})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}

// cache возвращает книги из кэша.
func cache(client *redis.Client, ids []int) ([]book, error) {
	var books []book
	for _, id := range ids {
		cmd := client.Get(context.Background(), "books:"+strconv.Itoa(id))
		var b book
		err := json.Unmarshal([]byte(cmd.Val()), &b)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

// updateCache обновляет данные в кэше Redis.
func updateCache(client *redis.Client, books []book) error {
	for _, b := range books {
		key := "books:" + strconv.Itoa(b.ID)
		val, err := json.Marshal(b)
		if err != nil {
			return err
		}
		err = client.Set(context.Background(), key, string(val), 0).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
