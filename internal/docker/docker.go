package docker

import (
	"log"

	"github.com/docker/docker/client"
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
	return &Client{
		cli: cli,
	}, err
}
