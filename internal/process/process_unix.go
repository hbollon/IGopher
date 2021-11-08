//go:build !windows

package process

import (
	"fmt"
	"os"
	"syscall"
)

// TerminateRunningInstance check if the pid stored in the pid file is running and, if yes, terminate it.
func TerminateRunningInstance() error {
	if res, psProcess := CheckIfAlreadyRunning(); res && psProcess != nil {
		process, err := os.FindProcess(psProcess.Pid())
		if err != nil {
			return err
		}

		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Failed to recover running igopher process")
}
