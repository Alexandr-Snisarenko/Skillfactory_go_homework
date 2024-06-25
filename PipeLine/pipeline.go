package main

import "log/slog"

//import "stages"

// пайплай выводит результат в кольцевой буфер
func ExecuteBufPipeline(in chIn, done chDone, bufSize int, stages ...Stage) *RingBufferInt {
	var (
		stgIn  chIn
		stgOut chOut
	)

	bufOut := NewRingBufferInt(bufSize)

	// Запуксаем массив стейджей.
	// выходной канал текущего сдейджа = входному каналу следующего стейджа
	stgIn = in
	for _, f := range stages {
		stgOut = f(stgIn, done)
		stgIn = stgOut
	}

	// загружаем результат крайнего сдейджа в выходной буфер
	slog.Debug("Start потока записи в буфер", "thread", "WriteToBufer")
	go func() {
		for {
			select {
			case v, ok := <-stgOut:
				if !ok {
					slog.Debug("Выход по закрытию канала данных", "thread", "WriteToBufer")
					return
				}
				slog.Debug("Чтение данных из входного канала", "thread", "WriteToBufer", "value", v)
				bufOut.Push(v)
			case <-done:
				slog.Debug("Выход по сигнальному каналу done", "thread", "WriteToBufer")
				return
			}
		}
	}()

	return bufOut
}
