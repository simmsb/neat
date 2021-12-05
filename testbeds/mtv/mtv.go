package mtv

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/willfantom/neat/testbeds"
	"github.com/willfantom/neat/testbeds/mtv/mnapi"
	"github.com/willfantom/neat/tools/docker"
	"github.com/willfantom/neat/types"
)

var containers = make(map[string]*docker.NeatContainer)

func validateConfiguration(config map[string]interface{}) (bool, error) {
	parsedConfig, err := parseConfig(config)
	if err != nil {
		return false, err
	}
	if parsedConfig.Files == "" {
		return false, fmt.Errorf("files must be provided to an mtv testbed")
	}
	return true, nil
}

func create(testbed *testbeds.Testbed) error {
	parsedConfig, err := parseConfig(testbed.VariantConfig)
	if err != nil {
		return err
	}
	if parsedConfig.Image == "" {
		parsedConfig.Image = dockerImage
	}
	start := time.Now()
	currentDir, _ := os.Getwd()
	container := docker.NeatContainer{
		Name:  testbed.Name,
		Image: parsedConfig.Image,
		Volumes: map[string]string{
			filepath.Join(currentDir, parsedConfig.Files): "/mnt",
		},
		Labels: map[string]string{
			"name":    testbed.Name,
			"variant": testbed.VariantName,
		},

		Environment: map[string]string{
			"SCRIPT":    "/mnt/topology.py",
			"ASSET_DIR": currentDir,
		},
		Privileged: true,
		TTY:        true,
		Command:    []string{},
	}
	if parsedConfig.Libvirt {
		container.Volumes["/var/run/libvirt/libvirt-sock"] = "/var/run/libvirt/libvirt-sock"
		// container.Volumes["/var/run/docker.sock"] = "/var/run/docker.sock"
	}

	if err := container.Create(); err != nil {
		return err
	}
	containers[testbed.Name] = &container
	testbed.Metrics.CreatedAt = time.Now()
	testbed.Metrics.CreationTime = time.Since(start)
	return nil
}

func start(testbed *testbeds.Testbed) error {
	if container, ok := containers[testbed.Name]; !ok {
		return fmt.Errorf("mtv testbed has no container")
	} else {
		start := time.Now()
		err := container.Start()
		if err != nil {
			return err
		}
		ip, err := container.GetIP()
		if err != nil {
			return err
		}
		for {
			//TODO: don't hardcode so much of this
			// if time.Since(start) > (300 * time.Second) {
			// 	return fmt.Errorf("failed to start mtv api")
			// }
			client, err := mnapi.NewClient("http://"+ip+":8080", nil)
			if err != nil {
				return fmt.Errorf("failed to create mtv client")
			}
			if _, err = client.GetNodes(); err == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		testbed.Metrics.Runs = append(testbed.Metrics.Runs, testbeds.RunMetrics{
			StartedAt: time.Now(),
			StartTime: time.Since(start),
		})
		return nil
	}
}

func stop(testbed *testbeds.Testbed) error {
	if container, ok := containers[testbed.Name]; !ok {
		return fmt.Errorf("mtv testbed has no container")
	} else {
		start := time.Now()
		err := container.Stop()
		if err != nil {
			return err
		}
		testbed.Metrics.Runs[len(testbed.Metrics.Runs)-1].StoppedAt = time.Now()
		testbed.Metrics.Runs[len(testbed.Metrics.Runs)-1].StopTime = time.Since(start)
		containerCPUUsage := container.StopStats[len(container.StopStats)-1].CPUStats.Usage.Total - container.StartStats[len(container.StartStats)-1].CPUStats.Usage.Total
		containerMemoryUsage := float64(container.StopStats[len(container.StopStats)-1].MemoryStats.MaxUsage) / float64(container.StopStats[len(container.StopStats)-1].MemoryStats.Limit)
		systemCPUUsage := container.StopStats[len(container.StopStats)-1].CPUStats.SystemUsage - container.StartStats[len(container.StartStats)-1].CPUStats.SystemUsage
		testbed.Metrics.Runs[len(testbed.Metrics.Runs)-1].CPUUsage = ((float64(containerCPUUsage) / float64(systemCPUUsage)) * float64(len(container.StopStats[len(container.StopStats)-1].CPUStats.Usage.PerCPU)) * 100)
		testbed.Metrics.Runs[len(testbed.Metrics.Runs)-1].PeakMemoryUsage = containerMemoryUsage * 100
		return nil
	}
}

func remove(testbed *testbeds.Testbed) error {
	if container, ok := containers[testbed.Name]; !ok {
		return fmt.Errorf("mtv testbed has no container")
	} else {
		start := time.Now()
		err := container.Remove()
		if err != nil {
			return err
		}
		testbed.Metrics.RemovedAt = time.Now()
		testbed.Metrics.RemoveTime = time.Since(start)

		return nil
	}
}

func doPing(testbed *testbeds.Testbed, request types.PingRequest) (*types.PingResponse, error) {
	if container, ok := containers[testbed.Name]; !ok {
		return nil, fmt.Errorf("mtv testbed has no container")
	} else {
		ip, err := container.GetIP()
		if err != nil {
			return nil, fmt.Errorf("failed to get container ip")
		}
		client, err := mnapi.NewClient("http://"+ip+":8080", nil)
		if err != nil {
			return nil, err
		}
		pingData, err := client.PingSet([]string{request.Sender, request.Target})
		if err != nil {
			return nil, err
		}
		for _, pingResponse := range pingData {
			if pingResponse.Sender == request.Sender {
				return &types.PingResponse{
					Sent:     uint(pingResponse.Sent),
					Received: uint(pingResponse.Received),
				}, nil
			}
		}
	}
	return nil, errors.New("failed to get ping response")
}

func getArguments(path string, testbed *testbeds.Testbed) []string {
	if container, ok := containers[testbed.Name]; !ok {
		return []string{path}
	} else {
		return []string{path, container.ID}
	}
}
