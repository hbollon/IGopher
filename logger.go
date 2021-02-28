package igopher

import (
	"flag"
	"os"
	"runtime"

	logRuntime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/rifflock/lfshook"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	setLoggerOutput()

	flag.Parse()
	level, err := log.ParseLevel(*flags.LogLevelFlag)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.InfoLevel)
		log.Warnf("Invalid log level '%s', use default one.", *flags.LogLevelFlag)
	}
}

func setLoggerOutput() {
	// Initialize logs folder
	if _, err := os.Stat("./logs/"); os.IsNotExist(err) {
		os.Mkdir("./logs/", os.ModePerm)
	}

	// Add formatter to logrus in order to display line and function with messages on Stdout
	formatter := logRuntime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)

	if runtime.GOOS == "windows" {
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	} else {
		log.SetOutput(os.Stdout)
	}

	// Add hook to logrus to also redirect logs to files with custom formatter
	log.AddHook(lfshook.NewHook(
		lfshook.PathMap{
			logrus.InfoLevel:  "./logs/igopher.log",
			logrus.WarnLevel:  "./logs/igopher.log",
			logrus.ErrorLevel: "./logs/igopher.log",
			logrus.FatalLevel: "./logs/igopher.log",
		},
		&log.TextFormatter{
			FullTimestamp: false,
			DisableColors: true,
		},
	))
}
