# 部署TiDB


## TiDB安装

在测试环境，部署6台独立ip的机器，如下所示：

| Host Name | IP            | Services   | Data Path      |
| --------- | ------------- | ---------- | ------------   |
| **host1** | 192.168.1.229 | PD1 & TiDB | /home/tidb     |
| **host2** | 192.168.1.221 | PD2        | /home/tidb     |
| **host3** | 192.168.1.22  | PD3        | /home/tidb     |
| **host4** | 192.168.1.147 | TiKV1      | /data/tidb     |
| **host5** | 192.168.1.155 | TiKV2      | /data/tidb     |
| **host6** | 192.168.1.164 | TiKV3      | /data/tidb     |

### 安装Docker

从已有文件中，解压缩docker.zip得到docker的安装文件。执行项目命令：
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
  --advertise-client-urls="http://192.168.1.229:2379" \
  --peer-urls="http://0.0.0.0:2380" \
  --advertise-peer-urls="http://192.168.1.229:2380" \
  --initial-cluster="pd1=http://192.168.1.229:2380,pd2=http://192.168.1.221:2380,pd3=http://192.168.1.22:2380"
  
  
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
  --advertise-client-urls="http://192.168.1.221:2379" \
  --peer-urls="http://0.0.0.0:2380" \
  --advertise-peer-urls="http://192.168.1.221:2380" \
  --initial-cluster="pd1=http://192.168.1.229:2380,pd2=http://192.168.1.221:2380,pd3=http://192.168.1.22:2380"
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
  --advertise-client-urls="http://192.168.1.22:2379" \
  --peer-urls="http://0.0.0.0:2380" \
  --advertise-peer-urls="http://192.168.1.22:2380" \
  --initial-cluster="pd1=http://192.168.1.229:2380,pd2=http://192.168.1.221:2380,pd3=http://192.168.1.22:2380"
 ```
 
 TiKV1:
 
 ```
 docker run --restart=always -d --name tikv1 \
  -p 20160:20160 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /data/tidb:/data \
  pingcap/tikv:latest \
  --addr="0.0.0.0:20160" \
  --advertise-addr="192.168.1.147:20160" \
  --data-dir="/data/tikv1" \
  --pd="192.168.1.229:2379,192.168.1.221:2379,192.168.1.22:2379"
 ```
 
  TiKV2:
  ```
 docker run --restart=always -d --name tikv2 \
  -p 20160:20160 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /data/tidb:/data \
  pingcap/tikv:latest \
  --addr="0.0.0.0:20160" \
  --advertise-addr="192.168.1.155:20160" \
  --data-dir="/data/tikv2" \
  --pd="192.168.1.229:2379,192.168.1.221:2379,192.168.1.22:2379"
 ```
 
 
   TiKV3:
   
  ```
 docker run --restart=always -d --name tikv3 \
  -p 20160:20160 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /data/tidb:/data \
  pingcap/tikv:latest \
  --addr="0.0.0.0:20160" \
  --advertise-addr="192.168.1.164:20160" \
  --data-dir="/data/tikv3" \
  --pd="192.168.1.229:2379,192.168.1.221:2379,192.168.1.22:2379"
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

The TiKV and PD can be started with a specified configuration file, which includes some advanced parameters, for the performance tuning.

Assume that the path to configuration file of PD and TiKV on the host is /path/to/config/pd.toml and /path/to/config/tikv.toml

You can start TiKV and PD as follows:

```
docker run -d --name tikv1 \
  -p 20160:20160 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /data:/data \
  -v /path/to/config/tikv.toml:/tikv.toml:ro \
  pingcap/tikv:latest \
  --addr="0.0.0.0:20160" \
  --advertise-addr="192.168.1.104:20160" \
  --data-dir="/data/tikv1" \
  --pd="192.168.1.101:2379,192.168.1.102:2379,192.168.1.103:2379" \
  --config="/tikv.toml"
```


```
docker run -d --name pd1 \
  -p 2379:2379 \
  -p 2380:2380 \
  -v /etc/localtime:/etc/localtime:ro \
  -v /data:/data \
  -v /path/to/config/pd.toml:/pd.toml:ro \
  pingcap/pd:latest \
  --name="pd1" \
  --data-dir="/data/pd1" \
  --client-urls="http://0.0.0.0:2379" \
  --advertise-client-urls="http://192.168.1.101:2379" \
  --peer-urls="http://0.0.0.0:2380" \
  --advertise-peer-urls="http://192.168.1.101:2380" \
  --initial-cluster="pd1=http://192.168.1.101:2380,pd2=http://192.168.1.102:2380,pd3=http://192.168.1.103:2380" \
  --config="/pd.toml"
  
 ```


# 迁移镜像


##Docker Image Save

docker save -o <path for generated tar file> <image name>

e.g.

```
docker save -o pd.zip pingcap/pd:latest
```

## Docker Migration

docker load -i <path to image tar file>

e.g.
```
docker load -i redis_cluster_image.zip 
```
