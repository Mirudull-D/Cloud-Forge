package docker

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Client struct {
	cli *client.Client
}

func NewDockerClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatal(err)
	}
	info, err := cli.Info(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Docker connected Successfully ...!!!!\n \t\t\t\t\tServer Version:", info.ServerVersion)

	return &Client{
		cli: cli,
	}, err
}

func (c *Client) BuildImage(
	workspace string,
	imageName string,
) error {

	cmd := exec.Command(
		"docker",
		"build",
		"-t",
		imageName,
		workspace,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
func (c *Client) RunContainer(
	ctx context.Context,
	imageName string,
	hostPort string,
) (string, error) {

	containerPort, err := nat.NewPort(
		"tcp",
		"8080",
	)
	if err != nil {
		return "", err
	}

	resp, err := c.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			ExposedPorts: nat.PortSet{
				containerPort: struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				containerPort: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: hostPort,
					},
				},
			},
		},
		nil,
		nil,
		"",
	)

	if err != nil {
		return "", err
	}

	err = c.cli.ContainerStart(
		ctx,
		resp.ID,
		container.StartOptions{},
	)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	if err := c.cli.ContainerStop(
		ctx,
		containerID,
		container.StopOptions{
			Timeout: new(10),
		},
	); err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveContainer(ctx context.Context, containerID string) error {
	err := c.cli.ContainerRemove(
		ctx,
		containerID,
		container.RemoveOptions{
			Force: true,
		})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveImage(ctx context.Context, imageName string) error {
	_, err := c.cli.ImageRemove(
		ctx,
		imageName,
		image.RemoveOptions{
			Force: true,
		})
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) RestartContainer(ctx context.Context, containerID string) error {
	err := c.cli.ContainerRestart(ctx, containerID, container.StopOptions{
		Timeout: new(10),
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) InspectContainer(
	ctx context.Context, containerID string) (container.InspectResponse, error) {
	inspect, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return container.InspectResponse{}, err
	}
	return inspect, nil
}
