package main

import "log/slog"

type (
	chIn   = <-chan int
	chOut  = chIn
	chDone = <-chan interface{}
)

// стейджи для текущего пайплайна
type Stage func(in chIn, done chDone) chOut

// фильтр негативных чисел
func NegativeNumberFilter(in chIn, done chDone) chOut {
	ch := make(chan int)
	slog.Debug("Start stage NegativeNumberFilter", "thread", "NegativeNumberFilter")
	go func() {
		defer close(ch)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					slog.Debug("Выход по закрытию канала данных", "thread", "NegativeNumberFilter")
					return
				}
				slog.Debug("Чтение данных из входного канала", "thread", "NegativeNumberFilter", "value", v)
				if v > 0 {
					slog.Debug("Запись фильтровнного значения в выходной канал", "thread", "NegativeNumberFilter", "value", v)
					ch <- v
				}
			case <-done:
				slog.Debug("Выход по сигнальному каналу done", "thread", "NegativeNumberFilter")
				return
			}
		}
	}()
	return ch
}

// фильтр чисел, которые не делятся на 3 (+0)
func Dev3NumberFilter(in chIn, done chDone) chOut {
	ch := make(chan int)
	slog.Debug("Start stage Dev3NumberFilter", "thread", "Dev3NumberFilter")
	go func() {
		defer close(ch)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					slog.Debug("Выход по закрытию канала данных", "thread", "Dev3NumberFilter")
					return
				}
				slog.Debug("Чтение данных из входного канала", "thread", "Dev3NumberFilter", "value", v)
				if v != 0 && (v%3) == 0 {
					slog.Debug("Запись фильтровнного значения в выходной канал", "thread", "Dev3NumberFilter", "value", v)
					ch <- v
				}
			case <-done:
				slog.Debug("Выход по сигнальному каналу done", "thread", "Dev3NumberFilter")
				return
			}
		}
	}()
	return ch
}
