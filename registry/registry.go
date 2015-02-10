package registry

import (
	"encoding/json"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
	"github.com/samalba/dockerclient"
	"strconv"
	"strings"
)

func getContainerResourceURI(id string) string {
	return CONTAINER_PREFIX + id
}

func getContainerResourceInfoURI(id string) string {
	return CONTAINER_PREFIX + id + INFO_POSTFIX
}

func DeleteContainer(etcdClient *etcd.Client, node string, id string) error {
	data, _ := getContainerResource(etcdClient, id)

	if data != nil {
		key := getContainerResourceURI(id)
		_, err := etcdClient.Delete(key, true)

		hostContainerKey := getHostContainerResourceURI(node, id)
		_, err = etcdClient.Delete(hostContainerKey, false)

		serviceId := getEnvValue(data.Config.Env, SERVICE_ID)

		if len(serviceId) > 0 {
			serviceContainerKey := getServiceContainerResourceURI(serviceId, id)
			_, err = etcdClient.Delete(serviceContainerKey, false)
		}
		return err
	} else {
		//TODO NOT FOUND
		return nil
	}

}

func getEnvValue(environments []string, name string) string {

	prefix := name + "="

	for _, env := range environments {
		log.Info("Env: ", env)
		if strings.HasPrefix(env, prefix) {
			return env[len(prefix):]
		}
	}
	return ""
}

func getContainerResource(etcdClient *etcd.Client, id string) (*ContainerResource, error) {

	key := getContainerResourceInfoURI(id)
	resp, err := etcdClient.Get(key, false, false)

	if err != nil {
		return nil, err
	}

	var data ContainerResource

	if resp != nil {
		if err := json.Unmarshal([]byte(resp.Node.Value), &data); err != nil {
			log.Warn(err)
			return nil, err
		}
		return &data, nil
	}
	//NOT FOUND
	return nil, nil
}

func SetContainer(client *dockerclient.DockerClient, etcdClient *etcd.Client, node string, id string) error {

	info, err := client.InspectContainer(id)

	if err != nil {
		if err == dockerclient.ErrNotFound {
			log.Info("Delete container info: ", id)
			//Delete
			DeleteContainer(etcdClient, node, id)
		}
		return err
	}

	key := getContainerResourceInfoURI(id)

	data, _ := getContainerResource(etcdClient, id)

	if data != nil { // Update only
		log.Info("Update container info: ", id)
	} else {
		log.Info("Create container info: ", id)
	}

	serviceId := getEnvValue(info.Config.Env, SERVICE_ID)

	log.Info("ServiceId = ", serviceId)

	data = &ContainerResource{
		"container",
		node,
		info.Id,
		info.Name,
		serviceId,
		info.Image,
		info.Config,
		info.Created,
		info.State,
		info.NetworkSettings}

	if jsonData, err := json.Marshal(data); err == nil {
		jsonStr := string(jsonData)
		log.Infoln(jsonStr)
		_, err = etcdClient.Set(key, jsonStr, 0)
	}

	hostContainerKey := getHostContainerResourceURI(node, id)
	_, err = etcdClient.Set(hostContainerKey, strconv.FormatBool(info.State.Running), 0)

	if len(serviceId) > 0 {
		serviceContainerKey := getServiceContainerResourceURI(serviceId, id)
		log.Info("serviceContainerKey = ", serviceContainerKey)
		_, err = etcdClient.Set(serviceContainerKey, strconv.FormatBool(info.State.Running), 0)
	}

	return err
}

func SetContainers(client *dockerclient.DockerClient, etcdClient *etcd.Client, node string) error {

	containers, err := client.ListContainers(true, false, "")
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		err = SetContainer(client, etcdClient, node, container.Id)
		if err != nil {
			log.Fatal(err)
		}
	}
	return err
}

func getHostResourceURI(node string) string {
	return HOST_PREFIX + node + INFO_POSTFIX
}

func getHostContainerResourceURI(node string, id string) string {
	return HOST_PREFIX + node + "/containers/" + id
}

func SetHost(client *dockerclient.DockerClient, etcdClient *etcd.Client, node string) error {

	info, err := client.Info()
	info.Containers = -1
	info.Images = -1

	key := getHostResourceURI(node)
	var data = HostResource{"node", *info}

	if jsonData, err := json.Marshal(data); err == nil {
		jsonStr := string(jsonData)
		_, err = etcdClient.Set(key, jsonStr, 0)
	}

	return err

}

func SetHostWithContainers(client *dockerclient.DockerClient, etcdClient *etcd.Client, node string) error {
	err := SetHost(client, etcdClient, node)
	if err != nil {
		log.Fatal(err)
	}

	err = SetContainers(client, etcdClient, node)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func getServiceResourceURI(serviceId string) string {
	return SERVICE_PREFIX + serviceId + INFO_POSTFIX
}

func getServiceContainerResourceURI(serviceId string, id string) string {
	return SERVICE_PREFIX + serviceId + "/containers/" + id
}

func SetService(etcdClient *etcd.Client, serviceId string, data *map[string]interface{}) error {

	key := getServiceResourceURI(serviceId)

	(*data)["Kind"] = "service"

	jsonData, err := json.Marshal(data)

	if err == nil {
		jsonStr := string(jsonData)
		log.Info("SetService: ", key, " -> ", jsonStr)

		_, err = etcdClient.Set(key, jsonStr, 0)
	}

	return err

}

func GetServiceDescription(etcdClient *etcd.Client, config *dockerclient.ContainerConfig) (map[string]interface{}, error) {
	serviceId := getEnvValue(config.Env, SERVICE_ID)
	key := getServiceResourceURI(serviceId)

	resp, err := etcdClient.Get(key, false, false)

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	if resp != nil {
		if err := json.Unmarshal([]byte(resp.Node.Value), &data); err != nil {
			log.Warn(err)
			return nil, err
		}
		return data, nil
	}
	//NOT FOUND
	return nil, nil
}
