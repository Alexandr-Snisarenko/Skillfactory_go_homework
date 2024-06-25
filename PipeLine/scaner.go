package main

import (
	"fmt"
	"log/slog"
	"strconv"
)

func StartScaner(ch chan<- int) {
	var (
		inVal string
		val   int
		err   error
	)
	fmt.Println("Для выхода введите 'exit'")
	fmt.Println("Введите целые числа: ")
	for {
		_, err = fmt.Scan(&inVal)
		if err != nil {
			fmt.Println(err)
		}

		if inVal == "exit" {
			slog.Debug("Scaner: получена команда exit", "thread", "main")
			break
		}

		if val, err = strconv.Atoi(inVal); err != nil {
			fmt.Println("Допустимы только целые числа")
			continue
		}

		slog.Debug("Scaner: передача числа в канал данных pipeline ", "thread", "main")
		ch <- val
	}

	return
}
