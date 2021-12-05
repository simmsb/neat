package docker

import (
	"github.com/sirupsen/logrus"
	"github.com/willfantom/neat/tools"
)

var (
	Tool = tools.Tool{
		Name:        "Docker",
		Description: "Interact with a local docker engine via a unix socket",

		Check: check,
	}
	setupOK bool          = false
	log     *logrus.Entry = logrus.WithField("tool", "docker")
)

func check() bool {
	if setupOK {
		return true
	}
	log.Traceln("checking tool")
	if err := setup(); err != nil {
		log.Errorln(err.Error())
		return false
	} else {
		setupOK = true
	}
	return true
}
