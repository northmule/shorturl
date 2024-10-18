package workers

import (
	"github.com/northmule/shorturl/internal/app/logger"
)

// workerNum количество воркеров.
const workerNum = 1

// Worker воркер удаления адресов пользователя.
type Worker struct {
	deleter  Deleter
	jobChan  chan job
	stopChan <-chan struct{}
}

// NewWorker конструктор.
func NewWorker(deleter Deleter, stop <-chan struct{}) *Worker {
	instance := Worker{
		deleter: deleter,
	}

	instance.jobChan = make(chan job, 1)
	instance.stopChan = stop
	for i := 0; i < workerNum; i++ {
		go instance.worker()
	}

	return &instance
}

// Deleter в фоне удаляет адреса пользователей.
type Deleter interface {
	SoftDeletedShortURL(userUUID string, shortURL ...string) error
}

type job struct {
	userUUID string
	url      []string
}

// Del удалить адреса у пользователя.
func (w *Worker) Del(userUUID string, input []string) {
	go w.producer(job{
		userUUID: userUUID,
		url:      input,
	})
}

func (w *Worker) producer(newJob job) {
	for {
		select {
		case <-w.stopChan:
			logger.LogSugar.Info("Поступил сигнал о закрытии продюсера")
			return
		case w.jobChan <- newJob:
			logger.LogSugar.Infof("В канал поступили ссылки для удаления: %v", newJob.url)
		}
	}
}

func (w *Worker) worker() {
	for {
		select {
		case <-w.stopChan:
			logger.LogSugar.Info("Поступил сигнал о закрытии воркера")
			return
		case jobs := <-w.jobChan:
			logger.LogSugar.Infof("Удаляю ссылки %v для пользователя %s", jobs.url, jobs.userUUID)
			err := w.deleter.SoftDeletedShortURL(jobs.userUUID, jobs.url...)
			if err != nil {
				logger.LogSugar.Infof(err.Error())
			}
		}
	}

}
