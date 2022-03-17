package docker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
)

const (
	neatLabelPrefix string = "neat"
)

type NeatContainer struct {
	ID          string            `mapstructure:"id"`
	Name        string            `mapstructure:"name"`
	Image       string            `mapstructure:"image"`
	Volumes     map[string]string `mapstructure:"volumes"`
	Labels      map[string]string `mapstructure:"labels"`
	Networks    []string          `mapstructure:"networks"`
	Environment map[string]string `mapstructure:"environment"`
	Privileged  bool              `mapstructure:"privileged"`
	TTY         bool              `mapstructure:"tty"`
	Command     []string          `mapstructure:"command"`

	StartStats []*ContainerStats `mapstructure:"start_stats"`
	StopStats  []*ContainerStats `mapstructure:"stop_stats"`
}

func (c *NeatContainer) ImageExists() bool {
	if c.Image == "" {
		return false
	}
	if _, _, err := docker.ImageInspectWithRaw(ctx, c.Image); err != nil {
		log.WithField("image", c.Image).Warnln("docker image could not be found")
		return false
	}
	return true
}

func (c *NeatContainer) Pull() error {
	if c.Image == "" {
		return fmt.Errorf("no image provided")
	}
	log.WithField("image", c.Image).Traceln("pulling docker container image")
	if _, err := docker.ImagePull(ctx, c.Image, types.ImagePullOptions{}); err != nil {
		log.Errorln(err.Error())
		return fmt.Errorf("image pull failed")
	}
	return nil
}

func (c *NeatContainer) Create() error {
	envStrings := make([]string, 0)
	for name, value := range c.Environment {
		envStrings = append(envStrings, fmt.Sprintf("%s=%s", name, value))
	}
	prefixedLabels := make(map[string]string)
	for label, value := range c.Labels {
		prefixedLabels[fmt.Sprintf("%s.%s", neatLabelPrefix, label)] = value
	}
	volumeBinds := make([]mount.Mount, 0)
	for host, container := range c.Volumes {
		volumeBinds = append(volumeBinds, mount.Mount{
			Type:   mount.TypeBind,
			Source: host,
			Target: container,
		})
	}
	containerConfig := &container.Config{
		Image:        c.Image,
		Cmd:          c.Command,
		Labels:       prefixedLabels,
		Env:          envStrings,
		Tty:          c.TTY,
		AttachStdin:  c.TTY,
		AttachStdout: c.TTY,
	}
	hostConfig := &container.HostConfig{
		Privileged:  c.Privileged,
		Mounts:      volumeBinds,
		NetworkMode: "bridge",
		CapAdd:      strslice.StrSlice{"sys_nice"},
		Resources:   container.Resources{
			// unstable, will add to configuration space soon
		},
	}
	respose, err := docker.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, c.Name)
	if err != nil {
		log.Errorln(err.Error())
		return fmt.Errorf("failed to create new container")
	}
	for _, net := range c.Networks {
		err := docker.NetworkConnect(ctx, net, respose.ID, nil)
		if err != nil {
			log.Errorln(err.Error())
			return fmt.Errorf("failed to add network to container")
		}
	}
	c.ID = respose.ID
	return nil
}

func (c *NeatContainer) Start() error {
	if c.ID == "" {
		return fmt.Errorf("container id needed to start a docker container")
	}
	if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		log.Errorln(err.Error())
		return fmt.Errorf("failed to start container")
	}
	if stats, err := c.Stats(); err != nil {
		log.WithField("id", c.ID).Warnln("no docker start stats could be collected")
	} else {
		c.StartStats = append(c.StartStats, stats)
	}
	return nil
}

func (c *NeatContainer) Stop() error {
	if c.ID == "" {
		return fmt.Errorf("container id needed to stop a docker container")
	}
	if stats, err := c.Stats(); err != nil {
		log.WithField("id", c.ID).Warnln("no docker stop stats could be collected")
	} else {
		c.StopStats = append(c.StartStats, stats)
	}
	if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
		log.Errorln(err.Error())
		return fmt.Errorf("failed to stop container")
	}
	return nil
}

func (c *NeatContainer) Remove() error {
	if c.ID == "" {
		return fmt.Errorf("container id needed to remove a docker container")
	}
	if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
		log.Errorln(err.Error())
		return fmt.Errorf("failed to remove container")
	}
	return nil
}

func (c *NeatContainer) GetIPS() ([]string, error) {
	container, err := c.inspect()
	if err != nil {
		return []string{}, err
	}

	ips := []string{}
	for _, net := range container.NetworkSettings.Networks {
		ips = append(ips, net.IPAddress)
	}

	return ips, nil
}

func (c *NeatContainer) Running() (bool, error) {
	container, err := c.inspect()
	if err != nil {
		return false, err
	}
	return container.State.Running, nil
}

func (c NeatContainer) Stats() (*ContainerStats, error) {
	stats, err := docker.ContainerStats(ctx, c.ID, false)
	if err != nil {
		log.WithField("id", c.ID).Errorln(err.Error())
		return nil, errors.New("failed to get container statistics")
	}
	buf := make([]byte, 128)
	for {
		tmpBuf := make([]byte, 128)
		_, err := stats.Body.Read(tmpBuf)
		buf = append(buf, tmpBuf...)
		if err == io.EOF {
			break
		}
	}
	buf = bytes.Trim(buf, "\x00")
	var containerStats ContainerStats
	if err := json.Unmarshal(buf, &containerStats); err != nil {
		log.WithField("id", c.ID).Errorln(err.Error())
		return nil, errors.New("container statistics format not valid")
	}
	return &containerStats, nil
}

func (c NeatContainer) inspect() (*types.ContainerJSON, error) {
	if c.ID == "" {
		return nil, fmt.Errorf("container id needed to inspect a docker container")
	}
	if container, err := docker.ContainerInspect(ctx, c.ID); err != nil {
		log.WithField("id", c.ID).Errorln(err.Error())
		return nil, fmt.Errorf("could not inspect docker container")
	} else {
		return &container, nil
	}
}
