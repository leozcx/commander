package main

import (
	"github.com/codegangsta/cli"
	"os"
	"path/filepath"
	"runtime"
)

func homepath(p string) string {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}
	return filepath.Join(home, p)
}

func getDiscovery(c *cli.Context) string {
	if len(c.Args()) == 1 {
		return c.Args()[0]
	}
	return os.Getenv("SWARM_DISCOVERY")
}

var (
	flStore = cli.StringFlag{
		Name:  "rootdir",
		Value: homepath(".swarm"),
		Usage: "",
	}
	flAddr = cli.StringFlag{
		Name:   "addr",
		Value:  "tcp://192.168.59.103:2376",
		Usage:  "ip to advertise",
		EnvVar: "SWARM_ADDR",
	}
	flHosts = cli.StringSliceFlag{
		Name:   "host, H",
		Value:  &cli.StringSlice{"tcp://127.0.0.1:2375"},
		Usage:  "ip/socket to listen on",
		EnvVar: "SWARM_HOST",
	}
	flTls = cli.BoolFlag{
		Name:  "tls",
		Usage: "use TLS; implied by --tlsverify=true",
	}
	flTlsCaCert = cli.StringFlag{
		Name:  "tlscacert",
		Usage: "trust only remotes providing a certificate signed by the CA given here",
	}
	flTlsCert = cli.StringFlag{
		Name:  "tlscert",
		Usage: "path to TLS certificate file",
	}
	flTlsKey = cli.StringFlag{
		Name:  "tlskey",
		Usage: "path to TLS key file",
	}
	flTlsVerify = cli.BoolFlag{
		Name:  "tlsverify",
		Usage: "use TLS and verify the remote",
	}
	flEtcd = cli.StringFlag{
		Name:  "etcd",
		Usage: "etcd URL",
	}
	flNode = cli.StringFlag{
		Name:  "node",
		Usage: "node ip address",
	}
)
