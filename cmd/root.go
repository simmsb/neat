package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logLevel string

	rootCmd = &cobra.Command{
		Use:   "neat",
		Short: "Automatically test virtual network functions",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if logLevel != "" {
				if logrus_level, err := logrus.ParseLevel(logLevel); err != nil {
					logrus.Fatalln("invalid log level given")
				} else {
					logrus.SetLevel(logrus_level)
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "v", "info", "logging level")
}
