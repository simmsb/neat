package testbeds

import (
	"github.com/sirupsen/logrus"
	"github.com/willfantom/neat/tools"
	"github.com/willfantom/neat/types"
)

type Variant struct {
	Name        string
	Description string

	Tools []tools.Tool

	ValidateConfiguration func(config map[string]interface{}) (bool, error)
	Create                func(testbed *Testbed) error
	Start                 func(testbed *Testbed) error
	Stop                  func(testbed *Testbed) error
	Remove                func(testbed *Testbed) error

	HookArguments func(path string, testbed *Testbed) []string

	DoPing func(testbed *Testbed, request types.PingRequest) (*types.PingResponse, error)
	DoPingFrom func(testbed *Testbed, request types.PingRequest) (*types.PingResponse, error)
}

func VariantExists(name string) bool {
	logrus.WithField("variant", name).Traceln("checking if testbed vairant exists")
	_, ok := Variants[name]
	return ok
}

var Variants = map[string]Variant{}
