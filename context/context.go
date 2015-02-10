package context

import (
	"crypto/tls"
	"github.com/coreos/go-etcd/etcd"
	"github.com/samalba/dockerclient"
)

type Context struct {
	Addr         string
	Debug        bool
	Version      string
	TLSConfig    *tls.Config
	EnableCores  bool
	EtcdClient   *etcd.Client
	DockerClient *dockerclient.DockerClient
}
