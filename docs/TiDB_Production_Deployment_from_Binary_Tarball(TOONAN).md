
##Delopy TiDB


### PD

192.168.1.147		data-node(tikv)
192.168.1.157		data-node(tikv)
192.168.1.164		data-node(tikv)
192.168.1.22		PD1
192.168.1.221		PD2
192.168.1.229		PD3、TiDB

### Create a database running user account

Log in to the machine using the root user account and create a database running user account (tidb) using the following command:

```
# useradd tidb -m
```

Switch the user from root to tidb by using the following command. You can use this tidb user account to deploy your TiDB cluster.

```
# su - tidb
```

PD1:
```
./pd-server --name=pd1 \
                --data-dir=/home/tidb/pd/pd \
                --client-urls="http://192.168.1.22:2379" \
                --peer-urls="http://192.168.1.22:2380" \
                --initial-cluster="pd1=http://192.168.1.22:2380,pd2=http://192.168.1.221:2380,pd3=http://192.168.1.229:2380" \
                -L "info" \
                --log-file=pd.log &
				
				
./pd-server --name=pd2 \
                --data-dir=/home/tidb/pd/pd \
                --client-urls="http://192.168.1.221:2379" \
                --peer-urls="http://192.168.1.221:2380" \
                --initial-cluster="pd1=http://192.168.1.22:2380,pd2=http://192.168.1.221:2380,pd3=http://192.168.1.229:2380" \
                -L "info" \
                --log-file=pd.log &
				
./pd-server --name=pd3 \
                --data-dir=/home/tidb/pd/pd \
                --client-urls="http://192.168.1.229:2379" \
                --peer-urls="http://192.168.1.229:2380" \
                --initial-cluster="pd1=http://192.168.1.22:2380,pd2=http://192.168.1.221:2380,pd3=http://192.168.1.229:2380" \
                -L "info" \
                --log-file=pd.log &
```



### TiKV

TiKV:

```
#192.168.1.147
$ ./tikv-server --pd="192.168.1.22:2379,192.168.1.221:2379,192.168.1.229:2379" \
                  --addr="192.168.1.147:20160" \
                  --status-addr="192.168.1.147:20180" \
                  --data-dir=tikv \
                  --log-file=tikv.log &

#192.168.1.155
$ ./tikv-server --pd="192.168.1.22:2379,192.168.1.221:2379,192.168.1.229:2379" \
                  --addr="192.168.1.155:20160" \
                  --status-addr="192.168.1.155:20180" \
                  --data-dir=tikv \
                  --log-file=tikv.log &

#192.168.1.164
$ ./tikv-server --pd="192.168.1.22:2379,192.168.1.221:2379,192.168.1.229:2379" \
                  --addr="192.168.1.164:20160" \
                  --status-addr="192.168.1.164:20180" \
                  --data-dir=tikv \
                  --log-file=tikv.log &
```			  
			
			
### TiDB

TiDB:

```
./tidb-server --store=tikv \
                  --path="192.168.1.22:2379,192.168.1.221:2379,192.168.1.229:2379" \
                  --log-file=tidb.log &
                 
```


### Create Accounts

```
#mysql -h 192.168.1.229 -P 4000 -u root -D test

mysql>

```

修改密码：

```
SET PASSWORD FOR 'root'@'%' = 'Root123##%TOO';

CREATE USER 'bigdata'@'%' IDENTIFIED BY 'BI%Gdata999@!*';
GRANT ALL PRIVILEGES ON *.* TO 'bigdata'@'%' WITH GRANT OPTION;
```

