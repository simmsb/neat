package main

import (
	"github.com/sirupsen/logrus"
	"github.com/willfantom/neat/cmd"
	_ "github.com/willfantom/neat/testbeds/mtv"
)

func main() {

	logrus.SetLevel(logrus.TraceLevel)

	if err := cmd.Execute(); err != nil {
		logrus.Fatalln(err.Error())
	}

}
