package proxy

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

const (
	CR = 0xd
	LF = 0xa
)

var (
	ConnectRequest = []byte{67, 79, 78, 78, 69, 67, 84, 32}
	GetRequest     = []byte{47, 45, 54, 20}
)

// Copy data between two connections. Return EOF on connection close.
func Pipe(a, b net.Conn) error {
	done := make(chan error, 1)

	cpWithInjection := func(b io.Writer, a io.Reader) {
		var err error
		var written int64
		var buf []byte

		size := 32 * 1024
		if l, ok := a.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)

		for {
			nr, er := a.Read(buf)
			if nr > 0 {
				if buf, err = injectProxyAuthorization(buf); err != nil {
					logrus.Error(err)
				}
				nw, ew := b.Write(buf)
				if nw < 0 {
					nw = 0
					if ew == nil {
						ew = fmt.Errorf("Write failed")
					}
				}
				written += int64(nw)
				if ew != nil {
					err = ew
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
		}

		logrus.Debugf("Copied %d bytes", written)
		done <- err
	}

	cp := func(a, b net.Conn) {
		n, err := io.Copy(a, b)
		logrus.Debugf("copied %d bytes from %s to %s", n, a.RemoteAddr(), b.RemoteAddr())
		done <- err
	}

	go cpWithInjection(b, a)
	go cp(a, b)

	err1 := <-done
	if err1 != nil {
		return err1
	}

	err2 := <-done
	if err2 != nil {
		return err2
	}

	return nil
}

func injectProxyAuthorization(buf []byte) ([]byte, error) {
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(remoteServerAuth))
	var connData string
	if buf != nil {
		logrus.Debugf("content init: %s", string(buf))
		if isValidRequest(buf) {
			buf = bytes.Trim(buf, "\x00")
			buf = buf[:len(buf)-2]

			connData = string(buf)
			connData = connData + "Proxy-Authorization: " + basicAuth

			buf = []byte(connData)
			buf = append(buf, 0x0d, 0x0a, 0x0d, 0x0a)
			logrus.Debugf("content edited: %s", string(buf))
		}
	} else {
		return nil, fmt.Errorf("Buffer is empty")
	}

	return buf, nil
}

func isValidRequest(buf []byte) bool {
	if bytes.Contains(buf, ConnectRequest) || bytes.Contains(buf, GetRequest) {
		return true
	}
	return false
}
