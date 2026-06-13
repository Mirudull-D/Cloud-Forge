package worker

import (
	"CloudHub/internal/deployments"
	"CloudHub/internal/docker"
	"CloudHub/internal/github"
	"CloudHub/internal/queue"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

		if err := w.processDeployment(ctx, deploymentID); err != nil {
			log.Printf(
				"deployment %s failed: %v",
				deploymentID,
				err,
			)
		}
	}
}

func (w *Worker) processDeployment(
	ctx context.Context,
	deploymentID string,
) error {

	log.Printf(
		"processing deployment %s",
		deploymentID,
	)

	parsedUUID, err := uuid.Parse(deploymentID)
	if err != nil {
		return err
	}

	deployment, err := w.DeploymentStore.
		UpdateDeploymentStatusToBuilding(
			ctx,
			parsedUUID,
		)

	if err != nil {
		return err
	}

	tempDir, err := github.CreateRepoInTemp(
		deployment.GitUrl,
	)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Println(
				"failed to cleanup temp dir:",
				err,
			)
		}
	}()

	dockerfile := filepath.Join(
		tempDir,
		"Dockerfile",
	)

	if _, err := os.Stat(dockerfile); err != nil {
		return fmt.Errorf("dockerfile not found in temp dir: %s", tempDir)
	}

	log.Printf(
		"dockerfile found for deployment %s",
		deploymentID,
	)

	imageName := "cloudhub-" + deploymentID

	err = w.DockerClient.BuildImage(
		tempDir,
		imageName,
	)
	if err != nil {
		return err
	}
	log.Printf("building sucessfull for image %s", imageName)

	containerID, err := w.DockerClient.RunContainer(
		ctx,
		imageName,
		"8081",
	)

	if err != nil {
		return err
	}

	log.Println("container started:", containerID)
	// TODO:
	// update deployment status
	// save container id

	return nil
}
