package proxy

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	localServerHost  string
	remoteServerHost string
	remoteServerAuth string
)

// ProxyConfig store all remote proxy configuration
type ProxyConfig struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Enabled  bool   `yaml:"activated"`
}

func PrintResponse(r *http.Response) error {
	logrus.Infof("Response: %+v\n", r)
	return nil
}

// LaunchForwardingProxy launch forward server used to inject proxy authentication header
// into outgoing requests
func LaunchForwardingProxy(localPort uint16, remoteProxy ProxyConfig) error {
	localServerHost = fmt.Sprintf("localhost:%d", localPort)
	remoteServerHost = fmt.Sprintf(
		"http://%s:%d",
		remoteProxy.IP,
		remoteProxy.Port,
	)
	remoteServerAuth = fmt.Sprintf(
		"%s:%s",
		remoteProxy.Username,
		remoteProxy.Password,
	)

	remote, err := url.Parse(remoteServerHost)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	d := func(req *http.Request) {
		logrus.Infof("Pre-Edited request: %+v\n", req)
		// Inject proxy authentication headers to outgoing request into new Header
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(remoteServerAuth))
		req.Header.Set("Proxy-Authorization", basicAuth)
		// Change host to the remote proxy
		req.URL = remote
		logrus.Infof("Edited Request: %+v\n", req)
		logrus.Infof("Scheme: %s, Host: %s, Port: %s\n", req.URL.Scheme, req.URL.Host, req.URL.Port())
	}
	proxy.Director = d
	proxy.ModifyResponse = PrintResponse
	http.ListenAndServe(localServerHost, proxy)

	return nil
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	// Inject proxy authentication headers to outgoing request into new Header
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(remoteServerAuth))
	r.Header.Add("Proxy-Authorization", basicAuth)

	// Prepare new request for remote proxy
	bodyRemote, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*Preparation of the new request
	Part I'm not sure about */

	// Create new request
	hostURL := fmt.Sprintf("%s://%s", "http", remoteServerHost)
	proxyReq, err := http.NewRequest(r.Method, hostURL, bytes.NewReader(bodyRemote))
	if err != nil {
		http.Error(w, "Could not create new request", 500)
		return
	}

	// Copy header
	proxyReq.Header = r.Header
	logrus.Info(proxyReq)

	/* end of request preparation */

	// Forward request to remote proxy server
	httpClient := http.Client{}
	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		logrus.Info(err)
		http.Error(w, "Could not reach origin server", 500)
		return
	}
	defer resp.Body.Close()

	logrus.Infof("Response: %v", resp)

	// Transfer header from origin server -> client
	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)

	// Transfer response from origin server -> client
	if resp.ContentLength > 0 {
		io.CopyN(w, resp.Body, resp.ContentLength)
	} else if resp.Close {
		// Copy until EOF or some other error occurs
		for {
			if _, err := io.Copy(w, resp.Body); err != nil {
				break
			}
		}
	}
}
