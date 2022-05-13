
## set selinux

``` 
[root@localhost ~]# setenforce 0
setenforce: SELinux is disabled
[root@localhost ~]# vi /etc/selinux/config

```

## Master Machine

vi /etc/my.cnf

``` 
bind-address =173.82.120.14
server-id =1
log_bin =mysql-bin
```

Finally, restart the MySQL service.

```
sudo systemctl restart mysqld
```

Create an account that do the sync data operation.

``` 
CREATE USER 'replica'@'172.16.30.44' IDENTIFIED BY 'replica6!@#$%^';
GRANT REPLICATION SLAVE ON *.*TO 'replica'@'172.16.30.44';
```

Show master status:

``` 
mysql> SHOW MASTER STATUS
    -> ;
+---------------+----------+--------------+------------------+------------------                                               -+
| File          | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set                                                |
+---------------+----------+--------------+------------------+------------------                                               -+
| binlog.000001 | 91352121 |              |                  |                                                                  |
+---------------+----------+--------------+------------------+------------------                                               -+
1 row in set (0.00 sec)
```


## Slave Machine


vi /etc/my.cnf 

``` 
bind-address=173.82.115.165
server-id=2
log_bin=mysql-bin
```


login MySQL Server and stop slave thread.

``` 
>./mysql -u root -p


mysql> STOP SLAVE;


mysql>CHANGE MASTER TO MASTER_HOST='172.16.30.43' ,MASTER_PORT=3309,  MASTER_USER='replica' , MASTER_PASSWORD='replica6!@#$%^' , MASTER_LOG_FILE='binlog.000001' ,MASTER_LOG_POS=91353166;

Query OK, 0 rows affected, 8 warnings (0.04 sec)
mysql>start slave;
Query OK, 0 rows affected, 1 warning (0.01 sec)


```

Test the master-slave MySQL8.

master machine:
``` 
CREATE DATABASE IF NOT EXISTS `baseplatform` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
```

Query all databases, then you can see the database baseplatform has been synced to the slave database;.

``` 
show databases;
```





