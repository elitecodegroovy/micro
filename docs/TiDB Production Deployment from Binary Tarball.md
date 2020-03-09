
##Delopy TiDB


### PD

172.29.205.14		data-node
172.29.205.15		data-node
172.29.205.16		data-node
172.29.205.17		PD2
172.29.205.13		PD3
172.29.205.18		PD1、TiDB


PD1:
```
./bin/pd-server --name=pd1 \
                --data-dir=/home/tidb/pd/pd \
                --client-urls="http://172.29.205.18:2379" \
                --peer-urls="http://172.29.205.18:2380" \
                --initial-cluster="pd1=http://172.29.205.18:2380,pd3=http://172.29.205.13:2380,pd2=http://172.29.205.17:2380" \
                -L "info" \
                --log-file=pd.log &
				
				
./bin/pd-server --name=pd2 \
                --data-dir=/home/tidb/pd/pd \
                --client-urls="http://172.29.205.17:2379" \
                --peer-urls="http://172.29.205.17:2380" \
                --initial-cluster="pd1=http://172.29.205.18:2380,pd3=http://172.29.205.13:2380,pd2=http://172.29.205.17:2380" \
                -L "info" \
                --log-file=pd.log &
				
./bin/pd-server --name=pd3 \
                --data-dir=/home/tidb/pd/pd \
                --client-urls="http://172.29.205.13:2379" \
                --peer-urls="http://172.29.205.13:2380" \
                --initial-cluster="pd1=http://172.29.205.18:2380,pd3=http://172.29.205.13:2380,pd2=http://172.29.205.17:2380" \
                -L "info" \
                --log-file=pd.log &
```



### TiKV

TiKV:

```
#172.29.205.14
$ ./bin/tikv-server --pd="172.29.205.18:2379,172.29.205.13:2379,172.29.205.17:2379" \
                  --addr="172.29.205.14:20160" \
                  --status-addr="172.29.205.14:20180" \
                  --data-dir=tikv \
                  --log-file=tikv.log &

#172.29.205.15
$ ./bin/tikv-server --pd="172.29.205.18:2379,172.29.205.13:2379,172.29.205.17:2379" \
                  --addr="172.29.205.15:20160" \
                  --status-addr="172.29.205.15:20180" \
                  --data-dir=tikv \
                  --log-file=tikv.log &

#172.29.205.16
$ ./bin/tikv-server --pd="172.29.205.18:2379,172.29.205.13:2379,172.29.205.17:2379" \
                  --addr="172.29.205.16:20160" \
                  --status-addr="172.29.205.16:20180" \
                  --data-dir=tikv \
                  --log-file=tikv.log &
```			  
			
			
### TiDB

TiDB:

```
./bin/tidb-server --store=tikv \
                  --path="172.29.205.18:2379,172.29.205.13:2379,172.29.205.17:2379" \
                  --log-file=tidb.log &
                 
```


### Create Accounts

```
#mysql -h 192.168.199.113 -P 4000 -u root -D test

mysql>

```

修改密码：

```
SET PASSWORD FOR 'root'@'%' = 'Root123##%TOO';

CREATE USER 'bigdata'@'%' IDENTIFIED BY 'BI%Gdata999@!*';
GRANT ALL PRIVILEGES ON *.* TO 'bigdata'@'%' WITH GRANT OPTION;
```

