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

	// warning: brain rot ahead
	volumes := make(map[string]string, len(parsedConfig.ExtraMounts))
	for k, v := range parsedConfig.ExtraMounts {
		volumes[k] = v
	}
	volumes[filepath.Join(currentDir, parsedConfig.Files)] = "/mnt"

	container := docker.NeatContainer{
		Name:  testbed.Name,
		Image: parsedConfig.Image,
		Volumes: volumes,
		Networks: parsedConfig.Networks,
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
		ips, err := container.GetIPS()
		if err != nil {
			return err
		}
		fmt.Printf("\tIP of test container: %s\n", ips)
	outer:
		for {
			//TODO: don't hardcode so much of this
			// if time.Since(start) > (300 * time.Second) {
			// 	return fmt.Errorf("failed to start mtv api")
			// }

			for _, ip := range ips {
				client, err := mnapi.NewClient("http://"+ip+":8080", nil)
				if err != nil {
					return fmt.Errorf("failed to create mtv client")
				}
				if nodes, err := client.GetNodes(); err == nil {
					fmt.Println("nodes: ", nodes)
					testbed.VariantConfig["MNApiIP"] = ip
					break outer
				}
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
	if _, ok := containers[testbed.Name]; !ok {
		return nil, fmt.Errorf("mtv testbed has no container")
	} else {
		ip := testbed.VariantConfig["MNApiIP"].(string)
		client, err := mnapi.NewClient("http://"+ip+":8080", nil)
		if err != nil {
			return nil, err
		}
		fmt.Println("trying to ping %s", ip)
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
