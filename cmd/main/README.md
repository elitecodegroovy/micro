
# Registry Settings

```shell
# environment set for my service
MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=127.0.0.1:2379 myservice

# parameters in application cmd.
myservice --registry=etcd --registry-address=127.0.0.1:2379
```

Client :
```shell script

#Using flags
micro --client=grpc --server=grpc api
```

## List Services

```shell script
micro --registry=etcd --registry-address=127.0.0.1:2379 list services
```


## New opts

```shell script
## default srv template
micro new github.com/micro/example

## specify the template type app, wapi, fnc, srv, web
micro new github.com/micro/example1 --namespace=micro.service.example1 --alias=example1 --type=app
```

## Web Cli
The web dashboard provides a visual tool for explorings services and a built-in web proxy for web based micro services.

```shell script
# localhost:8082
micro web
```

## Health

```shell script
micro health --check_service=go.micro.srv.greeter --check_address=localhost:9090
```

check operation:
```shell script

micro health check go.micro.srv.greeter
```


## Network


```shell script

# list the nodes
micro network nodes

# list the routes
micro network routes

# list the services
micro network services

# print the graph
micro network graph
```

## Tunnel

A service tunnel is a point to point tunnel used for accessing services in remote environments.

```shell script
#Start the tunnel server (Runs on port :8083)
micro tunnel

```
y default the token “micro” is used allowing anyone to connect via the tunnel.


## Run Service

The micro runtime provides a way to manage the lifecycle of services without the complexity of orchestration systems. 

## Verify

```shell script

#Check everythings working by using a few commands
# list local services

micro list services

# list network nodes

micro network nodes

# call a service

micro call go.micro.network Debug.Health

```


## Logger

To enable debug logging
```shell script
MICRO_LOG_LEVEL=debug
```

The log levels supported are

```
trace
debug
error
info
```

To view logs from a service
```shell script
micro log [service]
```


