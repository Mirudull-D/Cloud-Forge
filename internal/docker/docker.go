package docker

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/client"
)

type Client struct {
	Cli *client.Client
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
		Cli: cli,
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
