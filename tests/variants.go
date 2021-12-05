package tests

import (
	"fmt"

	"github.com/willfantom/neat/testbeds"
	"github.com/willfantom/neat/tests/ping"
)

type Variant struct {
	Name        string
	Description string

	ValidateConfiguration func(config map[string]interface{}) (bool, error)
	ValidateExpression    func(expression string) (bool, error)
	Run                   func(testbed *testbeds.Testbed, config map[string]interface{}) (map[string]interface{}, error)
	EvaluateExpression    func(result map[string]interface{}, expression string) (bool, error)
	EvaluateScript        func(result map[string]interface{}, script string) (bool, error)
}

type VaraintNotExistError struct {
	givenVariant     string
	suggestedVariant string
}

func (e *VaraintNotExistError) Error() string {
	base := fmt.Sprintf("variant '%s' does not exist", e.givenVariant)
	if e.suggestedVariant != "" {
		return fmt.Sprintf("%s: did you mean '%s'", base, e.suggestedVariant)
	}
	return base
}

func VariantExists(name string) bool {
	_, ok := variants[name]
	return ok
}

var variants = map[string]Variant{
	"ping": {
		Name:                  "Ping",
		Description:           "Check connectivity between 2 network nodes using ICMP echo packets",
		ValidateConfiguration: ping.ValidateConfiguration,
		Run:                   ping.Run,
		EvaluateExpression:    ping.Evaluate,
	},
}
