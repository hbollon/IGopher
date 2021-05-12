// +build !dev

package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hbollon/igopher"
	tui "github.com/hbollon/igopher/internal/tui"
	"github.com/mitchellh/go-ps"
	"github.com/sirupsen/logrus"
)

const pidFilePath = "./data/pid.txt"

func init() {
	igopher.Flags.BackgroundFlag = flag.Bool("background-task", false,
		"Run IGopher as background task with actual configuration (configure it normally and after re-run IGopher with this flag)")
}

// Change the current working directory by executable location one
func changeWorkingDir() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	} else {
		err = os.Chdir(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Check for ./data/pid.txt file existence. If exist, it'll get saved pid and check if the process is still running.
func checkIfAlreadyRunning() bool {
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

// Get program pid and save it to ./data/pid.txt file
func dumpProcessPidToFile() {
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

func main() {
	flag.Parse()
	changeWorkingDir()
	igopher.InitLogger()

	// Initialize environment
	igopher.CheckEnvironment()

	alreadyRunning := checkIfAlreadyRunning()
	if *igopher.Flags.BackgroundFlag {
		if alreadyRunning {
			logrus.Error("IGopher is already running! Kill it or close it through TUI interface and retry.")
			return
		}
		dumpProcessPidToFile()
		logrus.Debug("Successfully dump pid to tmp file!")
		igopher.LaunchBotTui()
	} else {
		if alreadyRunning {
			logrus.Error("IGopher is already running! Kill it or close it through TUI interface and retry.")
			return
		}

		// Clear terminal session
		igopher.ClearTerminal()

		// Launch TUI
		execBot := tui.InitTui()

		// Launch bot if option selected
		if execBot {
			igopher.LaunchBotTui()
		}
	}
}
