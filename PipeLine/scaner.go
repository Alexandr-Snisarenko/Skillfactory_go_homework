package main

import (
	"fmt"
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
			break
		}

		if val, err = strconv.Atoi(inVal); err != nil {
			fmt.Println("Допустимы только целые числа")
			continue
		}

		ch <- val
	}

	return
}
