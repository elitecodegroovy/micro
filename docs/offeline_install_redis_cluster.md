
# 离线安装Redis集群

```
wget http://download.redis.io/releases/redis-4.0.1.tar.gz

tar xvf redis-4.0.1.tar.gz
cd redis-4.0.1
make

#复制redis-server到目录/opt/redis_cluster

cp src/redis-server /opt/redis_cluster/
```

将准备好的redis集群conf文件复制到/opt/redis_cluster/，编写下面脚本start-redisCluster.sh：

```
/opt/redis_cluster/redis-server /opt/redis_cluster/redis-6380.conf &
/opt/redis_cluster/redis-server /opt/redis_cluster/redis-6381.conf &
/opt/redis_cluster/redis-server /opt/redis_cluster/redis-6382.conf &
/opt/redis_cluster/redis-server /opt/redis_cluster/redis-6383.conf &
/opt/redis_cluster/redis-server /opt/redis_cluster/redis-6384.conf &
/opt/redis_cluster/redis-server /opt/redis_cluster/redis-6385.conf &
```


``
redis-cli -h 10.124.132.8 -p 6380
>CLUSTER MEET 10.124.132.8 6381
>CLUSTER MEET 10.124.132.8 6382
>CLUSTER MEET 10.124.132.8 6383
>CLUSTER MEET 10.124.132.8 6384
>CLUSTER MEET 10.124.132.8 6385
>cluster nodes
>
bfa48ce85fecc1f48e327fc11bcc067b1ddbb438 10.124.132.8:6384@16384 master - 0 1560669575254 4 connected
da49a70e92676124cdfdd030b49830edb936d0d2 10.124.132.8:6381@16381 master - 0 1560669574000 0 connected
a6befa5f1f51a4f536c339dc3c825ad16a796805 10.124.132.8:6380@16380 myself,master - 0 1560669575000 1 connected
b9bb7ac525d52734b3e4b299db7c3d1394e8aacf 10.124.132.8:6383@16383 master - 0 1560669574000 3 connected
8d7cf6d8b83f70609b048bfa42be8d29324ec880 10.124.132.8:6382@16382 master - 0 1560669576256 2 connected
757b12d6368f7a2de5c8f94c158d229daa901b0e 10.124.132.8:6385@16385 master - 0 1560669574250 5 connected
```

运行添加slot，执行下面脚本：

```
#!/bin/bash
# node1 10.124.132.8   6380
n=0
for ((i=n;i<=5461;i++))
do
   /opt/redis_cluster/redis-cli -h 10.124.132.8 -p 6380 -a TOONAN  CLUSTER ADDSLOTS $i
done


# node2 10.124.132.8    6381
n=5462
for ((i=n;i<=10922;i++))
do
   /opt/redis_cluster/redis-cli -h 10.124.132.8 -p 6381 -a TOONAN CLUSTER ADDSLOTS $i
done


# node3 10.124.132.8    6382
n=10923
for ((i=n;i<=16383;i++))
do
   /opt/redis_cluster/redis-cli -h 10.124.132.8 -p 6382 -a TOONAN CLUSTER ADDSLOTS $i
done
```

其中， -a TOONAN  表示需要输入的密码。
```
192.168.1.229:6379> CLUSTER NODES
bfa48ce85fecc1f48e327fc11bcc067b1ddbb438 10.124.132.8:6384@16384 master - 0 1560671583295 4 connected
da49a70e92676124cdfdd030b49830edb936d0d2 10.124.132.8:6381@16381 master - 0 1560671584299 0 connected 5462-10922
a6befa5f1f51a4f536c339dc3c825ad16a796805 10.124.132.8:6380@16380 myself,master - 0 1560671585000 1 connected 0-5461
b9bb7ac525d52734b3e4b299db7c3d1394e8aacf 10.124.132.8:6383@16383 master - 0 1560671583000 3 connected
8d7cf6d8b83f70609b048bfa42be8d29324ec880 10.124.132.8:6382@16382 master - 0 1560671584000 2 connected 10923-16383
757b12d6368f7a2de5c8f94c158d229daa901b0e 10.124.132.8:6385@16385 master - 0 1560671585302 5 connected

```

添加slave实例。执行下面脚本：

```
#!/bin/bash
# replication 后面的是主master
/opt/redis_cluster/redis-cli -h 10.124.132.8 -p 6383 -a TOONAN  CLUSTER REPLICATE a6befa5f1f51a4f536c339dc3c825ad16a796805

/opt/redis_cluster/redis-cli -h 10.124.132.8 -p 6384 -a TOONAN  CLUSTER REPLICATE da49a70e92676124cdfdd030b49830edb936d0d2

/opt/redis_cluster/redis-cli -h 10.124.132.8 -p 6385 -a TOONAN  CLUSTER REPLICATE 8d7cf6d8b83f70609b048bfa42be8d29324ec880

```

注意REPLICATE后的值就是master节点的标识符，根据cluster node信息进行修改。显示下面信息表示安装成功了。

```

10.124.132.8:6380> clluster nodes
(error) ERR unknown command 'clluster'
10.124.132.8:6380> cluster nodes
bfa48ce85fecc1f48e327fc11bcc067b1ddbb438 10.124.132.8:6384@16384 slave da49a70e92676124cdfdd030b49830edb936d0d2 0 1560672029738 4 connected
da49a70e92676124cdfdd030b49830edb936d0d2 10.124.132.8:6381@16381 master - 0 1560672031000 0 connected 5462-10922
a6befa5f1f51a4f536c339dc3c825ad16a796805 10.124.132.8:6380@16380 myself,master - 0 1560672030000 1 connected 0-5461
b9bb7ac525d52734b3e4b299db7c3d1394e8aacf 10.124.132.8:6383@16383 slave a6befa5f1f51a4f536c339dc3c825ad16a796805 0 1560672031000 3 connected
8d7cf6d8b83f70609b048bfa42be8d29324ec880 10.124.132.8:6382@16382 master - 0 1560672031744 2 connected 10923-16383
757b12d6368f7a2de5c8f94c158d229daa901b0e 10.124.132.8:6385@16385 slave 8d7cf6d8b83f70609b048bfa42be8d29324ec880 0 1560672032746 5 connected
```



