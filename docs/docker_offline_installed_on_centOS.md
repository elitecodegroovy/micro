
# 安装Docker

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


## 加载Redis image

```
docker load -i redis_cluster_image.zip 

```
## 启动并创建container
```

docker run -d --restart=always --name redis-6379 --net host -v /opt/docker/redis-6379.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6380 --net host -v /opt/docker/redis-6380.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6381 --net host -v /opt/docker/redis-6381.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6382 --net host -v /opt/docker/redis-6382.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6383 --net host -v /opt/docker/redis-6383.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6384 --net host -v /opt/docker/redis-6384.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1

```




## 下载并安装 libtool-ltdl
```
rpm_url="http://mirror.centos.org/centos/7/os/x86_64/Packages/libtool-ltdl-2.4.2-22.el7_3.x86_64.rpm"
wget "${rpm_url}"
rpm -ivh --nodeps --replacefiles --replacepkgs $(basename "${rpm_url}")
```
# 下载并安装 docker-engine-selinux

```
rpm_url="https://yum.dockerproject.org/repo/main/centos/7/Packages/docker-engine-selinux-17.05.0.ce-1.el7.centos.noarch.rpm"
wget "${rpm_url}"
rpm -ivh --nodeps --replacefiles --replacepkgs $(basename "${rpm_url}")
```

# 下载并安装 docker-engine

```
rpm_url="https://yum.dockerproject.org/repo/main/centos/7/Packages/docker-engine-17.05.0.ce-1.el7.centos.x86_64.rpm"
wget "${rpm_url}"
rpm -ivh --nodeps --replacefiles --replacepkgs $(basename "${rpm_url}")

```

模版命令：

```
rpm -ivh --nodeps --replacefiles --replacepkgs ${file_name}
```

# 启动Docker

```
systemctl enable docker
systemctl start docker 

#docker version 显示服务、客户端
```


# Image Exporter And Importer

## Docker Image Save
```
docker save -o <path for generated tar file> <image name>

```

e.g.

```
docker save -o redis_cluster_image.zip jigang/redis4.0.1
```

## Docker Migration

```
docker load -i <path to image tar file>
```

e.g.

```
docker load -i redis_cluster_image.zip
```

# Install Redis Clustser

```
docker run -d --restart=always --name redis-6379 --net host -v /opt/docker/redis-6379.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6380 --net host -v /opt/docker/redis-6380.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6381 --net host -v /opt/docker/redis-6381.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6382 --net host -v /opt/docker/redis-6382.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6383 --net host -v /opt/docker/redis-6383.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6384 --net host -v /opt/docker/redis-6384.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1

```

# 配置集群服务

```
./redis-cli -h 172.29.205.18 -p 6379

#输入auth 密码
>auth TOONAN

#显示集群信息
>cluster nodes
f4772521a3586d8ca17effc1c6d39810dd1d7c12 :6379@16379 myself,master - 0 0 0 connected


#彼此通信，构建集群

172.29.205.18:6379> CLUSTER MEET 172.29.205.18 6380
OK
172.29.205.18:6379> CLUSTER MEET 172.29.205.18 6381
OK
172.29.205.18:6379> CLUSTER MEET 172.29.205.18 6382
OK
172.29.205.18:6379> CLUSTER MEET 172.29.205.18 6383
OK
172.29.205.18:6379> CLUSTER MEET 172.29.205.18 6384
OK

172.29.205.18:6379> cluster nodes
f4772521a3586d8ca17effc1c6d39810dd1d7c12 172.29.205.18:6379@16379 myself,master - 0 1568791667000 0 connected
819a358cc21e27ab36ee13afa667c6e45af16d08 172.29.205.18:6383@16383 master - 0 1568791667820 4 connected
e4206040aa7d969d817eeb76b96f5d1f2090f27d 172.29.205.18:6381@16381 master - 0 1568791668000 2 connected
0855cc5a2b1cc1ba2c97ba8f41d2363fe964360d 172.29.205.18:6382@16382 master - 0 1568791667000 3 connected
de5b3003ca7251f3a2c53dc8c44af74d551b69a5 172.29.205.18:6384@16384 master - 0 1568791668821 5 connected
60767c59b94eb308c93ec1424c60e8a136d977ca 172.29.205.18:6380@16380 master - 0 1568791669823 1 connected
```


# 分配槽信息


上面看到集群状态是失败的，原因是槽位没有分配，而且需要一次性把16384个槽位完全分配了，集群才可用。

分配槽位： CLUSTER ADDSLOTS  槽位，一个槽位只能分配一个节点，16384个槽位必须分配完，不同节点不能冲突。
所以通过脚本进行分配 do_slots_script.sh：

取代ip 192.168.1.229 为目标ip地址（redis所在机器的IP地址）：
```
#!/bin/bash
# node1 192.168.1.229   6379
n=0
for ((i=n;i<=5461;i++))
do
   /usr/docker/redis-cli -h 192.168.1.229 -p 6379 -a TOONAN  CLUSTER ADDSLOTS $i
done
 
 
# node2 1192.168.1.229    6380
n=5462
for ((i=n;i<=10922;i++))
do
   /usr/docker/redis-cli -h 192.168.1.229 -p 6380 -a TOONAN CLUSTER ADDSLOTS $i
done
 
 
# node3 192.168.1.229    6381
n=10923
for ((i=n;i<=16383;i++))
do
   /usr/docker/redis-cli -h 192.168.1.229 -p 6381 -a TOONAN CLUSTER ADDSLOTS $i
done
```

使用下面命令完成字符串取代工作。

```
sed -i "s/oldstring/newstring/g" `grep oldstring -rl yourdir`
```
其中， -a TOONAN  表示需要输入的密码。
```
192.168.1.229:6379> CLUSTER INFO
cluster_state:ok       # 集群状态为成功
cluster_slots_assigned:16384   # 已经全部分配完成
cluster_slots_ok:16384
cluster_slots_pfail:0
cluster_slots_fail:0
cluster_known_nodes:6  # 总共6个节点
cluster_size:3         # 集群为 3 个节点
cluster_current_epoch:5
cluster_my_epoch:2
cluster_stats_messages_ping_sent:1550
cluster_stats_messages_pong_sent:1509
cluster_stats_messages_meet_sent:5
cluster_stats_messages_sent:3064
cluster_stats_messages_ping_received:1509
cluster_stats_messages_pong_received:1555
cluster_stats_messages_received:3064

```

综上可知，当全部槽位分配完成之后，集群还是可行的，如果我们手欠，移除一个槽位，那么集群就立马那不行了，自己去试试吧 ——CLUSTER DELSLOTS 0 。


# 可用性集群环境

[root@svr8 docker]# ./redis-cli -h 172.29.205.18 -p 6379
172.29.205.18:6379> auth TOONAN
OK
172.29.205.18:6379> cluster nodes
f4772521a3586d8ca17effc1c6d39810dd1d7c12 172.29.205.18:6379@16379 myself,master - 0 1568792268000 0 connected 0-5461
819a358cc21e27ab36ee13afa667c6e45af16d08 172.29.205.18:6383@16383 master - 0 1568792267000 4 connected
e4206040aa7d969d817eeb76b96f5d1f2090f27d 172.29.205.18:6381@16381 master - 0 1568792270865 2 connected 10923-16383
0855cc5a2b1cc1ba2c97ba8f41d2363fe964360d 172.29.205.18:6382@16382 master - 0 1568792267860 3 connected
de5b3003ca7251f3a2c53dc8c44af74d551b69a5 172.29.205.18:6384@16384 master - 0 1568792270000 5 connected
60767c59b94eb308c93ec1424c60e8a136d977ca 172.29.205.18:6380@16380 master - 0 1568792270000 1 connected 5462-10922


编写脚本，添加副本节点

```
>i addSlaveNodes.sh 
#!/bin/bash
# replication 后面的是主master  6379的ID值， 6380的ID值，6381的ID值
/opt/docker/redis-cli -h 172.29.205.18  -p 6382 -a TOONAN  CLUSTER REPLICATE f4772521a3586d8ca17effc1c6d39810dd1d7c12

/opt/docker/redis-cli -h 172.29.205.18 -p 6383 -a TOONAN  CLUSTER REPLICATE  60767c59b94eb308c93ec1424c60e8a136d977ca

/opt/docker/redis-cli -h 172.29.205.18 -p 6384 -a TOONAN  CLUSTER REPLICATE  e4206040aa7d969d817eeb76b96f5d1f2090f27d

```
执行脚本

```
./addSlaveNodes.sh 
Warning: Using a password with '-a' option on the command line interface may not be safe.
OK
Warning: Using a password with '-a' option on the command line interface may not be safe.
OK
Warning: Using a password with '-a' option on the command line interface may not be safe.
OK


```
查看集群信息。

```
[root@svr8 docker]# ./redis-cli -h 172.29.205.18 -p 6379
172.29.205.18:6379> auth TOONAN
OK
172.29.205.18:6379> cluster
(error) ERR wrong number of arguments for 'cluster' command
172.29.205.18:6379> cluster nodes
f4772521a3586d8ca17effc1c6d39810dd1d7c12 172.29.205.18:6379@16379 myself,master - 0 1568792578000 0 connected 0-5461
819a358cc21e27ab36ee13afa667c6e45af16d08 172.29.205.18:6383@16383 slave 60767c59b94eb308c93ec1424c60e8a136d977ca 0 1568792583397 4 connected
e4206040aa7d969d817eeb76b96f5d1f2090f27d 172.29.205.18:6381@16381 master - 0 1568792580000 2 connected 10923-16383
0855cc5a2b1cc1ba2c97ba8f41d2363fe964360d 172.29.205.18:6382@16382 slave f4772521a3586d8ca17effc1c6d39810dd1d7c12 0 1568792582000 3 connected
de5b3003ca7251f3a2c53dc8c44af74d551b69a5 172.29.205.18:6384@16384 slave e4206040aa7d969d817eeb76b96f5d1f2090f27d 0 1568792582394 5 connected
60767c59b94eb308c93ec1424c60e8a136d977ca 172.29.205.18:6380@16380 master - 0 1568792581391 1 connected 5462-10922
```
上面信息说明成功创建了Redis集群。