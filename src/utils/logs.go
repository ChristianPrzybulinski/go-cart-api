// Copyright Christian Przybulinski
// All Rights Reserved

//Package utils
package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//InitLog starts the logrus log setting up a log level, being info the default
func InitLog(logLevel string) {
	log.SetOutput(os.Stdout)

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	case "error", "err":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
