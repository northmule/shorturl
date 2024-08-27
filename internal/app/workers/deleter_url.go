package workers

import (
	"github.com/northmule/shorturl/internal/app/logger"
)

const workerNum = 50

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
	SoftDeletedShortURL(userUUID string, shortURL ...string) error
}

type job struct {
	uuid string
	url  string
}

func (w *Worker) Del(userUUID string, input []string) {
	jobCh := make(chan job, len(input))

	for i := 0; i < workerNum; i++ {
		go w.worker(jobCh)
	}

	go w.producer(userUUID, input, jobCh)
}

func (w *Worker) producer(uuid string, urls []string, jobCh chan<- job) {
	for _, url := range urls {
		jobCh <- job{uuid: uuid, url: url}
	}
}

func (w *Worker) worker(jobCh <-chan job) {
	for item := range jobCh {
		err := w.deleter.SoftDeletedShortURL(item.uuid, item.url)
		if err != nil {
			logger.LogSugar.Infof(err.Error())
		}
	}
}
