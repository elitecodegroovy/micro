
##install MySQL8.0.24

```

docker pull mysql/mysql-server:8.0.24-1.2.2-server 
```

数据和配置文件目录
``` 
mkdir -p /opt/mysql/conf
mkdir -p /opt/app/mysql/data
mkdir -p /opt/app/mysql/logs
```

编辑 
``` 
vi /opt/mysql/conf/my.cnf
```

内容：

``` 

# Remove leading # and set to the amount of RAM for the most important data
# cache in MySQL. Start at 70% of total RAM for dedicated server, else 10%.
# innodb_buffer_pool_size = 128M
#
# Remove leading # to turn on a very important data integrity option: logging
# changes to the binary log between backups.
# log_bin
#
# Remove leading # to set options mainly useful for reporting servers.
# The server defaults are faster for transactions and fast SELECTs.
# Adjust sizes as needed, experiment to find the optimal values.
# join_buffer_size = 128M
# sort_buffer_size = 2M
# read_rnd_buffer_size = 2M

# Remove leading # to revert to previous value for default_authentication_plugin,
# this will increase compatibility with older clients. For background, see:
# https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_default_authentication_plugin
# default-authentication-plugin=mysql_native_password
# For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/8.0/en/server-configuration-defaults.html

[mysqld]
lower_case_table_names=1
skip-host-cache
skip-name-resolve
datadir=/var/lib/mysql
socket=/var/lib/mysql/mysql.sock
secure-file-priv=/var/lib/mysql-files
user=mysql

default-time-zone ='+8:00'
default_authentication_plugin=mysql_native_password
tmpdir                          =/tmp
sql_mode = STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION


wait_timeout=1814400
#
# Remove leading # and set to the amount of RAM for the most important data
# cache in MySQL. Start at 70% of total RAM for dedicated server, else 10%.
#################innodb########################
#innodb_buffer_pool_size = 90000M
#
# Remove the leading "# " to disable binary logging
# Binary logging captures changes between backups and is enabled by
# default. It's default setting is log_bin=binlog
disable_log_bin

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

#datadir=/usr/local/mysql/data

log-error=/var/log/mysqld.log
#slow_query_log_file             = /usr/local/mysql/log/slow.log
#slow_query_log                  = 1
#long_query_time                 = 10
#log_slow_admin_statements       = ON
pid-file=/var/run/mysqld/mysqld.pid
#TIMESTAMP如果没有显示声明NOT NULL，允许NULL值
explicit_defaults_for_timestamp = true

max_connections=2000
max_connect_errors=10

[client]
port=3306
socket=/var/lib/mysql/mysql.sock

```

执行命令：
``` 
	docker run \
    --network=host \
    -e MYSQL_ROOT_PASSWORD=1234567890 \
    -v /opt/mysql/data:/var/lib/mysql:rw \
    -v /opt/mysql/log:/var/log/mysql:rw \
	-v /opt/mysql/conf/my.cnf:/etc/my.cnf:rw \
    -v /etc/localtime:/etc/localtime:ro \
   --name mysql8 \
   --restart=always \
    -d mysql/mysql-server:8.0.24-1.2.2-server 
```

参数说明
``` 
--restart=always                                            -> 开机启动容器,容器异常自动重启
-d                                                          -> 以守护进程的方式启动容器
-v /opt/app/mysql/conf.d/my.cnf:/etc/my.cnf                 -> 映射配置文件，此版本映射，其他版本可能不一样
-v /opt/app/mysql/logs://var/log/mysql                      -> 映射日志
-v /opt/app/mysql/data:/var/lib/mysql                       -> 映射数据
--network=host                                              -> 本机和docker网络端口一一映射
--name mysql                                                -> 指定容器名称
-e MYSQL_ROOT_PASSWORD=123456                               -> 写入配置root密码
```

查看运行状态：
``` 

[root@host-172-16-27-32 ~]# docker ps -a
CONTAINER ID        IMAGE                                    COMMAND                  CREATED             STATUS                    PORTS                                            NAMES
c00b6e6539bd        mysql/mysql-server:8.0.24-1.2.2-server   "/entrypoint.sh mysq…"   23 minutes ago      Up 23 minutes (healthy)      

```

添加root@% 账号

``` 
docker exec -it mysql /bin/bash

>mysql -uroot -p  ## 输入root密码

>use mysql

> select user , host from user;   ## 查看账号信息

>CREATE USER 'root'@'%' IDENTIFIED BY 'Root123##%2019';

>GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;

> flush privileges;

>exit;

```

查看是否区分表的大小写：

``` 
show variables like 'lower%'

```

修改root@%的密码

``` 
 ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '新密码';
```