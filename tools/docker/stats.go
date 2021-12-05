package docker

import "time"

type ContainerStats struct {
	ReadTime    time.Time         `json:"read"`
	PreReadTime time.Time         `json:"preread"`
	CPUStats    ContainerCPUStats `json:"cpu_stats"`
	PreCPUStats ContainerCPUStats `json:"precpu_stats"`
	MemoryStats ContainerMemory   `json:"memory_stats"`
}

type ContainerCPUStats struct {
	Usage       ContainerCPUUsageStats `json:"cpu_usage"`
	SystemUsage uint64                 `json:"system_cpu_usage"`
	OnlineCPUs  uint64                 `json:"online_cpus"`
}

type ContainerCPUUsageStats struct {
	Usermode   uint64   `json:"usage_in_usermode"`
	Total      uint64   `json:"total_usage"`
	KernelMode uint64   `json:"usage_in_kernelmode"`
	PerCPU     []uint64 `json:"percpu_usage"`
}

type ContainerMemory struct {
	Usage    uint64 `json:"usage"`
	MaxUsage uint64 `json:"max_usage"`
	Limit    uint64 `json:"limit"`
	Stats    ContainerMemoryStats
}

type ContainerMemoryStats struct {
	Cache uint64 `json:"cache"`
}
