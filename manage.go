package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/coreos/go-etcd/etcd"
	"github.com/denverdino/commander/api"
	"github.com/denverdino/commander/context"
	"github.com/samalba/dockerclient"
)

func manage(c *cli.Context) {
	tlsConfig, err := parseTlsConfig(c)

	if err != nil {
		log.Fatal(err)
	}

	// see https://github.com/codegangsta/cli/issues/160
	addr := c.String("addr")
	if addr == "" {
		log.Fatalf("addr required to access a cluster. See '%s manage --help'.", c.App.Name)
	}

	dockerClient, err := dockerclient.NewDockerClient(addr, tlsConfig)

	if err != nil {
		log.Fatal(err)
	}

	// see https://github.	com/codegangsta/cli/issues/160
	hosts := c.StringSlice("host")
	if c.IsSet("host") || c.IsSet("H") {
		hosts = hosts[1:]
	}

	// see https://github.com/codegangsta/cli/issues/160
	etcdURL := c.String("etcd")
	if etcdURL == "" {
		log.Fatalf("etcd required to access a etcd server. See '%s watch --help'.", c.App.Name)
	}

	etcdClient := etcd.NewClient([]string{etcdURL})

	context := &context.Context{
		Addr:         addr,
		Version:      c.App.Version,
		EtcdClient:   etcdClient,
		TLSConfig:    tlsConfig,
		DockerClient: dockerClient,
	}

	log.Fatal(api.ListenAndServe(context, hosts, c.Bool("cors")))
}
