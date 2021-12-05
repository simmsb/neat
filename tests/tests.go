package tests

import (
	"fmt"
	"time"

	"github.com/willfantom/neat/testbeds"
)

type Test struct {
	ID      string `mapstructure:"id" json:"id"`
	Name    string `mapstructure:"name" json:"name"`
	Variant string `mapstructure:"variant" json:"variant"`
	variant Variant
	Order   uint `mapstructure:"order" json:"order"`
	Repeats uint `mapstructure:"repeats" json:"repeats"`

	TestbedNames []string `mapstructure:"testbeds" json:"testbeds"`
	testbeds     []*testbeds.Testbed

	Expression       string `mapstructure:"expression" json:"expression"`
	Evaluate         string `mapstructure:"evaluate" json:"evaluate"`
	EvaluationScript string `mapstructure:"eval_script" json:"eval_script"`

	PreRunScript  string `mapstructure:"pre_run_script" json:"pre_run_script"`
	PreRun        string `mapstructure:"pre_run" json:"pre_run"`
	PostRunScript string `mapstructure:"post_run_script" json:"post_run_script"`
	PostRun       string `mapstructure:"post_run" json:"post_run"`

	VariantConfig map[string]interface{} `mapstructure:"config" json:"config"`

	Metrics map[string]Metrics `mapstructure:"metrics" json:"metrics"`
}

type Metrics struct {
	StartedAt     time.Time     `mapstructure:"started_at" json:"started_at"`
	ExecutionTime time.Duration `mapstructure:"execution_time" json:"execution_time"`
}

func (test *Test) Validate() (bool, error) {
	if !VariantExists(test.Variant) {
		return false, fmt.Errorf("test variant '%s' does not exist", test.Variant)
	}
	test.variant = variants[test.Variant]

	for _, testbedName := range test.TestbedNames {
		if testbed, err := testbeds.GetTestbed(testbedName); err != nil {
			return false, err
		} else {
			test.testbeds = append(test.testbeds, testbed)
		}
	}

	if validConfig, err := variants[test.Variant].ValidateConfiguration(test.VariantConfig); err != nil && !validConfig {
		return false, err
	} else if err == nil && !validConfig {
		return false, fmt.Errorf("test variant specific configuration is not valid")
	}
	if test.Evaluate == "" && test.Expression == "" && test.EvaluationScript == "" {
		return false, fmt.Errorf("no method to evaluate test provided")
	}
	// if test.Expression != "" {
	// 	if validConfig, err := variants[test.Variant].ValidateExpression(test.Expression); err != nil && !validConfig {
	// 		return false, err
	// 	} else if err == nil && !validConfig {
	// 		return false, fmt.Errorf("test expression is not valid for variant")
	// 	}
	// }
	return true, nil
}

//RunValid executes and evaluates the test on all given testbeds and for the given repeat value
//and also validates the tests configuration
func (test *Test) RunValid() (bool, error) {
	if valid, err := test.Validate(); !valid && err == nil {
		return false, fmt.Errorf("test configuration is not valid")
	} else if !valid && err != nil {
		return false, err
	}

	for _, testbed := range test.testbeds {
		result, err := test.variant.Run(testbed, test.VariantConfig)
		if err != nil {
			return false, fmt.Errorf("test %s failed to run on %s", test.Name, testbed.Name)
		}
		pass, err := test.variant.EvaluateExpression(result, test.Expression)
		if err != nil {
			return false, fmt.Errorf("test %s failed on %s", test.Name, testbed.Name)
		}
		if !pass {
			return false, nil
		}
	}
	return true, nil
}

//Run executes and evaluates the test on all given testbeds and for the given repeat value
func (test *Test) Run() (bool, error) {

	return true, nil
}
