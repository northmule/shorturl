package workers

import (
	"fmt"
	"github.com/northmule/shorturl/internal/app/logger"
	"sync"
)

const workerNum = 20

type Worker struct {
	deleter Deleter
}

func NewWorker(deleter Deleter) *Worker {
	instance := Worker{
		deleter: deleter,
	}

	return &instance
}

type Deleter interface {
	SoftDeletedShortURL(shortURL string) error
}

func (w *Worker) Del(input []string) {

	// сигнальный канал для завершения горутин
	doneCh := make(chan struct{})
	// закрываем его при завершении программы
	defer close(doneCh)

	// канал с данными
	inputCh := w.fillInputChan(doneCh, input)

	// получаем слайс каналов из 10 рабочих add
	channels := w.fanOut(doneCh, inputCh)

	// а теперь объединяем десять каналов в один
	resultCh := w.fanIn(doneCh, channels...)

	// выводим результаты расчетов из канала
	for res := range resultCh {
		fmt.Println(res)
	}
}

func (w *Worker) fillInputChan(doneCh chan struct{}, input []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

func (w *Worker) deleteShortURL(doneCh chan struct{}, inputCh chan string) chan string {
	addRes := make(chan string)

	go func() {
		defer close(addRes)

		for data := range inputCh {
			err := w.deleter.SoftDeletedShortURL(data)
			if err != nil {
				logger.LogSugar.Infof(err.Error())
				continue
			}

			select {
			case <-doneCh:
				return
			case addRes <- data:
			}
		}
	}()
	return addRes
}

// fanOut принимает канал данных
func (w *Worker) fanOut(doneCh chan struct{}, inputCh chan string) []chan string {
	// каналы, в которые отправляются результаты
	channels := make([]chan string, workerNum)

	for i := 0; i < workerNum; i++ {
		// получаем канал из горутины add
		addResultCh := w.deleteShortURL(doneCh, inputCh)
		// отправляем его в слайс каналов
		channels[i] = addResultCh
	}

	// возвращаем слайс каналов
	return channels
}

// fanIn объединяет несколько каналов resultChs в один.
func (w *Worker) fanIn(doneCh chan struct{}, resultChs ...chan string) chan string {
	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
	finalCh := make(chan string)

	// понадобится для ожидания всех горутин
	var wg sync.WaitGroup

	// перебираем все входящие каналы
	for _, ch := range resultChs {
		// в горутину передавать переменную цикла нельзя, поэтому делаем так
		chClosure := ch

		// инкрементируем счётчик горутин, которые нужно подождать
		wg.Add(1)

		go func() {
			// откладываем сообщение о том, что горутина завершилась
			defer wg.Done()

			// получаем данные из канала
			for data := range chClosure {
				select {
				// выходим из горутины, если канал закрылся
				case <-doneCh:
					return
				// если не закрылся, отправляем данные в конечный выходной канал
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		// ждём завершения всех горутин
		wg.Wait()
		// когда все горутины завершились, закрываем результирующий канал
		close(finalCh)
	}()

	// возвращаем результирующий канал
	return finalCh
}
