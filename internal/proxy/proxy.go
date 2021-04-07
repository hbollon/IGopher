package proxy

import (
	"fmt"
	"os/exec"
	"runtime"
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

	running                 bool
	stopProxyForwarderChan  chan bool
	errorProxyForwarderChan chan error
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

	p.stopProxyForwarderChan = make(chan bool)
	go func() {
		defer close(p.stopProxyForwarderChan)
		cmd := exec.Command(executable, options...)

		// Removed atm due to its incompatibility with OS other than Linux
		// cmd.SysProcAttr = &syscall.SysProcAttr{
		// 	Pdeathsig: syscall.SIGKILL,
		// }

		if err := cmd.Start(); err != nil {
			logrus.Errorf("Failed to launch local proxy-login-automator server: %v", err)
		}
		logrus.Debug("proxy-login-automator server successfully launched ! ")
		p.running = true

		p.errorProxyForwarderChan = make(chan error)
		defer close(p.errorProxyForwarderChan)
		go func() {
			p.errorProxyForwarderChan <- cmd.Wait()
		}()

		for {
			select {
			case <-p.stopProxyForwarderChan:
				cmd.Process.Kill()
				logrus.Debug("Successfully stopped proxy-login-automator server.")
				p.running = false
				return
			case err := <-p.errorProxyForwarderChan:
				logrus.Error(err)
				p.running = false
				return
			default:
				break
			}
		}
	}()
	time.Sleep(5 * time.Second)

	return nil
}

func (p *Proxy) RestartForwarderProxy() error {
	logrus.Debug("Restarting proxy-login-automator...")
	if p.running && p.stopProxyForwarderChan != nil {
		logrus.Debug("-> Stopping current proxy instance...")
		p.stopProxyForwarderChan <- true
	}
	if err := p.LaunchLocalForwarder(); err != nil {
		return err
	}
	logrus.Debug("Successfully restarted proxy-login-automator.")
	return nil
}

func (p *Proxy) StopForwarderProxy() {
	logrus.Debug("Stopping proxy-login-automator...")
	if p.running && p.stopProxyForwarderChan != nil {
		p.stopProxyForwarderChan <- true
	} else {
		logrus.Debug("proxy-login-automator isn't running.")
	}
}
