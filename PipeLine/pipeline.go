package main

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
	go func() {
		for {
			select {
			case v, ok := <-stgOut:
				if !ok {
					return
				}
				bufOut.Push(v)
			case <-done:
				return
			}
		}
	}()

	return bufOut
}
