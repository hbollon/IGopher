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
	RemoteIP       string `json:"ip" yaml:"ip" validate:"required,contains=."`
	RemotePort     int    `json:"port,string" yaml:"port" validate:"required,numeric,min=1,max=65535"`
	RemoteUsername string `json:"username" yaml:"username"`
	RemotePassword string `json:"password" yaml:"password"`

	WithAuth bool `json:"auth,string" yaml:"auth"`
	Enabled  bool `json:"proxyActivation,string" yaml:"activated"`

	running                 bool
	stopProxyForwarderChan  chan bool
	errorProxyForwarderChan chan error
}

// LaunchLocalForwarder launch an instance of proxy-login-automator (https://github.com/hbollon/proxy-login-automator) which starts
// a local forwarder proxy server in order to be able to automatically inject the "Proxy-Authorization" header
// to all outgoing Selenium requests and forward them to the remote proxy configured by the user.
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
				<-p.errorProxyForwarderChan // ignore cmd.Wait() output
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
			time.Sleep(10 * time.Millisecond)
		}
	}()
	time.Sleep(5 * time.Second)

	return nil
}

// RestartForwarderProxy check for running instance of proxy-login-automator, stop it if exist and finally start a new one
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

// StopForwarderProxy stop current running instance of proxy-login-automator
func (p *Proxy) StopForwarderProxy() {
	logrus.Debug("Stopping proxy-login-automator...")
	if p.running && p.stopProxyForwarderChan != nil {
		p.stopProxyForwarderChan <- true
	} else {
		logrus.Debug("proxy-login-automator isn't running.")
	}
}
