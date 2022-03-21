package ping_from

import (
	"errors"

	"github.com/antonmedv/expr"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/willfantom/neat/testbeds"
	"github.com/willfantom/neat/types"
)

func ValidateConfiguration(config map[string]interface{}) (bool, error) {

	var pingRequest types.PingRequest
	if err := mapstructure.Decode(config, &pingRequest); err != nil {
		return false, err
	}
	return true, nil
}

func Run(testbed *testbeds.Testbed, config map[string]interface{}) (map[string]interface{}, error) {
	var pingRequest types.PingRequest
	if err := mapstructure.Decode(config, &pingRequest); err != nil {
		return nil, err
	}

	result, err := testbed.DoPingFrom(pingRequest)
	if err != nil {
		return nil, err
	}
	return structs.Map(result), nil
}

func Evaluate(result map[string]interface{}, expression string) (bool, error) {
	program, err := expr.Compile(expression, expr.Env(result))
	if err != nil {
		return false, err
	}
	output, err := expr.Run(program, result)
	if err != nil {
		return false, err
	}

	if pass, ok := output.(bool); !ok {
		return false, errors.New("expression did not evaluate to bool")
	} else {
		return pass, nil
	}
}
