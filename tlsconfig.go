package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/denverdino/commander/state"
)

// Load the TLS certificates/keys and, if verify is true, the CA.
func loadTlsConfig(ca, cert, key string, verify bool) (*tls.Config, error) {
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load X509 key pair (%s, %s): %s. Key encrypted?",
			cert, key, err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{c},
		MinVersion:   tls.VersionTLS10,
	}

	if verify {
		certPool := x509.NewCertPool()
		file, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("Couldn't read CA certificate: %s", err)
		}
		certPool.AppendCertsFromPEM(file)
		config.RootCAs = certPool
	} else {
		// If --tlsverify is not supplied, disable CA validation.
		config.InsecureSkipVerify = true
	}

	return config, nil
}

func parseTlsConfig(c *cli.Context) (*tls.Config, error) {
	var (
		tlsConfig *tls.Config = nil
		err       error
	)

	// If either --tls or --tlsverify are specified, load the certificates.
	if c.Bool("tls") || c.Bool("tlsverify") {
		if !c.IsSet("tlscert") || !c.IsSet("tlskey") {
			log.Fatal("--tlscert and --tlskey must be provided when using --tls")
		}
		if c.Bool("tlsverify") && !c.IsSet("tlscacert") {
			log.Fatal("--tlscacert must be provided when using --tlsverify")
		}
		tlsConfig, err = loadTlsConfig(
			c.String("tlscacert"),
			c.String("tlscert"),
			c.String("tlskey"),
			c.Bool("tlsverify"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Otherwise, if neither --tls nor --tlsverify are specified, abort if
		// the other flags are passed as they will be ignored.
		if c.IsSet("tlscert") || c.IsSet("tlskey") || c.IsSet("tlscacert") {
			log.Fatal("--tlscert, --tlskey and --tlscacert require the use of either --tls or --tlsverify")
		}
	}

	store := state.NewStore(path.Join(c.String("rootdir"), "state"))
	if err := store.Initialize(); err != nil {
		log.Fatal(err)
	}
	return tlsConfig, err
}
