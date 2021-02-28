package igopher

import (
	"bufio"
	"flag"
	"os"
	"runtime"

	logRuntime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/rifflock/lfshook"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const logFilePath = "./logs/logs.log"

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
			logrus.InfoLevel:  logFilePath,
			logrus.WarnLevel:  logFilePath,
			logrus.ErrorLevel: logFilePath,
			logrus.FatalLevel: logFilePath,
		},
		&logrus.JSONFormatter{},
	))
}

// Read and parse log file to json array string
func parseLogsToString() (string, error) {
	// Open log file
	file, err := os.Open(logFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	out := `[`
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	for {
		out += scanner.Text()
		if !scanner.Scan() {
			break
		}
		out += `,`
	}
	out += `]`

	return out, nil
}
