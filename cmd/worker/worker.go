package worker

import (
	"CloudHub/internal/deployments"
	"CloudHub/internal/docker"
	"CloudHub/internal/queue"
	"context"
	"log"

	"github.com/google/uuid"
)

type Worker struct {
	RedisStore      *queue.RedisStore
	DeploymentStore *deployments.Store
	DockerClient    *docker.Client
}

func NewWorker(redisStore *queue.RedisStore, deploymentStore *deployments.Store, client *docker.Client) *Worker {
	return &Worker{
		RedisStore:      redisStore,
		DeploymentStore: deploymentStore,
		DockerClient:    client,
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
		parsedUuid, err := uuid.Parse(deploymentID)
		if err != nil {
			log.Println("failed to parse deployment id:", err)
			continue
		}

		err = w.DeploymentStore.UpdateDeploymentStatusToBuilding(ctx, parsedUuid)
		if err != nil {
			log.Println("failed to update deployment status:", err)
			return
		}

		// TODO:
		// clone repo
		// build image
		// run container
	}
}
