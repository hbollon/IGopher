// +build windows

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

		dll, e := syscall.LoadDLL("kernel32.dll")
		if e != nil {
			logrus.Fatalf("LoadDLL: %v\n", e)
		}
		process, e := dll.FindProc("GenerateConsoleCtrlEvent")
		if e != nil {
			log.Fatalf("FindProc: %v\n", e)
		}
		r, _, e := process.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
		if r == 0 {
			log.Fatalf("GenerateConsoleCtrlEvent: %v\n", e)
		}

		return nil
	}

	return fmt.Errorf("Failed to recover running igopher process")
}
