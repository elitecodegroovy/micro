
### Setting up a cluster on 3 different servers provisioned beforehand

server 1

```
# Cluster Server A192.168.1.224

port: 6222

cluster {
  host: '0.0.0.0'
  port: 6223

  routes = [
    nats-route://192.168.1.225:6223
    nats-route://192.168.1.226:6223
  ]
}
```

server 2
```
# Cluster Server B 192.168.1.225
# 
port: 6222

cluster {
  host: '0.0.0.0'
  port: 6223

  routes = [
    nats-route://192.168.1.224:6223
    nats-route://192.168.1.226:6223
  ]
}

```

server 3
```
# Cluster Server C 192.168.1.226
port: 6222

cluster {
  host: '0.0.0.0'
  port: 6223

  routes = [
    nats-route://192.168.1.224:6223
    nats-route://192.168.1.225:6223
  ]
}
```

Execute Comand on each node:

```
gmessaged -config ./gmessage.conf -D
```