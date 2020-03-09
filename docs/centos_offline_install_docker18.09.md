### 安装Docker

从已有文件中，解压缩docker.zipo得到docker的安装文件。执行项目命令：
```
rpm -ivh --nodeps --replacefiles --replacepkgs libseccomp-2.3.1-3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs libseccomp-devel-2.3.1-3.el7.x86_64.rpm 
rpm -ivh --nodeps --replacefiles --replacepkgs libltdl7-2.4.2-alt8.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-selinux-17.03.3.ce-1.el7.noarch.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-cli-18.09.2-3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs containerd.io-1.2.2-3.3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-18.09.2-3.el7.x86_64.rpm
```

启动docker服务：

```
systemctl start docker
```

开机启动docker服务：

```
systemctl enable docker
```

显示启动docker版本信息：
```
[root@node5 Downloads]# docker version
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


# 加载Redis image

```
docker load -i redis_cluster_image.zip 

```
# 启动并创建container
```

docker run -d --restart=always --name redis-6379 --net host -v /opt/docker/redis-6379.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6380 --net host -v /opt/docker/redis-6380.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6381 --net host -v /opt/docker/redis-6381.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6382 --net host -v /opt/docker/redis-6382.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6383 --net host -v /opt/docker/redis-6383.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6384 --net host -v /opt/docker/redis-6384.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1

```

