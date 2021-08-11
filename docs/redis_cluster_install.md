
进入/opt/docker_redic_cluster目录，执行下面命令

```
[root@etcd1 tmp]# mkdir -p /opt/docker_redis_cluster
[root@etcd1 tmp]# cd /opt/docker_redis_cluster/
[root@etcd2 docker_redis_cluster]# wget http://download.redis.io/releases/redis-6.0.9.tar.gz
```


解压编译redis
```
[root@etcd1 docker_redis_cluster]# tar zxvf redis-6.0.9.tar.gz
[root@etcd1 docker_redis_cluster]# cd redis-6.0.9/
[root@etcd1 redis-4.0.1]# make

```
                                                                 

修改redis配置

```
[root@etcd3 redis-4.0.1]# vi /tmp/docker_redis_cluster/redis-6.0.9/redis.conf
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

进入/opt/docker_redis_cluster目录，创建Dockerfile文件，编辑内容如下：
```
# Redis
# Version 6.0.9

FROM centos:7.5.1804
ENV REDIS_HOME /opt/docker_redis_cluster
ADD redis-6.0.9.tar.gz /
RUN mkdir -p $REDIS_HOME/redis 
ADD redis-6.0.9/redis.conf $REDIS_HOME/redis/ 

RUN yum -y update  
RUN yum install -y gcc make 

WORKDIR /redis-6.0.9
RUN make
RUN mv /redis-6.0.9/src/redis-server  $REDIS_HOME/redis/  

WORKDIR /
RUN rm -rf /redis-6.0.9          

RUN yum remove -y gcc make  

VOLUME ["/var/log/redis"] 

EXPOSE 6379  
```

``` 
# Redis
# Version 6.0.9

FROM centos:7.5.1804
ENV REDIS_HOME /opt/docker_redis_cluster
RUN mkdir -p $REDIS_HOME/redis 
ADD redis.conf $REDIS_HOME/redis/ 

ADD src/redis-server  $REDIS_HOME/redis/  

VOLUME ["/var/log/redis"] 

EXPOSE 6379

```
构建镜像

### 切换中国源

```
[root@etcd3 docker_redis_cluster]# vi /etc/docker/daemon.json
{
  "registry-mirrors": ["https://registry.docker-cn.com"]
}
```

运行编译：
```
docker build -t jigang/docker-redis .	
```
编辑定义的镜像：
```
# Redis Node
# Version 6.0.9
FROM jigang/docker-redis:latest

# MAINTAINER_INFO
MAINTAINER jigang 517356906@qq.com

ENTRYPOINT ["/opt/docker_redis_cluster/redis/redis-server", "/opt/docker_redis_cluster/redis/redis.conf"]

```

构建nodes-redis:6.0.9镜像。
```
docker build -t jigang/nodes-redis:6.0.9 .	
```

```
docker build -t jigang/nodes-redis:6.0.9 .	
```

修改配置文件redis-6379.conf，redis-6380.conf，redis-6381.conf,redis-6382.conf,redis-6383.conf,redis-6384.conf的端口号。例如，redis-6383.conf的配置如下：






创建容器并启动：
						
```
docker run -d --name redis-6379 --net host -v /opt/docker_redis_cluster/redis-6379.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --name redis-6380 --net host -v /opt/docker_redis_cluster/redis-6380.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --name redis-6381 --net host -v /opt/docker_redis_cluster/redis-6381.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --name redis-6382 --net host -v /opt/docker_redis_cluster/redis-6382.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --name redis-6383 --net host -v /opt/docker_redis_cluster/redis-6383.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --name redis-6384 --net host -v /opt/docker_redis_cluster/redis-6384.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
```

jigang/nodes-redis:4.0.1 的镜像，地址变成`/root/docker_redis_cluster/redis-6379.conf:/usr/local/redis/redis.conf`. 
alway running :

``` 
--restart=always

e.g.
docker run -d --restart=always --name redis-6379 --net host -v /opt/docker_redis_cluster/redis-6379.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --restart=always --name redis-6380 --net host -v /opt/docker_redis_cluster/redis-6380.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --restart=always --name redis-6381 --net host -v /opt/docker_redis_cluster/redis-6381.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --restart=always --name redis-6382 --net host -v /opt/docker_redis_cluster/redis-6382.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --restart=always --name redis-6383 --net host -v /opt/docker_redis_cluster/redis-6383.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
docker run -d --restart=always --name redis-6384 --net host -v /opt/docker_redis_cluster/redis-6384.conf:/opt/docker_redis_cluster/redis/redis.conf  jigang/nodes-redis:6.0.9
```


```
#!/bin/bash
# node1 172.16.27.24   6379
n=0
for ((i=n;i<=5461;i++))
do
   /opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6379 -a Tisson123!  CLUSTER ADDSLOTS $i
done


# node2 172.16.27.24    6380
n=5462
for ((i=n;i<=10922;i++))
do
   /opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6380 -a Tisson123! CLUSTER ADDSLOTS $i
done


# node3 172.16.27.24    6381
n=10923
for ((i=n;i<=16383;i++))
do
   /opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6381 -a Tisson123! CLUSTER ADDSLOTS $i
done
```

redis 集群感知：节点握手——是指一批运行在集群模式的节点通过 Gossip 协议彼此通信， 达到感知对方的过程。

```
/opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6382

192.168.1.229:6379> auth Tisson123!
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

192.168.1.229:6379> CLUSTER NODES connected

```


查看 172.17.0.2:6379 的槽个数
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

6382对应是6379的ID标识，以此类推。编写脚本，添加副本节点:

```
#!/bin/bash
# replication 后面的是主master
/opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6382 -a Tisson123!  CLUSTER REPLICATE 8142e5a3826762f82f4a5980f58711e10efa9cd1

/opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6383 -a Tisson123!  CLUSTER REPLICATE 1974b62a14aea9cf9dfa219d44bc80f3ef8b14f5

/opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6384 -a Tisson123!  CLUSTER REPLICATE 7e3baf5146e3c766901b797795a2918b5dbc316c
```

``` 
e.g.
172.16.27.24:6379> cluster nodes
842fbcd90a5943b3369da432e1fa66114b9280c3 172.16.27.24:6383@16383 master - 0 1604308063526 4 connected
0cdc30904f0333f57169fdb82ef92bd33c31ebfb 172.16.27.24:6380@16380 master - 0 1604308062521 1 connected 5462-10922
b3d3f6aaddb7e272c1a2457ab8297fa97d6de958 172.16.27.24:6384@16384 master - 0 1604308062000 5 connected
ccc3b83574a60016d65625b83b4eb90f7d78ad48 172.16.27.24:6382@16382 master - 0 1604308062000 0 connected
a2310bb922c3cc4f2202ff586c204d6df37b6da2 172.16.27.24:6379@16379 myself,master - 0 1604308063000 2 connected 0-5461
a5f0e4e3f8ee32e9fb4bbf7216ad6933e4179d00 172.16.27.24:6381@16381 master - 0 1604308062000 3 connected 10923-16383
172.16.27.24:6379> 
e.g.
/opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6382 -a Tisson123!  CLUSTER REPLICATE a2310bb922c3cc4f2202ff586c204d6df37b6da2
/opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6383 -a Tisson123!  CLUSTER REPLICATE 0cdc30904f0333f57169fdb82ef92bd33c31ebfb
/opt/docker_redis_cluster/src/redis-cli -h 172.16.27.24 -p 6384 -a Tisson123!  CLUSTER REPLICATE a5f0e4e3f8ee32e9fb4bbf7216ad6933e4179d00
```