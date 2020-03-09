
172.29.205.14	M:64G	tikv1(data-node)
172.29.205.15	M:64G	tikv2(data-node)
172.29.205.16	M:64G	tikv3(data-node)
172.29.205.17	M:128G	PD3
172.29.205.13	M:64G	PD2
172.29.205.18	M:128G	PD1、TiDB


### 启动TiDB

启动TiDB的pd、tikv和tidb组件。 启动PD1、PD2、PD3、TiKV1，TiKV2，TiKV3.

PD1:

```
docker run --restart=always -d --name pd1  \
  -p 2379:2379 \
  -p 2380:2380 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /home/tidb:/data \
  pingcap/pd:latest \
  --name="pd1" \
  --data-dir="/data/pd1" \
  --client-urls="http://0.0.0.0:2379" \
  --advertise-client-urls="http://http://172.29.205.18:2379" \
  --peer-urls="http://0.0.0.0:2380" \
  --advertise-peer-urls="http://http://172.29.205.18:2380" \
  --initial-cluster="pd1=http://172.29.205.18:2380,pd2=http://172.29.205.13:2380,pd3=http://172.29.205.17:2380"
  
  
  ```
  
 PD2:
 
  ```
  docker run  --restart=always -d --name pd2 \
  -p 2379:2379 \
  -p 2380:2380 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /home/tidb:/data \
  pingcap/pd:latest \
  --name="pd2" \
  --data-dir="/data/pd2" \
  --client-urls="http://0.0.0.0:2379" \
  --advertise-client-urls="http://172.29.205.13:2379" \
  --peer-urls="http://0.0.0.0:2380" \
  --advertise-peer-urls="http://172.29.205.13:2380" \
  --initial-cluster="pd1=http://172.29.205.18:2380,pd2=http://172.29.205.13:2380,pd3=http://172.29.205.17:2380"
  ```
  
PD3
  
 ```
 docker run --restart=always -d --name pd3 \
  -p 2379:2379 \
  -p 2380:2380 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /home/tidb:/data \
  pingcap/pd:latest \
  --name="pd3" \
  --data-dir="/data/pd3" \
  --client-urls="http://0.0.0.0:2379" \
  --advertise-client-urls="http://172.29.205.17:2379" \
  --peer-urls="http://0.0.0.0:2380" \
  --advertise-peer-urls="http://172.29.205.17:2380" \
  --initial-cluster="pd1=http://172.29.205.18:2380,pd2=http://172.29.205.13:2380,pd3=http://172.29.205.17:2380"
 
 ```
 
 TiKV1:
 
 ```
 docker run --restart=always -d --name tikv1 \
  -p 20160:20160 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /home/tidb/data:/data \
  pingcap/tikv:latest \
  --addr="0.0.0.0:20160" \
  --advertise-addr="192.168.1.14:20160" \
  --data-dir="/data/tikv1" \
  --pd="172.29.205.18:2379,172.29.205.13:2379,172.29.205.17:2379"
 ```
 
  TiKV2:
  ```
 docker run --restart=always -d --name tikv2 \
  -p 20160:20160 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /home/tidb/data:/data \
  pingcap/tikv:latest \
  --addr="0.0.0.0:20160" \
  --advertise-addr="192.168.1.15:20160" \
  --data-dir="/data/tikv2" \
  --pd="172.29.205.18:2379,172.29.205.13:2379,172.29.205.17:2379"
 ```
 
 
   TiKV3:
   
  ```
 docker run --restart=always -d --name tikv3 \
  -p 20160:20160 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /home/tidb/data:/data \
  pingcap/tikv:latest \
  --addr="0.0.0.0:20160" \
  --advertise-addr="192.168.1.16:20160" \
  --data-dir="/data/tikv3" \
  --pd="172.29.205.18:2379,172.29.205.13:2379,172.29.205.17:2379"
 ```
 
 启动TiDB：
 
 ```
 docker run --restart=always -d --name tidb \
  -p 4000:4000 \
  -p 10080:10080 \
  -v /etc/localtime:/etc/localtime:ro \
  pingcap/tidb:latest \
  --store=tikv \
  --path="192.168.1.229:2379,192.168.1.221:2379,192.168.1.22:2379"
 ```
 
 
 
 ### 测试TiDB
 
 使用mysql客户端连接TiDB服务。
```
 mysql -h 127.0.0.1 -P 4000 -u root -D test
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| INFORMATION_SCHEMA |
| PERFORMANCE_SCHEMA |
| mysql              |
| test               |
+--------------------+
4 rows in set (0.00 sec)

```


change account's passwd:
```

SET PASSWORD FOR 'root'@'%' = '##%TOONAN2019)#@*';

```

Add the account :
```
CREATE USER 'bigdata'@'%' IDENTIFIED BY 'BI%Gdata#999@!*@@@@@';
GRANT ALL PRIVILEGES ON bigdata.* TO 'bigdata'@'%' WITH GRANT OPTION;

```