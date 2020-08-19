package dockle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func Scan(ctx context.Context, config Config) (*ScanResult, error) {

	if config.ScanImageName == "" {
		err := fmt.Errorf("invalid parameter: %s", errors.New("the name of the image to be scanned is empty"))
		return nil, err
	}

	imageName, err := GetLatestImageName()
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest tag: %s", err)
	}

	client, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create DockerClient: %s", err)
	}

	if _, err := client.ImagePull(ctx, imageName, types.ImagePullOptions{}); err != nil {
		return nil, fmt.Errorf("failed to pull docker image: %s", err)
	}

	dockerConfig := container.Config{
		Image: imageName,
		Cmd:   []string{"-f", "json", config.ScanImageName},
		Tty:   true,
	}

	var hostConfig container.HostConfig
	if config.IsLocalImage {
		hostConfig.Mounts = append(hostConfig.Mounts,
			mount.Mount{ // -volume source:target
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		)
	}

	container, err := client.ContainerCreate(ctx, &dockerConfig, &hostConfig, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create docker container of image[%s]: %s", imageName, err)
	}

	if err := client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start docker container of image[%s]: %s", imageName, err)
	}

	if _, err := client.ContainerWait(ctx, container.ID); err != nil {
		return nil, fmt.Errorf("failed to wait for docker container to finish running: %s", err)
	}

	reader, err := client.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get container execution Log: %s", err)
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, reader); err != nil {
		return nil, fmt.Errorf("failed to read io.ReadCloser of container's execution has been written: %s", err)
	}

	var scanResult ScanResult
	err = json.Unmarshal(buf.Bytes(), &scanResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the results of a scan in the json format: %s", err)
	}
	return &scanResult, nil
}
