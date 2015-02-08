package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/denverdino/commander/api"
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

	log.Fatal(api.ListenAndServe(addr, hosts, c.App.Version, c.Bool("cors"), etcdURL, tlsConfig))
}
