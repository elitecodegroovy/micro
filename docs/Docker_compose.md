### Install Docker Compose

```
$wget  https://github.com/docker/compose/releases/download/1.21.0/docker-compose-Linux-x86_64
$ mv ~/docker-compose-Linux-x86_64 /usr/local/bin/docker-compose
$ chmod a+x  /usr/local/bin/docker-compose
$ cd harbor
$ docker load -i harbor.v1.1.0.tar.gz
```