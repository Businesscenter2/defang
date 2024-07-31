package docker

import (
	"context"
	"errors"

	"github.com/DefangLabs/defang/src/pkg/types"
	"github.com/docker/docker/client"
)

type ContainerID = types.TaskID

type Docker struct {
	*client.Client

	image    string
	memory   uint64
	platform string
}

func New() *Docker {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &Docker{
		Client: cli,
	}
}

var _ types.Driver = (*Docker)(nil)

func (Docker) PutConfig(ctx context.Context, name, value string, isSensitive bool) error {
	return errors.New("docker does not support secrets")
}

func (Docker) GetConfig(ctx context.Context, name []string) (types.ConfigData, error) {
	return nil, errors.New("docker does not support secrets")
}

func (Docker) ListConfigs(ctx context.Context) ([]string, error) {
	return nil, errors.New("docker does not support secrets")
}

func (Docker) CreateUploadURL(ctx context.Context, name string) (string, error) {
	return "", errors.New("docker does not support uploads")
}
