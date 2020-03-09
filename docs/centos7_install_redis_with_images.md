

假设redis的集群配置文件在目录`/opt/redis_cluster`下。

```
docker run -d --restart=always --name redis-6379 --net host -v /opt/redis_cluster/redis-6379.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6380 --net host -v /opt/redis_cluster/redis-6380.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6381 --net host -v /opt/redis_cluster/redis-6381.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6382 --net host -v /opt/redis_cluster/redis-6382.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6383 --net host -v /opt/redis_cluster/redis-6383.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
docker run -d --restart=always --name redis-6384 --net host -v /opt/redis_cluster/redis-6384.conf:/usr/local/redis/redis.conf  jigang/nodes-redis:4.0.1
```


redis 集群感知：节点握手——是指一批运行在集群模式的节点通过Gossip协议彼此通信，达到感知对方的过程。假设提供的主机IP是`192.168.1.101`（依据服务的IP而变化）
```
chmod +x *
redis-cli -h 192.168.1.101 -p 6379
192.168.1.101>CLUSTER MEET 192.168.1.101 6380
192.168.1.101>CLUSTER MEET 192.168.1.101 6381
192.168.1.101>CLUSTER MEET 192.168.1.101 6382
192.168.1.101>CLUSTER MEET 192.168.1.101 6383
192.168.1.101>CLUSTER MEET 192.168.1.101 6384
#显示集群信息
192.168.1.101>CLUSTER NODES
40810a3c97f28e1c6bd4239a12f90b4d324e723b 192.168.1.101:6383@16383 master - 0 1550896365928 4 connected
74796c70d932bca41320d6be34bc172a19302258 192.168.1.101:6380@16380 master - 0 1550896365000 2 connected 5462-10922
d0cab0bc0bc731495248bfe338c1394874cbc8c6 192.168.1.101:6382@16382 master - 0 1550896364924 0 connected
50a882c3b20f852a666b866cbeeeb3a6c1926abe 192.168.1.101:6384@16384 master - 0 1550896362000 5 connected
107277df1b546a95120bfb034e108f20fe8b1148 192.168.1.101:6379@16379 myself,master - 0 1550896365000 1 connected 0-5461
78ebae87007f4be13161efde4708721c805e4132 192.168.1.101:6381@16381 master - 0 1550896366930 3 connected 10923-16383

192.168.1.101>exit
```


分配槽信息，假设提供的主机IP是`192.168.1.101`（依据服务的IP而变化）。修改do_slots_script.sh脚本。确定IP地址和redis-cli所在的位置。修改成功后，执行命令。
```
chmod +x *
./do_slots_script.sh
```

**编写脚本*addSlaveNodes.sh***，添加副本节点。
```
./addSlaveNodes.sh

```

显示下面信息表示成功了：
```
Warning: Using a password with '-a' option on the command line interface may not be safe.
OK
Warning: Using a password with '-a' option on the command line interface may not be safe.
OK
Warning: Using a password with '-a' option on the command line interface may not be safe.
OK
```

```
CREATE USER 'gc'@'%' IDENTIFIED BY '019#^(963';

GRANT ALL PRIVILEGES ON *.* TO 'gc'@'%' identified by '019#^(963';

```