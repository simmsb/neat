package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/willfantom/neat/testbeds"
	"github.com/willfantom/neat/tests"
)

type Compose struct {
	Testbeds []*testbeds.Testbed `mapstructure:"testbeds"`
	Tests    []*tests.Test       `mapstructure:"tests"`
}

var uiLock = sync.Mutex{}

var (
	composeCmd = &cobra.Command{
		Use:   "compose",
		Short: "Run NEAT based on the NEAT Compose spec",
		Run: func(cmd *cobra.Command, args []string) {
			compose, err := parseComposeFile()
			if err != nil {
				logrus.WithField("extended", err.Error()).Fatalln("failed to parse compose file")
			}

			err = addTestbeds(compose.Testbeds)
			if err != nil {
				logrus.WithField("extended", err.Error()).Fatalln("failed to add testbed")
			}

			start := time.Now()

			err = createTestbeds(compose.Testbeds)
			if err != nil {
				logrus.WithField("extended", err.Error()).Fatalln("failed to create/start testbeds")
			}

			seenFailure := false

			for _, test := range compose.Tests {
				success, err := test.RunValid()
				if err != nil {
					logrus.WithField("extended", err.Error()).Fatalln("test failed")
				}
				if success {
					uiTestPassed(test.Name)
				} else {
					uiTestFailed(test.Name)
					seenFailure = true
				}
			}

			err = removeTestbeds(compose.Testbeds)
			if err != nil {
				logrus.WithField("extended", err.Error()).Fatalln("failed to stop/remove testbeds")
			}

			fmt.Printf("Total Time: %d\n", time.Since(start).Milliseconds())

			dumpStats(compose.Testbeds)

			if seenFailure {
				os.Exit(1)
			} else {
				os.Exit(0)
			}

		},
	}
)

func parseComposeFile() (*Compose, error) {
	composeViper := viper.New()
	composeViper.SetConfigName("neat-compose")
	composeViper.SetConfigType("yaml")
	composeViper.AddConfigPath("./.neat/.")
	composeViper.AddConfigPath(".")

	if err := composeViper.ReadInConfig(); err != nil {
		return nil, err
	}

	var compose Compose
	if err := composeViper.Unmarshal(&compose); err != nil {
		return nil, err
	}
	return &compose, nil
}

func addTestbeds(allTestbeds []*testbeds.Testbed) error {
	for _, testbed := range allTestbeds {
		fmt.Printf("Adding Testbed: %s\n", testbed.Name)
		if _, err := testbed.Add(); err != nil {
			return err
		}
	}
	return nil
}

func createTestbeds(allTestbeds []*testbeds.Testbed) error {
	wg := sync.WaitGroup{}
	for _, testbed := range allTestbeds {
		fmt.Printf("Creating Testbed: %s\n", testbed.Name)
		wg.Add(1)
		go func(testbed *testbeds.Testbed) {
			defer wg.Done()
			if err := testbeds.Create(testbed.Name); err != nil {
				fmt.Printf("Failed to Create Testbed: %s\n", testbed.Name)
				panic(err)
			}
			uiTestbedCreated(testbed.Name)
			if err := testbeds.Start(testbed.Name); err != nil {
				fmt.Printf("Failed to Start Testbed: %s\n", testbed.Name)
				panic(err)
			}
			uiTestbedStarted(testbed.Name)
		}(testbed)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
	return nil
}

func removeTestbeds(allTestbeds []*testbeds.Testbed) error {
	wg := sync.WaitGroup{}
	for _, testbed := range allTestbeds {
		wg.Add(1)
		go func(testbed *testbeds.Testbed) {
			defer wg.Done()
			fmt.Printf("Stopping Testbed: %s\n", testbed.Name)
			if err := testbeds.Stop(testbed.Name); err != nil {
				fmt.Printf("Failed to Stop Testbed: %s\n", testbed.Name)
				panic(err)
			}
			if err := testbeds.Remove(testbed.Name); err != nil {
				fmt.Printf("Failed to Remove Testbed: %s\n", testbed.Name)
				panic(err)
			}
			fmt.Printf("Stopped Testbed: %s\n", testbed.Name)
		}(testbed)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
	return nil
}

func uiTestbedCreated(tbName string) {
	uiLock.Lock()
	defer uiLock.Unlock()
	fmt.Printf("\n✔️\tTestbed Created: %s\n", strings.ToLower(tbName))
}

func uiTestbedStarted(tbName string) {
	uiLock.Lock()
	defer uiLock.Unlock()
	fmt.Printf("\n✅\tTestbed Started: %s\n", strings.ToLower(tbName))
}

func uiTestPassed(tName string) {
	uiLock.Lock()
	defer uiLock.Unlock()
	fmt.Printf("\n💯\tTest Passed: %s\n", strings.ToLower(tName))
}

func uiTestFailed(tName string) {
	uiLock.Lock()
	defer uiLock.Unlock()
	fmt.Printf("\n❌\tTest Failed: %s\n", strings.ToLower(tName))
}

func dumpStats(allTestbeds []*testbeds.Testbed) {
	for _, testbed := range allTestbeds {
		fmt.Printf("----------\nTestbed %s\n", testbed.Name)
		fmt.Printf("\tCreated %s\n", testbed.Metrics.CreatedAt.Format("15:04:05.0000"))
		fmt.Printf("\tRemoved %s\n", testbed.Metrics.RemovedAt.Format("15:04:05.0000"))
		fmt.Printf("\tStart Time %dms\n", testbed.Metrics.Runs[0].StartTime.Milliseconds())
		fmt.Printf("\tTotal Time %dms\n", testbed.Metrics.RemovedAt.Sub(testbed.Metrics.CreatedAt).Milliseconds())
		fmt.Printf("\tCPU Usage %f pct\n", testbed.Metrics.Runs[0].CPUUsage)
		fmt.Printf("\tMemory Usage %f pct\n", testbed.Metrics.Runs[0].PeakMemoryUsage)
		fmt.Printf("----------\n")
	}
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
