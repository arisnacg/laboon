package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerContainer struct {
	ID    string
	Name  string
	Image string
	State string
}

type DockerWrapper struct {
	client *client.Client
}

func (dc *DockerWrapper) NewClient() {
	var err error
	dc.client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err.Error())
	}
}

func (dc *DockerWrapper) CloseClient() {
	dc.client.Close()
}

func (dc *DockerWrapper) GetContainers() []DockerContainer {
	containers, err := dc.client.ContainerList(
		context.Background(),
		container.ListOptions{All: true},
	)
	if err != nil {
		panic(err)
	}
	var dockerContainers []DockerContainer
	for _, container := range containers {
		dockerContainers = append(dockerContainers, DockerContainer{
			ID:    container.ID,
			Name:  container.Names[0][1:],
			Image: container.Image,
			State: container.State,
		})
	}
	return dockerContainers
}

func (dc *DockerWrapper) GetContainerState(id string) string {
	container, _ := dc.client.ContainerInspect(context.Background(), id)
	return container.State.Status
}

func (dc *DockerWrapper) PauseContainer(id string) {
	dc.client.ContainerPause(context.Background(), id)
}

func (dc *DockerWrapper) PauseContainers(ids []string) {
	for _, id := range ids {
		dc.PauseContainer(id)
	}
}

func (dc *DockerWrapper) UnpauseContainer(id string) {
	dc.client.ContainerUnpause(context.Background(), id)
}

func (dc *DockerWrapper) UnpauseContainers(ids []string) {
	for _, id := range ids {
		dc.UnpauseContainer(id)
	}
}

func (dc *DockerWrapper) StartContainer(id string) {
	dc.client.ContainerStart(context.Background(), id, container.StartOptions{})
}

func (dc *DockerWrapper) StartContainers(ids []string) {
	for _, id := range ids {
		dc.StartContainer(id)
	}
}

func (dc *DockerWrapper) StopContainer(id string) {
	dc.client.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (dc *DockerWrapper) StopContainers(ids []string) {
	for _, id := range ids {
		dc.StopContainer(id)
	}
}
