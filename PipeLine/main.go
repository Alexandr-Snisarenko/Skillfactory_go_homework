package main

import (
	"fmt"
	"time"
)

const bufferSize int = 10
const printBufferDilay = 10 // задержка при печати буфера. в секундах

func PrintBufer(buf *RingBufferInt, done chDone) {
	for {
		// ждём 5 сек перед следующим чтением буфера
		time.Sleep(time.Second * printBufferDilay)
		select {
		case <-done:
			return
		default:
			// читаем текущий буфер
			arr := buf.GetAll()
			// печатаем буфер
			fmt.Println("Out Buffer: ", arr)
		}
	}
}

func main() {
	// создаём входной и сигнальный каналы
	in := make(chan int)
	done := make(chan interface{})

	// формируем массив стейджей для обработки
	stages := make([]Stage, 2)
	stages[0] = NegativeNumberFilter
	stages[1] = Dev3NumberFilter

	// запускаем пайплайн
	bufOut := ExecuteBufPipeline(in, done, bufferSize, stages...)

	// горутина печати буфера
	go PrintBufer(bufOut, done)

	// стартуем сканер (сканер в рутине main)
	StartScaner(in)

	// сворачиваем все рутины (закрываем сигнальный канал) они бы и так все закрылись.
	// но ... для примера done канала
	close(done)
}
