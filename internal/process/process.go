package process

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/mitchellh/go-ps"
	"github.com/sirupsen/logrus"
)

// CheckIfAlreadyRunning check for pid file (location is given by pidFilePath parameter) file existence.
// If exist, it'll get saved pid and check if the process is still running.
func CheckIfAlreadyRunning(pidFilePath string) bool {
	if _, err := os.Stat(pidFilePath); err == nil {
		var file *os.File
		file, err = os.Open(pidFilePath)
		if err != nil {
			logrus.Error("Failed to open existing pid file located at './data/pid.txt'.")
			logrus.Error(err)
			return true
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		if res := scanner.Scan(); !res {
			logrus.Warn("Pid file exist but without content, IGopher may be already running.")
			logrus.Info("Delete corrupt pid file and continue.")
			if err = os.Remove(pidFilePath); err != nil {
				logrus.Error("Failed to delete corrupt pid file!")
			}
			return false
		}
		pidStr := scanner.Text()

		pid, _ := strconv.Atoi(pidStr)
		var process ps.Process
		process, err = ps.FindProcess(pid)
		if process == nil && err == nil {
			logrus.Warnf("Failed to find process: %s\n. The pid must be outdated.", err)
			logrus.Info("Delete outdated pid file and continue.")
			if err = os.Remove(pidFilePath); err != nil {
				logrus.Error("Failed to delete corrupt pid file!")
			}
			return false
		}

		fmt.Println("running")
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		logrus.Fatalf(
			"Unknown issue during pid file checking: try to manually check if './data/pid.txt' exist and delete it. Detailed error: %v\n",
			err,
		)
	}

	return false
}

// DumpProcessPidToFile get program pid and save it to pidFilePath file
func DumpProcessPidToFile(pidFilePath string) {
	pid := strconv.Itoa(os.Getpid())
	file, err := os.Create(pidFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(pid)
	if err != nil {
		logrus.Fatalf("Failed to dump IGopher pid to file! Exit program. Detailed error: %v\n", err)
	}
}
