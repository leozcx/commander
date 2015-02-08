package main

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Commander of the Docker Cluster"
	app.Version = "0.1.0"
	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "DEBUG",
		},

		cli.StringFlag{
			Name:  "log-level, l",
			Value: "info",
			Usage: fmt.Sprintf("Log level (options: debug, info, warn, error, fatal, panic)"),
		},
	}

	// logs
	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.SetLevel(level)

		// If a log level wasn't specified and we are running in debug mode,
		// enforce log-level=debug.
		if !c.IsSet("log-level") && !c.IsSet("l") && c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      "manage",
			ShortName: "m",
			Usage:     "manage a docker cluster",
			Flags: []cli.Flag{
				flStore,
				flAddr, flHosts,
				flTls, flTlsCaCert, flTlsCert, flTlsKey, flTlsVerify,
				flEtcd,
				flNode, //TODO: Remove
			},
			Action: manage,
		},
		{
			Name:      "watch",
			ShortName: "w",
			Usage:     "watch a docker daemon",
			Flags: []cli.Flag{
				flStore,
				flAddr, flHosts,
				flTls, flTlsCaCert, flTlsCert, flTlsKey, flTlsVerify,
				flNode, flEtcd,
			},
			Action: watch,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
