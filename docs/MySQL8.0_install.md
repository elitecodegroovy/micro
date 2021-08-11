## MySQL Server安装


```
rpm -Uvh mysql-community-common-8.0.16-2.el7.x86_64.rpm 
rpm -Uvh mysql-community-libs-8.0.16-2.el7.x86_64.rpm
rpm -Uvh mysql-community-libs-compat-8.0.16-2.el7.x86_64.rpm
rpm -Uvh mysql-community-devel-8.0.16-2.el7.x86_64.rpm
rpm -Uvh mysql-community-client-8.0.16-2.el7.x86_64.rpm
rpm -Uvh mysql-community-server-8.0.16-2.el7.x86_64.rpm

```

```shell script
groupadd mysql 
useradd -r -g mysql -s /bin/false mysql
```



## my.cnf配置

vi /etc/my.cnf

```

# For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/8.0/en/server-configuration-defaults.html

[mysqld]
server-id = 1
default-time-zone ='+8:00'
#lower_case_table_names=1
default_authentication_plugin=mysql_native_password
tmpdir                          =/home/mysql/tmp/
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
disable_log_bin
#log-bin = /data1/mysql/mysql-bin.log
#binlog_format = ROW
#max_binlog_size = 600M
#sync_binlog = 1
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

datadir=/home/mysql/data
socket=/var/lib/mysql/mysql.sock

log-error=/var/log/mysqld.log
slow_query_log_file             = /home/mysql/log/slow.log
slow_query_log                  = 1
long_query_time                 = 10
log_slow_admin_statements       = ON
pid-file=/var/run/mysqld/mysqld.pid
#TIMESTAMP如果没有显示声明NOT NULL，允许NULL值
explicit_defaults_for_timestamp = true

# 允许最大连接数
max_connections=2000
# 允许连接失败的次数。
max_connect_errors=10
```

启动mysql

```
systemctl start mysqld
```

开启启动：

```
systemctl enable mysqld
```

修改root密码：

```
>grep 'temporary password' /var/log/mysqld.log
>
#login mysql
>mysql -u root -p
>
```

```
ALTER USER 'root'@'localhost' IDENTIFIED BY 'TOONAN2019)#@*t1236';

```

ALTER USER 'root'@'localhost' IDENTIFIED BY 'Tisson2021)#@*t1236';
添加账号：

```
CREATE USER 'gc'@'%' IDENTIFIED BY '019#^(963adhCswGC';

GRANT ALL PRIVILEGES ON nsmjpt2.* TO 'gc'@'%' ;

##兼容旧客户端
ALTER USER 'gc'@'%' IDENTIFIED WITH mysql_native_password BY '019#^(963adhCsw';
##ALTER USER 'root'@'%' IDENTIFIED BY 'password' PASSWORD EXPIRE NEVER; #修改加密规则 
##ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'password'; #更新一下用户的密码 
FLUSH PRIVILEGES; #刷新权限 

ALTER USER 'root'@'%' IDENTIFIED BY 'Tisson2021)Remote#' PASSWORD EXPIRE NEVER;

CREATE USER 'root'@'%' IDENTIFIED BY '019#^(963adhCsw';

create database grafana default character set utf8mb4 collate utf8mb4_unicode_ci;
create database nsmjpt2 default character set utf8mb4 collate utf8mb4_unicode_ci;
CREATE USER 'gc'@'%' IDENTIFIED BY '019#^(963adhCsw';
GRANT ALL PRIVILEGES ON nsmjpt2.* TO 'gc'@'%' WITH GRANT OPTION;

create database nacos default character set utf8mb4 collate utf8mb4_unicode_ci;

GRANT ALL PRIVILEGES ON baseplatform.* TO 'gc'@'%' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON bpworkflow.* TO 'gc'@'%' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON nacos.* TO 'gc'@'%' WITH GRANT OPTION;
```

## 安装插件

-- 插件动态安装启用
```
mysql> INSTALL PLUGIN CONNECTION_CONTROL SONAME 'connection_control.so';
mysql> INSTALL PLUGIN CONNECTION_CONTROL_FAILED_LOGIN_ATTEMPTS SONAME 'connection_control.so';
```
验证是否正常安装
```
mysql> SELECT PLUGIN_NAME, PLUGIN_STATUS 
FROM INFORMATION_SCHEMA.PLUGINS
WHERE PLUGIN_NAME LIKE 'connection%';

mysql> SHOW PLUGINS;

```


- 查看默认相关变量
```
mysql> show variables like 'connection_control%';
+-------------------------------------------------+------------+
| Variable_name                                   | Value      |
+-------------------------------------------------+------------+
| connection_control_failed_connections_threshold | 3          |
| connection_control_max_connection_delay         | 2147483647 |
| connection_control_min_connection_delay         | 1000       |
+-------------------------------------------------+------------+
```
-- 定制化配置

```
mysql> SET GLOBAL connection_control_failed_connections_threshold = 30;
mysql> SET GLOBAL connection_control_min_connection_delay = 600000;
mysql> SET GLOBAL connection_control_max_connection_delay = 1200000;
```

-- 查看修改后的配置
```
mysql> show variables like 'connection_control%';

```