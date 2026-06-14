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
	"strconv"

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
				err.Error(),
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

	//imageName := "cloudhub-" + deploymentId.String()

	imageName := "cloudhub-7f79e8bd-22e8-4c96-81f8-1118b39611ac"
	//err = w.DockerClient.BuildImage(
	//	tempDir,
	//	imageName,
	//)
	//if err != nil {
	//	return err
	//}
	//log.Printf("building successful for image %s", imageName)
	port, err := w.DeploymentStore.GetNextAvailablePort(ctx)
	if err != nil {
		return err
	}

	containerID, err := w.DockerClient.RunContainer(
		ctx,
		imageName,
		strconv.Itoa(port),
	)

	if err != nil {
		return err
	}

	log.Println("container started:", containerID)

	err = w.DeploymentStore.UpdateDeploymentStatusToRunning(
		ctx,
		deploymentId,
		imageName,
		containerID,
		port)
	if err != nil {
		return err
	}
	return nil
}
