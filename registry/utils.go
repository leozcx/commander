package registry

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
)

// GetValues queries etcd for keys prefixed by prefix.
func GetValues(client etcd.Client, keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, key := range keys {
		resp, err := client.Get(key, true, true)
		if err != nil {
			return vars, err
		}
		err = nodeWalk(resp.Node, vars)
		if err != nil {
			return vars, err
		}
	}
	return vars, nil
}

// nodeWalk recursively descends nodes, updating vars.
func nodeWalk(node *etcd.Node, vars map[string]string) error {
	if node != nil {
		key := node.Key
		if !node.Dir {
			vars[key] = node.Value
		} else {
			for _, node := range node.Nodes {
				nodeWalk(node, vars)
			}
		}
	}
	return nil
}

func watchPrefix(client *etcd.Client, prefix string, waitIndex uint64, stopChan chan bool) (*etcd.Response, uint64, error) {
	if waitIndex == 0 {
		resp, err := client.Get(prefix, false, true)
		if err != nil {
			return resp, 0, err
		}
		return resp, resp.EtcdIndex, nil
	}
	resp, err := client.Watch(prefix, waitIndex+1, true, nil, stopChan)
	if err != nil {
		return resp, waitIndex, err
	}
	return resp, resp.Node.ModifiedIndex, err
}
