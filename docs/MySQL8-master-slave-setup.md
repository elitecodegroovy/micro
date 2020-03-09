
## 安装MySQL

假设已经成功安装了MySQL服务器，IP信息如下：

| Host Name | IP            | Services   | Data Path        |
| --------- | ------------- | ---------- | ------------     |
| **host1** | 10.195.254.69 | MASTER     | /data1/mysql     |
| **host2** | 10.195.254.70 | SLAVE      | /data1/mysql     |

## Master机器设置

```
vi /etc/my.cnf
```
修改内容后，显示效果如下：
```
# For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/8.0/en/server-configuration-defaults.html

[mysqld]
server-id = 1
default-time-zone ='+8:00'
#lower_case_table_names=1
default_authentication_plugin=mysql_native_password
tmpdir=/data1/mysql/tmp/
sql_mode = STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION

#等待超时设置24天
wait_timeout=1814400
#
# Remove leading # and set to the amount of RAM for the most important data
# cache in MySQL. Start at 70% of total RAM for dedicated server, else 10%.
#################innodb########################
innodb_buffer_pool_size = 90000M
#
# Remove the leading "# " to disable binary logging
# Binary logging captures changes between backups and is enabled by
# default. It's default setting is log_bin=binlog
# disable_log_bin
log-bin = /data1/mysql/mysql-bin.log
binlog_format = ROW
max_binlog_size = 600M
sync_binlog = 1
expire-logs-days = 7
#
# Remove leading # to set options mainly useful for reporting servers.
# The server defaults are faster for transactions and fast SELECTs.
# Adjust sizes as needed, experiment to find the optimal values.
join_buffer_size = 256M
read_buffer_size = 16M
sort_buffer_size = 64M
read_rnd_buffer_size = 64M

#
# Remove leading # to revert to previous value for default_authentication_plugin,
# this will increase compatibility with older clients. For background, see:
# https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_default_authentication_plugin
# default-authentication-plugin=mysql_native_password

datadir=/data1/mysql/data
socket=/var/lib/mysql/mysql.sock

log-error=/var/log/mysqld.log
slow_query_log_file             = /data1/mysql/log/slow.log
slow_query_log                  = 1
long_query_time                 = 10
log_slow_admin_statements       = ON
pid-file=/var/run/mysqld/mysqld.pid
#TIMESTAMP如果没有显示声明NOT NULL，允许NULL值
explicit_defaults_for_timestamp = true

# 允许最大连接数
max_connections=10000
# 允许连接失败的次数。
max_connect_errors=10


```
重启服务：

```
systemctl restart mysqld

```
登陆mysql 服务，显示BIN_LOG信息：

```
>mysql -u root -p


mysql>show master status \G;
*************************** 1. row ***************************
             File: mysql-bin.000001
         Position: 1532
     Binlog_Do_DB: 
 Binlog_Ignore_DB: 
Executed_Gtid_Set: 
1 row in set (0.00 sec)

ERROR: 
No query specified
```

为slave设置特权访问账号：

```
 #master配置 
mysql>create user replication@172.29.205.17 identified by '10.T195.2s54.7oo0';

mysql>grant replication slave on *.* to replication@172.29.205.17;

mysql> flush privileges;
 #Confirm grants for created user:
mysql> show grants for replication@172.29.205.17 ;
+-----------------------------------------------------------------+
| Grants for replication@10.195.254.70                            |
+-----------------------------------------------------------------+
| GRANT REPLICATION SLAVE ON *.* TO `replication`@`10.195.254.70` |
+-----------------------------------------------------------------+
1 row in set (0.00 sec)
```

## SLAVE设置
编辑文件：
```
vi /etc/my.cnf
```
编辑后，显示效果如下：

```
# For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/8.0/en/server-configuration-defaults.html

[mysqld]
server-id = 2
read_only = 1
default-time-zone ='+8:00'
#lower_case_table_names=1
default_authentication_plugin=mysql_native_password
tmpdir                          =/data1/mysql/tmp/
sql_mode = STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION

#等待超时设置24天
wait_timeout=1814400
#
# Remove leading # and set to the amount of RAM for the most important data
# cache in MySQL. Start at 70% of total RAM for dedicated server, else 10%.
#################innodb########################
innodb_buffer_pool_size = 90000M
#
# Remove the leading "# " to disable binary logging
# Binary logging captures changes between backups and is enabled by
# default. It's default setting is log_bin=binlog
#disable_log_bin
log_bin = /data1/mysql/mysql-bin.log
binlog_format = ROW
max_binlog_size = 600M
sync_binlog = 1
expire-logs-days = 7
#
# Remove leading # to set options mainly useful for reporting servers.
# The server defaults are faster for transactions and fast SELECTs.
# Adjust sizes as needed, experiment to find the optimal values.
join_buffer_size = 256M
read_buffer_size = 16M
sort_buffer_size = 64M
read_rnd_buffer_size = 64M

#
# Remove leading # to revert to previous value for default_authentication_plugin,
# this will increase compatibility with older clients. For background, see:
# https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_default_authentication_plugin
# default-authentication-plugin=mysql_native_password

datadir=/data1/mysql/data
socket=/var/lib/mysql/mysql.sock

log-error=/var/log/mysqld.log
slow_query_log_file             = /data1/mysql/log/slow.log
slow_query_log                  = 1
long_query_time                 = 10
log_slow_admin_statements       = ON
pid-file=/var/run/mysqld/mysqld.pid
#TIMESTAMP如果没有显示声明NOT NULL，允许NULL值
explicit_defaults_for_timestamp = true

# 允许最大连接数
max_connections=10000
# 允许连接失败的次数。
max_connect_errors=10

```

重启mysql服务:

```
systemctl restart mysqld
```

设置全局变量：
```
>mysql -u root -p

mysql> CHANGE MASTER TO MASTER_HOST='10.195.254.69', MASTER_USER='replication',MASTER_PASSWORD='10.T195.2s54.7oo0',MASTER_LOG_FILE='mysql-bin.000001',MASTER_LOG_POS=1532;
```
> 上面的参数根据master机器的binlog信息得到。比如： MASTER_LOG_FILE、MASTER_LOG_POS（master机器上显示binlog信息：`show master status \G;`）.
启动slave：

```
mysql>start slave;
```


## 验证

SLAVE机器上，执行下面命令：
```
mysql> show slave status\G
*************************** 1. row ***************************
               Slave_IO_State: Waiting for master to send event
                  Master_Host: 10.195.254.69
                  Master_User: replication
                  Master_Port: 3306
                Connect_Retry: 60
              Master_Log_File: mysql-bin.000001
          Read_Master_Log_Pos: 1532
               Relay_Log_File: dbsvr-1-relay-bin.000002
                Relay_Log_Pos: 991
        Relay_Master_Log_File: mysql-bin.000001
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
              Replicate_Do_DB: 
          Replicate_Ignore_DB: 
           Replicate_Do_Table: 
       Replicate_Ignore_Table: 
      Replicate_Wild_Do_Table: 
  Replicate_Wild_Ignore_Table: 
                   Last_Errno: 0
                   Last_Error: 
                 Skip_Counter: 0
          Exec_Master_Log_Pos: 1532
              Relay_Log_Space: 1201
              Until_Condition: None
               Until_Log_File: 
                Until_Log_Pos: 0
           Master_SSL_Allowed: No
           Master_SSL_CA_File: 
           Master_SSL_CA_Path: 
              Master_SSL_Cert: 
            Master_SSL_Cipher: 
               Master_SSL_Key: 
        Seconds_Behind_Master: 0
Master_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error: 
               Last_SQL_Errno: 0
               Last_SQL_Error: 
  Replicate_Ignore_Server_Ids: 
             Master_Server_Id: 1
                  Master_UUID: a1da6f11-9fbe-11e9-809b-286ed45ac635
             Master_Info_File: mysql.slave_master_info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Slave_SQL_Running_State: Slave has read all relay log; waiting for more updates
           Master_Retry_Count: 86400
                  Master_Bind: 
      Last_IO_Error_Timestamp: 
     Last_SQL_Error_Timestamp: 
               Master_SSL_Crl: 
           Master_SSL_Crlpath: 
           Retrieved_Gtid_Set: 
            Executed_Gtid_Set: 
                Auto_Position: 0
         Replicate_Rewrite_DB: 
                 Channel_Name: 
           Master_TLS_Version: 
       Master_public_key_path: 
        Get_master_public_key: 0
            Network_Namespace: 
1 row in set (0.00 sec)
```

有下面信息，标识SLAVE配置成功：

```
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
```