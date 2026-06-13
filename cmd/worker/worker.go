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

func NewWorker(
	redisStore *queue.RedisStore,
	deploymentStore *deployments.Store,
	client *docker.Client) *Worker {
	return &Worker{
		RedisStore:      redisStore,
		DeploymentStore: deploymentStore,
		DockerClient:    client,
	}
}

func (w *Worker) Run(ctx context.Context) {
	for {
		idInStr, err := w.RedisStore.PopDeployment(ctx)
		if err != nil {
			log.Println("failed to pop deployment:", err)
			continue
		}
		deploymentID, err := uuid.Parse(idInStr)
		if err != nil {
			log.Println("invalid deployment id:", idInStr)
			continue
		}

		if err = w.processDeployment(ctx, deploymentID); err != nil {

			err = w.DeploymentStore.UpdateDeploymentStatusToFailed(ctx, deploymentID, err.Error())
			if err != nil {
				log.Println("failed to update deployment status:", err)
				continue
			}

			log.Printf(
				"deployment %s failed: %v",
				deploymentID,
				err,
			)

			continue
		}
	}
}

func (w *Worker) processDeployment(
	ctx context.Context,
	deploymentId uuid.UUID,
) error {

	log.Printf(
		"processing deployment %s",
		deploymentId,
	)

	deployment, err := w.DeploymentStore.
		UpdateDeploymentStatusToBuilding(
			ctx,
			deploymentId,
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
		deploymentId,
	)

	imageName := "cloudhub-" + deploymentId.String()

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
	err = w.DeploymentStore.UpdateDeploymentStatusToRunning(
		ctx,
		deploymentId,
		imageName,
		containerID,
		port)
	if err != nil {
		return err
	}
	// update deployment status
	// save container id

	return nil
}
