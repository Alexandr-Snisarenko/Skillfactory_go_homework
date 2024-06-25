package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

const bufferSize int = 10
const printBufferDilay = 10 // задержка при печати буфера. в секундах

func PrintBufer(buf *RingBufferInt, done chDone) {
	for {
		// ждём 5 сек перед следующим чтением буфера
		slog.Debug("Wait 5s.", "thread", "printBufer")
		time.Sleep(time.Second * printBufferDilay)
		select {
		case <-done:
			slog.Debug("Get Done message", "thread", "printBufer")
			return
		default:
			// читаем текущий буфер
			slog.Debug("Read Bufer", "thread", "printBufer")
			arr := buf.GetAll()
			// печатаем буфер
			slog.Debug("Print Bufer", "thread", "printBufer", "value", arr)
			fmt.Println("Out Buffer: ", arr)
		}
	}
}

func main() {
	// создаём входной и сигнальный каналы
	in := make(chan int)
	done := make(chan interface{})
	// настраиваем логер
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(log)

	// формируем массив стейджей для обработки
	slog.Debug("Create stage collection")
	stages := make([]Stage, 2)
	slog.Debug("Add stage to collection", "thread", "main", "index", 0, "Stage", "NegativeNumberFilter")
	stages[0] = NegativeNumberFilter
	slog.Debug("Add stage to collection", "thread", "main", "index", 1, "Stage", "Dev3NumberFilter")
	stages[1] = Dev3NumberFilter

	// запускаем пайплайн
	slog.Debug("Start PipeLine", "thread", "main", "bufferSize", bufferSize)
	bufOut := ExecuteBufPipeline(in, done, bufferSize, stages...)

	// горутина печати буфера
	slog.Debug("Start PrintBufer", "thread", "main")
	go PrintBufer(bufOut, done)

	// стартуем сканер (сканер в рутине main)
	slog.Debug("Start ConsoleScaner", "thread", "main")
	StartScaner(in)

	// сворачиваем все рутины (закрываем сигнальный канал) они бы и так все закрылись.
	// но ... для примера done канала
	close(done)
}
