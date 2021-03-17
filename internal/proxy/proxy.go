package proxy

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	localServerHost  string
	remoteServerHost string
	remoteServerAuth string
)

// Proxy store all remote proxy configuration
type Proxy struct {
	RemoteIP       string `yaml:"ip"`
	RemotePort     int    `yaml:"port"`
	RemoteUsername string `yaml:"username"`
	RemotePassword string `yaml:"password"`
	Enabled        bool   `yaml:"activated"`

	listener net.Listener
}

func (p *Proxy) Run() error {
	var err error
	if p.listener, err = net.Listen("tcp", ":8880"); err != nil {
		return err
	}
	logrus.Infof("Port forwarding server up and listening on %s", localServerHost)

	wg := &sync.WaitGroup{}
	for {
		if conn, err := p.listener.Accept(); err == nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.handle(conn)
			}()
		} else {
			wg.Wait()
			return nil
		}
	}
}

func (p *Proxy) Close() error {
	return p.listener.Close()
}

func (p *Proxy) handle(upConn net.Conn) {
	defer upConn.Close()
	logrus.Infof("accepted: %s", upConn.RemoteAddr())
	downConn, err := net.Dial("tcp", remoteServerHost)
	if err != nil {
		logrus.Errorf("unable to connect to %s: %s", remoteServerHost, err)
		return
	}
	defer downConn.Close()
	if err := Pipe(upConn, downConn); err != nil {
		logrus.Errorf("pipe failed: %s", err)
	} else {
		logrus.Infof("disconnected: %s", upConn.RemoteAddr())
	}
}

// LaunchForwardingProxy launch forward server used to inject proxy authentication header
// into outgoing requests
func InitForwardingProxy(localPort uint16, remoteProxy Proxy) error {
	localServerHost = fmt.Sprintf("localhost:%d", localPort)
	remoteServerHost = fmt.Sprintf(
		"%s:%d",
		remoteProxy.RemoteIP,
		remoteProxy.RemotePort,
	)
	remoteServerAuth = fmt.Sprintf(
		"%s:%s",
		remoteProxy.RemoteUsername,
		remoteProxy.RemotePassword,
	)

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		if err := remoteProxy.Close(); err != nil {
			logrus.Fatal(err.Error())
		}
	}()

	if err := remoteProxy.Run(); err != nil {
		logrus.Fatal(err.Error())
	}

	return nil
}
