package filter

import (
	"crypto/tls"
	"github.com/coreos/go-etcd/etcd"
)

type Context struct {
	Addr        string
	Debug       bool
	Version     string
	TLSConfig   *tls.Config
	EnableCores bool
	EtcdClient  *etcd.Client
}
