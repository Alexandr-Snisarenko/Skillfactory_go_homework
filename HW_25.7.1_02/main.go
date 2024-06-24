package main

import (
	"fmt"
	"log"
)

func main() {
	n := ""
	fmt.Print("Введите целое число: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Вы ввели число: %s\n", n)
}
