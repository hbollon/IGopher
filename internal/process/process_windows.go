//go:build windows

package process

import (
	"fmt"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// TerminateRunningInstance check if the pid stored in the pid file is running and, if yes, terminate it.
func TerminateRunningInstance() error {
	if res, psProcess := CheckIfAlreadyRunning(); res && psProcess != nil {
		process, err := os.FindProcess(psProcess.Pid())
		if err != nil {
			return err
		}

		dll, err := syscall.LoadDLL("kernel32.dll")
		if err != nil {
			logrus.Fatalf("LoadDLL: %v\n", err)
		}
		dllProc, err := dll.FindProc("GenerateConsoleCtrlEvent")
		if err != nil {
			logrus.Fatalf("FindProc: %v\n", err)
		}
		r, _, e := dllProc.Call(syscall.CTRL_BREAK_EVENT, uintptr(process.Pid))
		if r == 0 {
			logrus.Fatalf("GenerateConsoleCtrlEvent: %v\n", e)
		}

		return nil
	}

	return fmt.Errorf("Failed to recover running igopher process")
}
