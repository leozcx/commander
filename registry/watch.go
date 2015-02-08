package registry

import (
	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
)

func watch(client etcd.Client) {
	watchChan := make(chan *etcd.Response)
	go client.Watch(SERVICE_PREFIX, 0, false, watchChan, nil)
	log.Println("Waiting for an update...")
	r := <-watchChan
	log.Printf("Got updated creds: %s: %s\n", r.Node.Key, r.Node.Value)
}
