package registry

import (
	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
)

func WatchServiceChange(client *etcd.Client) {

	var resp *etcd.Response
	var waitIndex uint64
	var err error

	for {
		resp, waitIndex, err = watchPrefix(client, SERVICE_PREFIX, waitIndex, nil)
		if resp != nil {
			log.Printf("Got updated: %s %s: %s\n", resp.Action, resp.Node.Key, resp.Node.Value)
		}
		if err != nil {
			log.Error(err)
		}
	}
	/*
		watchChan := make(chan *etcd.Response)
		for {
			go client.Watch(SERVICE_PREFIX, 0, false, watchChan, nil)
			log.Println("Waiting for an update...")
			r := <-watchChan
			log.Printf("Got updated creds: %s: %s\n", r.Node.Key, r.Node.Value)
		}
	*/

}
