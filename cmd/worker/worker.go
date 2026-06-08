package worker

import (
	"CloudHub/internal/queue"
	"context"
	"log"
)

type Worker struct {
	RedisStore *queue.RedisStore
}

func NewWorker(redisStore *queue.RedisStore) *Worker {
	return &Worker{
		RedisStore: redisStore,
	}
}

func (w *Worker) Run(ctx context.Context) {
	for {
		deploymentID, err := w.RedisStore.PopDeployment(ctx)
		if err != nil {
			log.Println("failed to pop deployment:", err)
			continue
		}

		log.Printf("processing deployment %s\n", deploymentID)

		// TODO:
		// update deployment status
		// clone repo
		// build image
		// run container
	}
}
