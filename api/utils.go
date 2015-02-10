package api

import (
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/denverdino/commander/api/filter"
	"github.com/samalba/dockerclient"
	"io"
	"net"
	"net/http"
	"strings"
)

func newClientAndScheme(tlsConfig *tls.Config) (*http.Client, string) {
	if tlsConfig != nil {
		return &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}, "https"
	}
	return &http.Client{}, "http"
	//return &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}, "https"
}

// from https://github.com/golang/go/blob/master/src/net/http/httputil/reverseproxy.go#L82
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func getContainerFromVars(c *filter.Context, vars map[string]string) (*dockerclient.ContainerInfo, error) {
	client, err := newDockerClient(c)
	if err != nil {
		return nil, err
	}
	if name, ok := vars["name"]; ok {
		if container, _ := client.InspectContainer(name); container != nil {
			return container, nil
		}
		return nil, fmt.Errorf("No such container: %s", name)

	}
	//TODO: Optimize with the etcd access
	if ID, ok := vars["execid"]; ok {
		containers, _ := client.ListContainers(true, false, "")
		if containers != nil {
			for _, container := range containers {
				containerInfo, _ := client.InspectContainer(container.Id)
				if containerInfo != nil {
					for _, execID := range containerInfo.ExecIDs {
						if ID == execID {
							return containerInfo, nil
						}
					}
				}
			}
		}
		return nil, fmt.Errorf("Exec %s not found", ID)
	}
	return nil, errors.New("Not found")
}

func proxy(tlsConfig *tls.Config, addr string, w http.ResponseWriter, r *http.Request) (int, error) {
	// Use a new client for each request
	client, scheme := newClientAndScheme(tlsConfig)
	// RequestURI may not be sent to client
	r.RequestURI = ""

	r.URL.Scheme = scheme
	r.URL.Host = addr

	log.WithFields(log.Fields{"method": r.Method, "url": r.URL}).Debug("Proxy request")
	resp, err := client.Do(r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(NewWriteFlusher(w), resp.Body)

	return resp.StatusCode, nil
}

func hijack(tlsConfig *tls.Config, addr string, w http.ResponseWriter, r *http.Request) error {
	if parts := strings.SplitN(addr, "://", 2); len(parts) == 2 {
		addr = parts[1]
	}

	log.WithField("addr", addr).Debug("Proxy hijack request")

	var (
		d   net.Conn
		err error
	)

	if tlsConfig != nil {
		d, err = tls.Dial("tcp", addr, tlsConfig)
	} else {
		d, err = net.Dial("tcp", addr)
	}
	if err != nil {
		return err
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		return err
	}
	nc, _, err := hj.Hijack()
	if err != nil {
		return err
	}
	defer nc.Close()
	defer d.Close()

	err = r.Write(d)
	if err != nil {
		return err
	}

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		if conn, ok := dst.(interface {
			CloseWrite() error
		}); ok {
			conn.CloseWrite()
		}
		errc <- err
	}
	go cp(d, nc)
	go cp(nc, d)
	<-errc
	<-errc

	return nil
}

func newDockerClient(c *filter.Context) (dockerclient.Client, error) {
	docker, err := dockerclient.NewDockerClient(c.Addr, c.TLSConfig)
	return docker, err
}
