package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	logrus.SetOutput(os.Stdout)
}
