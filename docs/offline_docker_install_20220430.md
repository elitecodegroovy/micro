
## Install Docker

18 版本：
```
rpm -ivh --nodeps --replacefiles --replacepkgs containerd.io-1.2.2-3.3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-selinux-17.03.3.ce-1.el7.noarch.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-cli-18.09.2-3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-18.09.2-3.el7.x86_64.rpm
```

19 版本：
``` 
rpm -ivh  -ivh --nodeps --replacefiles --replacepkgs docker-ce-cli-19.03.13-3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs --force selinux-policy-3.13.1-268.el7.noarch.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs --force selinux-policy-targeted-3.13.1-268.el7.noarch.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs container-selinux-2.119.2-1.911c772.el7_8.noarch.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-19.03.13-3.el7.x86_64.rpm
```

## Configure Docker startup

```
systemctl enable docker 
systemctl start docker 
docker version
```
Show following information that proves it does successfully. 

```
Client:
 Version:           18.09.2
 API version:       1.39
 Go version:        go1.10.6
 Git commit:        6247962
 Built:             Sun Feb 10 04:13:27 2019
 OS/Arch:           linux/amd64
 Experimental:      false

Server: Docker Engine - Community
 Engine:
  Version:          18.09.2
  API version:      1.39 (minimum version 1.12)
  Go version:       go1.10.6
  Git commit:       6247962
  Built:            Sun Feb 10 03:47:25 2019
  OS/Arch:          linux/amd64
  Experimental:     false
```

## Import Image File

```
docker load -i mysql_docker_image.zip

```
Create user and account group：

```
groupadd mysql
useradd -g mysql mysql
```


Create the data storage directory `/opt/docker_mysql/data`:

```
mkdir -p /opt/docker_mysql/data
chown -R mysql:mysql /opt/docker_mysql
```

## Run Container

```
docker run --restart=always -d --name mysqld  \
-e MYSQL_ROOT_PASSWORD=TOONAN123456  \
-e MYSQL_MAX_CONNECTIONS=320 \
-v /opt/docker_mysql/data:/var/lib/mysql/data \
--net centos/mysql-57-centos7:5.7
```