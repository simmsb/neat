package docker

import (
	"context"

	"github.com/docker/docker/client"
)

var (
	docker *client.Client
	ctx    context.Context = context.Background()
)

func setup() error {
	var err error
	docker, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return err
}
