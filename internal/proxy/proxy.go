package proxy

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
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
		"-local_host 127.0.0.1",
		"-local_port 8880",
		fmt.Sprintf("-remote_host %s", p.RemoteIP),
		fmt.Sprintf("-remote_port %d", p.RemotePort),
		fmt.Sprintf("-usr %s", p.RemoteUsername),
		fmt.Sprintf("-pwd %s", p.RemotePassword),
	}

	cmd := exec.Command(executable, options...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}
	cmd.Start()

	return nil
}
