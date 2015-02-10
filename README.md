## Commander: a Commander of Docker clustering system



### Installation

######1 - Download and install the current source code.
Ensure you have golang and git client installed (e.g. `apt-get install golang git` on Ubuntu).
You may need to set `$GOPATH`, e.g `mkdir ~/gocode; export GOPATH=~/gocode`.


```sh
go get -u github.com/denverdino/commander
```


######2 - Docker Compose setup


```sh
git clone https://github.com/denverdino/fig.git
cd fig
python setup.py develop
```


######3 - Etcd Setup in boot2docker

Setup the Etcd with web console for debugging
ETCD_HOST
Please change the 

```sh
docker run -d -p 4001:4001 jadetest.cn.ibm.com:5000/coreos/etcd:v0.4.6 -cors='*'
docker run -d -p 8000:8000 --env ETCD_HOST=192.168.59.103 jadetest.cn.ibm.com:5000/pure/etcd-browser
```

######4 - "commander" of Container Service

```sh
commander manage  -H 0.0.0.0:4243 -node 192.168.59.103 -addr 192.168.59.103:2375 -etcd http://192.168.59.103:4001
```



######4 - Node watcher


```sh
commander watch  -node 192.168.59.103 -addr 192.168.59.103:2375 -etcd http://192.168.59.103:4001
```