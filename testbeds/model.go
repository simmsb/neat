package testbeds

import "time"

type Testbed struct {
	Name string `mapstructure:"name" json:"name"`

	VariantName string `mapstructure:"variant" json:"variant"`
	variant     Variant

	ResourceCap bool `mapstructure:"resource_cap" json:"resource_cap"`

	PreStartScript  string `mapstructure:"pre_start_script" json:"pre_start_script"`
	PreStart        string `mapstructure:"pre_start" json:"pre_start"`
	PostStartScript string `mapstructure:"post_start_script" json:"post_start_script"`
	PostStart       string `mapstructure:"post_start" json:"post_start"`

	PreStopScript  string `mapstructure:"pre_stop_script" json:"pre_stop_script"`
	PreStop        string `mapstructure:"pre_stop" json:"pre_stop"`
	PostStopScript string `mapstructure:"post_stop_script" json:"post_stop_script"`
	PostStop       string `mapstructure:"post_stop" json:"post_stop"`

	VariantConfig map[string]interface{} `mapstructure:"config" json:"config"`

	Metrics Metrics `mapstructure:"metrics" json:"metrics"`
}

type Metrics struct {
	CreatedAt    time.Time     `mapstructure:"created_at" json:"created_at"`
	CreationTime time.Duration `mapstructure:"creation_time" json:"creation_time"`
	RemovedAt    time.Time     `mapstructure:"removed_at" json:"removed_at"`
	RemoveTime   time.Duration `mapstructure:"remove_time" json:"remove_time"`
	Runs         []RunMetrics  `mapstructure:"runs" json:"runs"`
}

type RunMetrics struct {
	StartedAt       time.Time     `mapstructure:"started_at" json:"started_at"`
	StoppedAt       time.Time     `mapstructure:"stopped_at" json:"stopped_at"`
	StartTime       time.Duration `mapstructure:"start_time" json:"start_time"`
	StopTime        time.Duration `mapstructure:"stop_time" json:"stop_time"`
	CPUUsage        float64       `mapstructure:"cpu_usage" json:"cpu_usage"`
	PeakMemoryUsage float64       `mapstructure:"peak_memory_usage" json:"peak_memory_usage"`
}
