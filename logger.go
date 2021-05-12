package igopher

import (
	"bufio"
	"os"
	"runtime"

	logRuntime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/rifflock/lfshook"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const logFilePath = "./logs/logs.log"

func InitLogger() {
	setLoggerOutput()
	level, err := log.ParseLevel(*Flags.LogLevelFlag)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.InfoLevel)
		log.Warnf("Invalid log level '%s', use default one.", *Flags.LogLevelFlag)
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

	// Parse logs to string array
	var logs []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}

	// Build json array string with logs from newer to older
	out := `[`
	for i := len(logs) - 1; i >= 0; i-- {
		out += logs[i]
		if i == 0 {
			break
		}
		out += `,`
	}
	out += `]`

	return out, nil
}
