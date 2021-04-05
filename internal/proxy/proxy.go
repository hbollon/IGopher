package proxy

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// Proxy store all remote proxy configuration
type Proxy struct {
	RemoteIP       string `yaml:"ip"`
	RemotePort     int    `yaml:"port"`
	RemoteUsername string `yaml:"username"`
	RemotePassword string `yaml:"password"`

	WithAuth bool `yaml:"auth"`
	Enabled  bool `yaml:"activated"`
}

func (p *Proxy) LaunchLocalForwarder() error {
	var executable string
	if runtime.GOOS == "windows" {
		executable = "./lib/proxy-login-automator.exe"
	} else {
		executable = "./lib/proxy-login-automator"
	}

	options := []string{
		"-local_host",
		"127.0.0.1",
		"-local_port",
		"8880",
		"-remote_host",
		p.RemoteIP,
		"-remote_port",
		fmt.Sprintf("%d", p.RemotePort),
		"-usr",
		p.RemoteUsername,
		"-pwd",
		p.RemotePassword,
	}

	stopProxyForwarderChan := make(chan bool)
	go func() {
		defer close(stopProxyForwarderChan)
		cmd := exec.Command(executable, options...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Pdeathsig: syscall.SIGKILL,
		}
		if err := cmd.Start(); err != nil {
			logrus.Errorf("Failed to launch local proxy-login-automator server: %v", err)
		}
		logrus.Debug("proxy-login-automator server successfully launched ! ")

		errorProxyForwarderChan := make(chan error)
		defer close(errorProxyForwarderChan)
		go func() {
			errorProxyForwarderChan <- cmd.Wait()
		}()

		for {
			select {
			case <-stopProxyForwarderChan:
				cmd.Process.Kill()
				logrus.Debug("Successfully stopped proxy-login-automator server.")
				return
			case err := <-errorProxyForwarderChan:
				logrus.Error(err)
				return
			default:
				break
			}
		}
	}()
	time.Sleep(5 * time.Second)

	return nil
}
