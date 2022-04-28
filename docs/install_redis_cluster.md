## 创建redis docker基础镜像

下载redis安装包，使用版本为：4.0.1
```
[root@etcd1 tmp]# cd /tmp
[root@etcd1 tmp]# mkdir docker_redis_cluster
[root@etcd1 tmp]# cd docker_redis_cluster/
[root@etcd2 docker_redis_cluster]# wget http://download.redis.io/releases/redis-4.0.1.tar.gz
```

redis6.0.3 需要gcc5.4版本，需要升级gcc版本：
``` 
yum -y install centos-release-scl
yum -y install devtoolset-9-gcc devtoolset-9-gcc-c++ devtoolset-9-binutils
#scl命令启用只是临时的，新开的会话默认还是原gcc版本。
scl enable devtoolset-9 bash
#如果要长期使用gcc 9.1的话执行下面的命令即可：

echo -e "\nsource /opt/rh/devtoolset-9/enable" >>/etc/profile
```
解压编译redis
```
[root@etcd1 docker_redis_cluster]# tar zxvf redis-4.0.1.tar.gz
[root@etcd1 docker_redis_cluster]# cd redis-4.0.1/
[root@etcd1 redis-4.0.1]# make

```


修改redis配置

```
[root@etcd3 redis-4.0.1]# vi /tmp/docker_redis_cluster/redis-4.0.1/redis.conf
```

修改bind ip地址
```
# ~~~ WARNING ~~~ If the computer running Redis is directly exposed to the
# internet, binding to all the interfaces is dangerous and will expose the
# instance to everybody on the internet. So by default we uncomment the
# following bind directive, that will force Redis to listen only into
# the IPv4 lookback interface address (this means Redis will be able to
# accept connections only from clients running into the same computer it
# is running).
#
# IF YOU ARE SURE YOU WANT YOUR INSTANCE TO LISTEN TO ALL THE INTERFACES
# JUST COMMENT THE FOLLOWING LINE.
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
#bind 127.0.0.1
bind 0.0.0.0
```

将守护进程yes改成no
```
# By default Redis does not run as a daemon. Use 'yes' if you need it.
# Note that Redis will write a pid file in /var/run/redis.pid when daemonized.
daemonize no
```

将密码项注释去掉，添加新密码
```
# Warning: since Redis is pretty fast an outside user can try up to
# 150k passwords per second against a good box. This means that you should
# use a very strong password otherwise it will be very easy to break.
#
# requirepass foobared
```
修改为
```
# Warning: since Redis is pretty fast an outside user can try up to
# 150k passwords per second against a good box. This means that you should
# use a very strong password otherwise it will be very easy to break.
#
requirepass 123456
```
因为配置了密码，所以，配置中另外一处主从连接也需要配置密码
```
# If the master is password protected (using the "requirepass" configuration
# directive below) it is possible to tell the slave to authenticate before
# starting the replication synchronization process, otherwise the master will
# refuse the slave request.
#
# masterauth <master-password>
```
　　修改为

```
# If the master is password protected (using the "requirepass" configuration
# directive below) it is possible to tell the slave to authenticate before
# starting the replication synchronization process, otherwise the master will
# refuse the slave request.
#
# masterauth <master-password>
masterauth 123456
```
设置日志路径
```
# Specify the log file name. Also the empty string can be used to force
# Redis to log on the standard output. Note that if you use standard
# output for logging but daemonize, logs will be sent to /dev/null
logfile "/var/log/redis/redis-server.log"

```

配置集群相关信息，去掉配置项前面的注释
```
# Normal Redis instances can't be part of a Redis Cluster; only nodes that are
# started as cluster nodes can. In order to start a Redis instance as a
# cluster node enable the cluster support uncommenting the following:
#
cluster-enabled yes
 
# Every cluster node has a cluster configuration file. This file is not
# intended to be edited by hand. It is created and updated by Redis nodes.
# Every Redis Cluster node requires a different cluster configuration file.
# Make sure that instances running in the same system do not have
# overlapping cluster configuration file names.
#
cluster-config-file nodes-6379.conf
 
# Cluster node timeout is the amount of milliseconds a node must be unreachable
# for it to be considered in failure state.
# Most other internal time limits are multiple of the node timeout.
#
cluster-node-timeout 15000
````

## 镜像制作

```
[root@etcd3 docker_redis_cluster]# cd /tmp/docker_redis_cluster
[root@etcd3 docker_redis_cluster]# vi Dockerfile
# Redis
# Version 4.0.1
 
FROM centos7.5.1804
ENV REDIS_HOME /usr/local
ADD redis-4.0.1.tar.gz / # 本地的redis源码包
RUN mkdir -p $REDIS_HOME/redis # 创建安装目录
ADD redis-4.0.1/redis.conf $REDIS_HOME/redis/  # 将一开始编译产生并修改后的配置复制到安装目录
 
RUN yum -y update  # 更新yum源
RUN yum install -y gcc make # 安装编译需要的工具
 
WORKDIR /redis-4.0.1
RUN make
RUN mv /redis-4.0.1/src/redis-server  $REDIS_HOME/redis/   # 编译后，容器中只需要可执行文件redis-server
 
WORKDIR /
RUN rm -rf /redis-4.0.1          # 删除解压文件
 
RUN yum remove -y gcc make   # 安装编译完成之后，可以删除多余的gcc跟make
 
VOLUME ["/var/log/redis"]  # 添加数据卷
 
EXPOSE 6379   # 暴露6379端口，也可以暴露多个端口，这里不需要如此
```
pure version(上面注释会报错，使用下面纯粹版本):
```
# Redis
# Version 4.0.1
 
FROM centos:7
ENV REDIS_HOME /usr/local
ADD redis-4.0.1.tar.gz / 
RUN mkdir -p $REDIS_HOME/redis 
ADD redis-4.0.1/redis.conf $REDIS_HOME/redis/
 
RUN yum -y update  
RUN yum install -y gcc make 
 
WORKDIR /redis-4.0.1
RUN make
RUN mv /redis-4.0.1/src/redis-server  $REDIS_HOME/redis/ 
RUN mv /redis-4.0.1/src/redis-cli  $REDIS_HOME/redis/

WORKDIR /
RUN rm -rf /redis-4.0.1          
 
RUN yum remove -y gcc make  
 
VOLUME ["/var/log/redis"] 
 
EXPOSE 6379

```

PS.当前镜像非可执行镜像，所以没有包含ENTRYPOINT和CMD指令

镜像制作完成，制作中间可能会报： Public key for glibc-headers-2.17-222.el7.x86_64.rpm is not installed 错误，这时候需要在镜像配置中添加一句命令：

```
RUN rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7
RUN yum -y update  # 更新yum源
RUN yum install -y gcc make # 安装编译需要的工具

```
构建镜像

### 切换中国源

```
[root@etcd3 docker_redis_cluster]# vi /etc/docker/daemon.json
{
  "registry-mirrors": ["https://registry.docker-cn.com"]
}
```
 
### 编译
[root@etcd3 docker_redis_cluster]# docker build -t jigang/cluster-redis .
```
 
Complete!
 ---> 546cb1d34f35
Removing intermediate container 6b6556c5f28d
Step 14/15 : VOLUME /var/log/redis
 ---> Running in 05a6642e4046
 ---> e7e2fb8676b2
Removing intermediate container 05a6642e4046
Step 15/15 : EXPOSE 6379
 ---> Running in 5d7abe1709e2
 ---> 2d1322475f79
Removing intermediate container 5d7abe1709e2
Successfully built 2d1322475f79

```
## 制作redis节点镜像

基于此前制作的redis基础镜像创建一个redis节点镜像

```
[root@etcd3 tmp]# mkdir docker_redis_nodes
[root@etcd3 tmp]# cd docker_redis_nodes
[root@etcd3 docker_redis_nodes]# vi Dockerfile
# Redis Node
# Version 4.0.1
FROM jigang/cluster-redis:latest
 
# MAINTAINER_INFO
MAINTAINER jigang 517356906@qq.com
 
ENTRYPOINT ["/usr/local/redis/redis-server", "/usr/local/redis/redis.conf"]
```
编译构建镜像

```
docker build -t jigang/nodes-redis:4.0.1 .     
```
## 运行redis集群

``` 
mkdir /root/docker_redis_cluster
cd /root/docker_redis_cluster 
cp /tmp/docker_redis_cluster/redis-4.0.1/redis.conf redis-6379.conf
cp /tmp/docker_redis_cluster/redis-4.0.1/redis.conf redis-6380.conf
cp /tmp/docker_redis_cluster/redis-4.0.1/redis.conf redis-6381.conf
cp /tmp/docker_redis_cluster/redis-4.0.1/redis.conf redis-6382.conf
cp /tmp/docker_redis_cluster/redis-4.0.1/redis.conf redis-6383.conf
cp /tmp/docker_redis_cluster/redis-4.0.1/redis.conf redis-6384.conf
```

将redis-6380.conf、redis-6381.conf、redis-6382.conf、redis-6383.conf和redis-6384.conf 中的端口换成对应的文件名称的后缀，例如
redis-6380.conf文件中，将：

``` 
# Accept connections on the specified port, default is 6379 (IANA #815344).
# If port 0 is specified Redis will not listen on a TCP socket.
port 6379


cluster-config-file nodes-6379.conf
```
修改为：
```
# Accept connections on the specified port, default is 6379 (IANA #815344).
# If port 0 is specified Redis will not listen on a TCP socket.
port 6380

cluster-config-file nodes-6380.conf
```


需要为不同的docker程序实例设置对应的配置文件，用这个配置文件覆盖docker镜像中的配置文件。端口名称与容器名称一一对象。假设已经配置了配置redis-6379.conf、redis-6380.conf、redis-6381.conf、redis-6382.conf、redis-6383.conf和redis-6384.conf（文件中的端口与文件名称对应，已经存在目录/root/docker_redis_cluster下）。

注意
```
#
docker run -d --name redis-6379 --net host -v /tmp/docker_redis_cluster/redis-6379.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
#
docker run -d --name redis-6380 --net host -v /tmp/docker_redis_cluster/redis-6380.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
#
docker run -d --name redis-6381 --net host -v /tmp/docker_redis_cluster/redis-6381.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
#
docker run -d --name redis-6382 --net host -v /tmp/docker_redis_cluster/redis-6382.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
#
docker run -d --name redis-6383 --net host -v /tmp/docker_redis_cluster/redis-6383.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
#
docker run -d --name redis-6384 --net host -v /tmp/docker_redis_cluster/redis-6384.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1

```

生产环境下添加参数`--restart=always` ，只要容器存在就自动启动。


查看容器信息
```
#docker ps
CONTAINER ID        IMAGE                      COMMAND                  CREATED              STATUS              PORTS                    NAMES
9389a2148db3        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   36 seconds ago       Up 34 seconds       0.0.0.0:6384->6379/tcp   redis-6384
95a0cd24b19e        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   50 seconds ago       Up 48 seconds       0.0.0.0:6383->6379/tcp   redis-6383
2c44da7ea131        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   About a minute ago   Up 59 seconds       0.0.0.0:6382->6379/tcp   redis-6382
a256a2895167        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   About a minute ago   Up About a minute   0.0.0.0:6381->6379/tcp   redis-6381
0356159fc66f        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   About a minute ago   Up About a minute   0.0.0.0:6380->6379/tcp   redis-6380
64d5df3dbae6        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   About a minute ago   Up About a minute   0.0.0.0:6379->6379/tcp   redis-6379

```
运行 redis 集群容器
通过远程连接，查看redis  info replication 信息
```
redis-cli -h 192.168.1.229 -p 6379
192.168.1.229:6379> info replication
NOAUTH Authentication required.
192.168.1.229:6379> auth TOONAN
OK
192.168.1.229:6379> info replication
# Replication
role:master
connected_slaves:0
master_replid:7ecf35a3810baa06e1b8a88d55bd56b1dbf646db
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0

```

可以看到，客户连接之后，因为之前设置了密码，所以需要先输入密码认证，否则就无法通过。以上信息，我们知道所有的redis都是master角色 role:master ，这显然不是我们所希望的。


在配置之前我们需要查看所有容器当前的IP地址
```
#docker ps
CONTAINER ID        IMAGE                      COMMAND                  CREATED             STATUS              PORTS                    NAMES
9389a2148db3        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   37 minutes ago      Up 37 minutes       0.0.0.0:6384->6379/tcp   redis-6384
95a0cd24b19e        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   37 minutes ago      Up 37 minutes       0.0.0.0:6383->6379/tcp   redis-6383
2c44da7ea131        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   37 minutes ago      Up 37 minutes       0.0.0.0:6382->6379/tcp   redis-6382
a256a2895167        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   37 minutes ago      Up 37 minutes       0.0.0.0:6381->6379/tcp   redis-6381
0356159fc66f        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   37 minutes ago      Up 37 minutes       0.0.0.0:6380->6379/tcp   redis-6380
64d5df3dbae6        jigang/nodes-redis:4.0.1   "/usr/local/redis/re…"   38 minutes ago      Up 38 minutes       0.0.0.0:6379->6379/tcp   redis-6379
```

查看docker 的IP地址信息：
```
#docker inspect 9389a2148db3 95a0cd24b19e 2c44da7ea131 a256a2895167  0356159fc66f 64d5df3dbae6 |grep "IPAddress"
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.7",
                    "IPAddress": "172.17.0.7",
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.6",
                    "IPAddress": "172.17.0.6",
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.5",
                    "IPAddress": "172.17.0.5",
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.4",
                    "IPAddress": "172.17.0.4",
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.3",
                    "IPAddress": "172.17.0.3",
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.2",
                    "IPAddress": "172.17.0.2",
```
可以知道：  `redis-6379：172.17.0.2，redis-6380：172.17.0.3，redis-6381：172.17.0.4，redis-6382：172.17.0.5，redis-6383：172.17.0.6，redis-6384：172.17.0.7  `



### Redis Cluster 的集群感知操作

```
//集群(cluster) 
CLUSTER INFO 打印集群的信息 
CLUSTER NODES 列出集群当前已知的所有节点（node），以及这些节点的相关信息。  
   
//节点(node) 
CLUSTER MEET <ip> <port> 将 ip 和 port 所指定的节点添加到集群当中，让它成为集群的一份子。 
CLUSTER FORGET <node_id> 从集群中移除 node_id 指定的节点。 
CLUSTER REPLICATE <node_id> 将当前节点设置为 node_id 指定的节点的从节点。 
CLUSTER SAVECONFIG 将节点的配置文件保存到硬盘里面。  
   
//槽(slot) 
CLUSTER ADDSLOTS <slot> [slot ...] 将一个或多个槽（slot）指派（assign）给当前节点。 
CLUSTER DELSLOTS <slot> [slot ...] 移除一个或多个槽对当前节点的指派。 
CLUSTER FLUSHSLOTS 移除指派给当前节点的所有槽，让当前节点变成一个没有指派任何槽的节点。 
CLUSTER SETSLOT <slot> NODE <node_id> 将槽 slot 指派给 node_id 指定的节点，如果槽已经指派给另一个节点，那么先让另一个节点删除该槽>，然后再进行指派。 
CLUSTER SETSLOT <slot> MIGRATING <node_id> 将本节点的槽 slot 迁移到 node_id 指定的节点中。 
CLUSTER SETSLOT <slot> IMPORTING <node_id> 从 node_id 指定的节点中导入槽 slot 到本节点。 
CLUSTER SETSLOT <slot> STABLE 取消对槽 slot 的导入（import）或者迁移（migrate）。  
   
//键 (key) 
CLUSTER KEYSLOT <key> 计算键 key 应该被放置在哪个槽上。 
CLUSTER COUNTKEYSINSLOT <slot> 返回槽 slot 目前包含的键值对数量。 
CLUSTER GETKEYSINSLOT <slot> <count> 返回 count 个 slot 槽中的键。 
```


redis 集群感知：节点握手——是指一批运行在集群模式的节点通过Gossip协议彼此通信，达到感知对方的过程。

```
192.168.1.229:6379> auth TOONAN
OK
192.168.1.229:6379> CLUSTER MEET 192.168.1.229 6380
OK
192.168.1.229:6379> CLUSTER MEET 192.168.1.229 6381
OK
192.168.1.229:6379> CLUSTER MEET 192.168.1.229 6382
OK
192.168.1.229:6379> CLUSTER MEET 192.168.1.229 6383
OK
192.168.1.229:6379> CLUSTER MEET 192.168.1.229 6384
192.168.1.229:6379> CLUSTER NODES
connected

```


### 分配槽信息

查看172.17.0.2:6379 的槽个数
```
192.168.1.229:6379> cluster info
cluster_state:fail
cluster_slots_assigned:0    ## 被分配槽的个数为0
cluster_slots_ok:0
cluster_slots_pfail:0
cluster_slots_fail:0
cluster_known_nodes:6
cluster_size:0
cluster_current_epoch:5
cluster_my_epoch:2
cluster_stats_messages_ping_sent:96
cluster_stats_messages_pong_sent:109
cluster_stats_messages_meet_sent:5
cluster_stats_messages_sent:210
cluster_stats_messages_ping_received:109
cluster_stats_messages_pong_received:101
cluster_stats_messages_received:210
```

上面看到集群状态是失败的，原因是槽位没有分配，而且需要一次性把16384个槽位完全分配了，集群才可用。

分配槽位： CLUSTER ADDSLOTS  槽位，一个槽位只能分配一个节点，16384个槽位必须分配完，不同节点不能冲突。
所以通过脚本进行分配 addslots.sh：

```
#!/bin/bash
# node1 192.168.1.229   6379
n=0
for ((i=n;i<=5461;i++))
do
   /usr/local/bin/redis-cli -h 192.168.1.229 -p 6379 -a TOONAN  CLUSTER ADDSLOTS $i
done
 
 
# node2 1192.168.1.229    6380
n=5462
for ((i=n;i<=10922;i++))
do
   /usr/local/bin/redis-cli -h 192.168.1.229 -p 6380 -a TOONAN CLUSTER ADDSLOTS $i
done
 
 
# node3 192.168.1.229    6381
n=10923
for ((i=n;i<=16383;i++))
do
   /usr/local/bin/redis-cli -h 192.168.1.229 -p 6381 -a TOONAN CLUSTER ADDSLOTS $i
done
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

### 可用性集群环境

查看所有节点的id
```
192.168.1.229:6379> CLUSTER NODES
1974b62a14aea9cf9dfa219d44bc80f3ef8b14f5 192.168.1.229:6380@16379 master - 0 1536819016000 5 connected
8142e5a3826762f82f4a5980f58711e10efa9cd1 192.168.1.229:6379@16379 myself,master - 0 1536819012000 2 connected 0-5461
7e3baf5146e3c766901b797795a2918b5dbc316c 192.168.1.229:6381@16379 master - 0 1536819016298 3 connected
18e8155ed60af005a93a2dc7a15af4a42c029a78 192.168.1.229:6382@16379 master - 0 1536819015000 4 connected
ce377a8f5f499a6bf33d670b66d7e85bee0b664f 192.168.1.229:6383@16379 master - 0 1536819014000 0 connected 10923-16383
7e58dbba300a09a944c32a5f72546ef85e4afd45 192.168.1.229:6384@16379 master - 0 1536819012292 1 connected 5462-10922

```
编写脚本，添加副本节点

```
#vi addSlaveNodes.sh

#!/bin/bash
 
/bin/redis-cli -h 192.168.1.229 -p 6382 -a TOONAN  CLUSTER REPLICATE 8142e5a3826762f82f4a5980f58711e10efa9cd1（6379：node1）
 
/bin/redis-cli -h 192.168.1.229 -p 6383 -a TOONAN  CLUSTER REPLICATE 7e58dbba300a09a944c32a5f72546ef85e4afd45（6380：node2）
 
/bin/redis-cli -h 192.168.1.229 -p 6384 -a TOONAN  CLUSTER REPLICATE ce377a8f5f499a6bf33d670b66d7e85bee0b664f（6381：node3）

```
/bin/redis-cli的路径需要根据具体需要指定对应文件路径

注意：1、作为备用的节点，必须是未分配槽位的，否者会操作失败 (error) ERR To set a master the node must be empty and without assigned slots 。
           2、需要从需要添加的节点上面执行操作，CLUSTER REPLICATE [node_id]  ，使当前节点成为 node_id 的副本节点。
           3、添加从节点（集群复制）： 复制的原理和单机的Redis复制原理一样，区别是：集群下的从节点也需要运行在cluster模式下，要先添加到集群里面，再做复制。
查看所有节点信息：
```
192.168.1.229:6379> CLUSTER NODES
f65ac908e01fdd178d3c3070241df7da5842d1c6 192.168.1.229:6381@16381 master - 0 1537005672794 0 connected 10923-16383
356f3f44887cc57fd2fa0a98893e41b28e746c79 192.168.1.229:6383@16383 slave e2884a6894b15c3aaf0ef7005432af411008ac5b 0 1537005674799 4 connected
cd4f617cce5068f88f239f69235c3e9b35f7df08 192.168.1.229:6379@16379 myself,master - 0 1537005672000 1 connected 0-5461
9dbf86cd38c19a871537ab78ef63ba982c135b2e 192.168.1.229:6384@16384 slave f65ac908e01fdd178d3c3070241df7da5842d1c6 0 1537005673000 5 connected
e2884a6894b15c3aaf0ef7005432af411008ac5b 192.168.1.229:6380@16380 master - 0 1537005673797 2 connected 5462-10922
2ea62b4bb89b48dffe0e92d56dc683acce0bf940 192.168.1.229:6382@16382 slave cd4f617cce5068f88f239f69235c3e9b35f7df08 0 1537005673000 3 connected


```
可以看到我们现在实现了三主三从的一个高可用集群。



## 高可用测试——故障转移（可选项）
查看当前运行状态：

192.168.10.52:6379> CLUSTER NODES
54cb5c2eb8e5f5aed2d2f7843f75a9284ef6785c 172.17.0.3:6379@16379 master - 0 1528705604149 1 connected 5462-10922
f45f9109f2297a83b1ac36f9e1db5e70bbc174ab 172.17.0.4:6379@16379 master - 0 1528705603545 0 connected 10923-16383
ae86224a3bc29c4854719c83979cb7506f37787a 172.17.0.7:6379@16379 slave f45f9109f2297a83b1ac36f9e1db5e70bbc174ab 0 1528705603144 5 connected
98aebcfe42d8aaa8a3375e4a16707107dc9da683 172.17.0.6:6379@16379 slave 54cb5c2eb8e5f5aed2d2f7843f75a9284ef6785c 0 1528705603000 4 connected
0bbdc4176884ef0e3bb9b2e7d03d91b0e7e11f44 172.17.0.5:6379@16379 slave 760e4d0039c5ac13d04aa4791c9e6dc28544d7c7 0 1528705603000 3 connected
760e4d0039c5ac13d04aa4791c9e6dc28544d7c7 172.17.0.2:6379@16379 myself,master - 0 1528705602000 2 connected 0-5461
　　以上，运行正常

尝试关闭一个master，选择端口为6380的容器，停掉之后：


192.168.10.52:6379> CLUSTER NODES
54cb5c2eb8e5f5aed2d2f7843f75a9284ef6785c 172.17.0.3:6379@16379 master,fail - 1528706408935 1528706408000 1 connected 5462-10922
f45f9109f2297a83b1ac36f9e1db5e70bbc174ab 172.17.0.4:6379@16379 master - 0 1528706463000 0 connected 10923-16383
ae86224a3bc29c4854719c83979cb7506f37787a 172.17.0.7:6379@16379 slave f45f9109f2297a83b1ac36f9e1db5e70bbc174ab 0 1528706462980 5 connected
98aebcfe42d8aaa8a3375e4a16707107dc9da683 172.17.0.6:6379@16379 slave 54cb5c2eb8e5f5aed2d2f7843f75a9284ef6785c 0 1528706463000 4 connected
0bbdc4176884ef0e3bb9b2e7d03d91b0e7e11f44 172.17.0.5:6379@16379 slave 760e4d0039c5ac13d04aa4791c9e6dc28544d7c7 0 1528706463985 3 connected
760e4d0039c5ac13d04aa4791c9e6dc28544d7c7 172.17.0.2:6379@16379 myself,master - 0 1528706462000 2 connected 0-5461
192.168.10.52:6379>
192.168.10.52:6379> CLUSTER INFO
cluster_state:fail
cluster_slots_assigned:16384
cluster_slots_ok:10923
cluster_slots_pfail:0
cluster_slots_fail:5461
cluster_known_nodes:6
cluster_size:3
cluster_current_epoch:5
cluster_my_epoch:2
cluster_stats_messages_ping_sent:275112
cluster_stats_messages_pong_sent:274819
cluster_stats_messages_meet_sent:10
cluster_stats_messages_fail_sent:5
cluster_stats_messages_sent:549946
cluster_stats_messages_ping_received:274818
cluster_stats_messages_pong_received:275004
cluster_stats_messages_meet_received:1
cluster_stats_messages_fail_received:1
cluster_stats_messages_received:549824
　　以上，发现整个集群都失败了，从节点没有自动升级为主节点，怎么回事？？
重启停掉的容器，经排查日志信息 [root@df6ebce6f12a /]# tail -f /var/log/redis/redis-server.log  ：


1:S 11 Jun 09:57:46.712 # Cluster state changed: ok
1:S 11 Jun 09:57:46.718 * (Non critical) Master does not understand REPLCONF listening-port: -NOAUTH Authentication required.
1:S 11 Jun 09:57:46.718 * (Non critical) Master does not understand REPLCONF capa: -NOAUTH Authentication required.
1:S 11 Jun 09:57:46.719 * Partial resynchronization not possible (no cached master)
1:S 11 Jun 09:57:46.719 # Unexpected reply to PSYNC from master: -NOAUTH Authentication required.
1:S 11 Jun 09:57:46.719 * Retrying with SYNC...
1:S 11 Jun 09:57:46.719 # MASTER aborted replication with an error: NOAUTH Authentication required.
1:S 11 Jun 09:57:46.782 * Connecting to MASTER 172.17.0.6:6379
1:S 11 Jun 09:57:46.782 * MASTER <-> SLAVE sync started
1:S 11 Jun 09:57:46.782 * Non blocking connect for SYNC fired the event.

可以看到，主从之间访问需要auth，之前忘记了配置 redis.conf  中的 # masterauth <master-password> ，所以导致主从之间无法通讯。修改配置之后，自动故障转移正常。


有时候需要实施人工故障转移：

登录6380端口的从节点：6383，执行 CLUSTER FAILOVER 命令：
1
2
192.168.10.52:6383> CLUSTER  FAILOVER
(error) ERR Master is down or failed, please use CLUSTER FAILOVER FORCE
　　
发现因为master已经down了，所以我们需要执行强制转移


192.168.10.52:6383> CLUSTER FAILOVER FORCE
OK
　　
查看当前 cluster node 情况：



192.168.10.52:6383>  CLUSTER NODES
0bbdc4176884ef0e3bb9b2e7d03d91b0e7e11f44 172.17.0.5:6379@16379 slave 760e4d0039c5ac13d04aa4791c9e6dc28544d7c7 0 1528707535332 3 connected
ae86224a3bc29c4854719c83979cb7506f37787a 172.17.0.7:6379@16379 slave f45f9109f2297a83b1ac36f9e1db5e70bbc174ab 0 1528707534829 5 connected
f45f9109f2297a83b1ac36f9e1db5e70bbc174ab 172.17.0.4:6379@16379 master - 0 1528707534527 0 connected 10923-16383
98aebcfe42d8aaa8a3375e4a16707107dc9da683 172.17.0.6:6379@16379 myself,master - 0 1528707535000 6 connected 5462-10922
760e4d0039c5ac13d04aa4791c9e6dc28544d7c7 172.17.0.2:6379@16379 master - 0 1528707535834 2 connected 0-5461
54cb5c2eb8e5f5aed2d2f7843f75a9284ef6785c 172.17.0.3:6379@16379 master,fail - 1528707472833 1528707472000 1 connected
　　
从节点已经升级为master节点。这时候，我们尝试重启了，6380节点的redis（其实是重新启动停掉的容器）：


192.168.10.52:6383>  CLUSTER NODES
0bbdc4176884ef0e3bb9b2e7d03d91b0e7e11f44 172.17.0.5:6379@16379 slave 760e4d0039c5ac13d04aa4791c9e6dc28544d7c7 0 1528707556044 3 connected
ae86224a3bc29c4854719c83979cb7506f37787a 172.17.0.7:6379@16379 slave f45f9109f2297a83b1ac36f9e1db5e70bbc174ab 0 1528707555000 5 connected
f45f9109f2297a83b1ac36f9e1db5e70bbc174ab 172.17.0.4:6379@16379 master - 0 1528707556000 0 connected 10923-16383
98aebcfe42d8aaa8a3375e4a16707107dc9da683 172.17.0.6:6379@16379 myself,master - 0 1528707556000 6 connected 5462-10922
760e4d0039c5ac13d04aa4791c9e6dc28544d7c7 172.17.0.2:6379@16379 master - 0 1528707556000 2 connected 0-5461
54cb5c2eb8e5f5aed2d2f7843f75a9284ef6785c 172.17.0.3:6379@16379 slave 98aebcfe42d8aaa8a3375e4a16707107dc9da683 0 1528707556547 6 connected
　　

我们发现，6380节点反而变成了 6383节点的从节点。

现在集群应该是完整的了，所以，集群状态应该已经恢复了，我们查看下：

192.168.10.52:6383> CLUSTER INFO
cluster_state:ok
cluster_slots_assigned:16384
cluster_slots_ok:16384
cluster_slots_pfail:0
cluster_slots_fail:0
cluster_known_nodes:6
cluster_size:3
cluster_current_epoch:6
cluster_my_epoch:6
cluster_stats_messages_ping_sent:19419
cluster_stats_messages_pong_sent:19443
cluster_stats_messages_meet_sent:1
cluster_stats_messages_auth-req_sent:5
cluster_stats_messages_update_sent:1
cluster_stats_messages_sent:38869
cluster_stats_messages_ping_received:19433
cluster_stats_messages_pong_received:19187
cluster_stats_messages_meet_received:5
cluster_stats_messages_fail_received:4
cluster_stats_messages_auth-ack_received:2
cluster_stats_messages_received:38631
　　

OK，没有问题。

集群访问
客户端在初始化的时候只需要知道一个节点的地址即可，客户端会先尝试向这个节点执行命令，比如  get key ，如果key所在的slot刚好在该节点上，则能够直接执行成功。如果slot不在该节点，则节点会返回MOVED错误，同时把该slot对应的节点告诉客户端，客户端可以去该节点执行命令

192.168.10.52:6383> get hello
(error) MOVED 866 172.17.0.2:6379
 
192.168.10.52:6379> set number 20004
(error) MOVED 7743 172.17.0.3:6379
　　另外，redis集群版只使用db0，select命令虽然能够支持select 0。其他的db都会返回错误。


192.168.10.52:6383> select 0
OK
192.168.10.52:6383> select 1
(error) ERR SELECT is not allowed in cluster mode


近期，有网友询问docker redis集群连接报错的问题，具体报错如下：
```
redis.clients.jedis.exceptions.JedisConnectionException: Could not get a resource from the pool

	at redis.clients.util.Pool.getResource(Pool.java:53)
	at redis.clients.jedis.JedisPool.getResource(JedisPool.java:226)
	at redis.clients.jedis.JedisSlotBasedConnectionHandler.getConnectionFromSlot(JedisSlotBasedConnectionHandler.java:66)
	at redis.clients.jedis.JedisClusterCommand.runWithRetries(JedisClusterCommand.java:116)
	at redis.clients.jedis.JedisClusterCommand.runWithRetries(JedisClusterCommand.java:141)
```
要是单机访问。但是，实际应用化场景中，多是公网跨主机访问，问题明朗了，想着集群这东西最好还是设置成共享主机公网ip比较好，于是解决如下：
在docker运行时，执行网络模式为：host。
端口冲突解决，毕竟host模式下，容器会占用宿主机的端口，于是，我们就从配置下手，在宿主机上生成配置redis-60001.conf,redis-60002.conf,redis-60003.conf...，有多少端口建多少个文件，最终运行一个容器，挂载一个配置到容器中用于覆盖主机中的配置。
最终的运行方式如下：
```
docker run -d --name redis-6380 --net host -v /tmp/redis.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1

```



Redis 优化设置：
``` 

echo 'vm.overcommit_memory = 1' >>  /etc/sysctl.conf
echo '1024' > /proc/sys/net/core/somaxconn

cat >> /etc/security/limits.conf <<-'EOF'
* soft nofile 204800
* hard nofile 204800
* soft nproc 204800
* hard nproc 204800
EOF



# centos 7.9 versioin
cat >/etc/security/limits.d/20-nproc.conf <<-'EOF'

# Default limit for number of user's processes to prevent
# accidental fork bombs.
# See rhbz #432903 for reasoning.

*          soft    nproc     204800
*          hard    nproc     204800
EOF
```

测试服务：

./redis-cli -c -h 172.16.30.42 -p 6383

>set abc 11111
-> Redirected to slot [7638] located at 172.16.30.30:6381
(error) NOAUTH Authentication required.
172.16.30.30:6381> auth Tisson123#
OK
172.16.30.30:6381> set abc 11111
OK
172.16.30.30:6381> get abc




