package mtv

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/neat/testbeds"
	"github.com/willfantom/neat/tools"
	"github.com/willfantom/neat/tools/docker"
)

type Config struct {
	Image       string            `mapstructure:"image"`
	Libvirt     bool              `mapstructure:"libvirt"`
	Files       string            `mapstructure:"files"`
	ExtraMounts map[string]string `mapstructure:"extra_mounts"`
	Command     string            `mapstructure:"command"`
	Networks    []string          `mapstructure:"networks"`
}

const (
	dockerImage string = "ghcr.io/ng-cdi/mtv:test"
)

var variant = testbeds.Variant{
	Name:        "MTV",
	Description: "A mininet fork designed for VNF testing",

	Tools: []tools.Tool{
		docker.Tool,
	},

	ValidateConfiguration: validateConfiguration,
	Create:                create,
	Start:                 start,
	Stop:                  stop,
	Remove:                remove,

	HookArguments: getArguments,

	DoPing: doPing,
}

func parseConfig(config map[string]interface{}) (*Config, error) {
	var parsedConfig Config
	if err := mapstructure.Decode(config, &parsedConfig); err != nil {
		return nil, err
	}
	return &parsedConfig, nil
}

func init() {
	for _, tool := range variant.Tools {
		if !tool.Check() {
			logrus.WithFields(logrus.Fields{
				"variant": strings.ToLower(variant.Name),
				"tool":    strings.ToLower(tool.Name),
			}).Warnln("testbed variant unavailable due to failed due to tool dependency check failure")
			return
		}
	}
	testbeds.Variants["mtv"] = variant
}
