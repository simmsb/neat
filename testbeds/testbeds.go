package testbeds

import (
	"fmt"
	"strings"
	"time"

	"github.com/willfantom/neat/tools/script"
	"github.com/willfantom/neat/types"
)

var testbeds = make(map[string]*Testbed)

func GetTestbed(searchTerm string) (*Testbed, error) {
	//TODO: maybe do some fuzzy search for names (not ids though)
	searchTerm = strings.ToLower(searchTerm)
	for id, testbed := range testbeds {
		if strings.ToLower(id) == searchTerm || strings.ToLower(testbed.Name) == searchTerm {
			return testbed, nil
		}
	}
	return nil, fmt.Errorf("testbed with name/id '%s' not found", searchTerm)
}

func (testbed *Testbed) Validate() (bool, error) {
	if !VariantExists(testbed.VariantName) {
		return false, fmt.Errorf("testbed variant '%s' does not exist", testbed.VariantName)
	} else {
		testbed.variant = Variants[testbed.VariantName]
	}

	return true, nil
}

func (testbed *Testbed) Add() (string, error) {
	if valid, err := testbed.Validate(); !valid {
		return "", err
	}
	id := generateID()
	testbeds[id] = testbed
	return id, nil
}

func Create(id string) error {
	testbed, err := GetTestbed(id)
	if err != nil {
		return err
	}
	start := time.Now()
	if err := testbed.variant.Create(testbed); err != nil {
		return err
	}
	testbed.Metrics.CreationTime = time.Since(start)
	testbed.Metrics.CreatedAt = time.Now()
	return nil
}

func Start(id string) error {
	testbed, err := GetTestbed(id)
	if err != nil {
		return err
	}
	if testbed.PreStartScript != "" {
		if err := script.Run(testbed.variant.HookArguments(testbed.PreStartScript, testbed)...); err != nil {
			return err
		}
	}
	if err := testbed.variant.Start(testbed); err != nil {
		return err
	}
	if testbed.PostStartScript != "" {
		if err := script.Run(testbed.variant.HookArguments(testbed.PostStartScript, testbed)...); err != nil {
			return err
		}
	}
	return nil
}

func Stop(id string) error {
	testbed, err := GetTestbed(id)
	if err != nil {
		return err
	}
	if testbed.PreStopScript != "" {
		if err := script.Run(testbed.variant.HookArguments(testbed.PreStopScript, testbed)...); err != nil {
			return err
		}
	}
	if err := testbed.variant.Stop(testbed); err != nil {
		return err
	}
	if testbed.PostStopScript != "" {
		if err := script.Run(testbed.variant.HookArguments(testbed.PostStopScript, testbed)...); err != nil {
			return err
		}
	}
	return nil
}

func Remove(id string) error {
	testbed, err := GetTestbed(id)
	if err != nil {
		return err
	}
	if err := testbed.variant.Remove(testbed); err != nil {
		return err
	}
	return nil
}

func (testbed *Testbed) DoPing(request types.PingRequest) (*types.PingResponse, error) {
	return testbed.variant.DoPing(testbed, request)
}
func (testbed *Testbed) DoPingFrom(request types.PingRequest) (*types.PingResponse, error) {
	return testbed.variant.DoPingFrom(testbed, request)
}
