package main

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
	go func() {
		defer close(ch)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				if v > 0 {
					ch <- v
				}
			case <-done:
				return
			}
		}
	}()
	return ch
}

// фильтр чисел, которые не делятся на 3 (+0)
func Dev3NumberFilter(in chIn, done chDone) chOut {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				if v != 0 && (v%3) == 0 {
					ch <- v
				}
			case <-done:
				return
			}
		}
	}()
	return ch
}
