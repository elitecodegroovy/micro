## MySQL Server安装

Download the source code and tar the source 

```shell script
groupadd mysql 
useradd -r -g mysql -s /bin/false mysql

passwd mysql
# set mysql passwd

#verify login
su - mysql

#it succeeded if you logged in.
```

### Install GCC 9.3.0

Please look the file gcc_setup.md.


### Install cmake 3.17.0

```
./bootstrap
make
make install
```


### Download Boost

Download the boost source code and tar to the specified path. The build command will like this.

```shell script
cmake . -DWITH_BOOST=/usr/local/boost_version_number
```

### Install Bison And  Ncurses 

Install the git , bison and ncurses.


### Building

```shell script
cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local/mysql \
-DWITH_DEBUG=OFF \
-DDEFAULT_CHARSET=utf8 \
-DDEFAULT_COLLATION=utf8_general_ci \
-DENABLED_LOCAL_INFILE=ON \
-DWITH_INNODB_MEMCACHED=ON \
-DWITH_SSL=system \
-DWITH_INNOBASE_STORAGE_ENGINE=1 \
-DWITH_FEDERATED_STORAGE_ENGINE=1 \
-DWITH_BLACKHOLE_STORAGE_ENGINE=1 \
-DWITH_ARCHIVE_STORAGE_ENGINE=1 \
-DWITHOUT_EXAMPLE_STORAGE_ENGINE=1 \
-DWITH_PERFSCHEMA_STORAGE_ENGINE=1 \
-DCOMPILATION_COMMENT="zsd edition" \
-DCMAKE_CXX_COMPILER=/usr/bin/g++ \
-DWITH_BOOST=/opt/mysql-8.0.28/boost \
-DMYSQL_UNIX_ADDR=/usr/local/mysql/mysql.sock \
-DSYSTEMD_PID_DIR=/usr/local/mysql \
-DSYSCONFDIR=/usr/local/mysql/    \
-DTMPDIR=/usr/local/mysql/tmp/

make -j$(nproc) 
make install
```

`-DSYSCONFDIR=dir_name ` The default my.cnf option file directory.

`-DSYSTEMD_PID_DIR=dir_name`  The name of the directory in which to 
create the PID file when MySQL is managed by systemd. 
The default is /var/run/mysqld; this might be changed 
implicitly according to the INSTALL_LAYOUT value.

`-DTMPDIR=dir_name` The default location to use for the tmpdir system variable.
 If unspecified, the value defaults to P_tmpdir in <stdio.h>.
###  my.cnf配置

vi /etc/my.cnf

```
# For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/8.0/en/server-configuration-defaults.html

[mysqld]
server-id = 1
default-time-zone ='+8:00'
#lower_case_table_names=1
default_authentication_plugin=mysql_native_password
tmpdir                          =/usr/local/mysql/tmp/
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
#log-bin = /usr/local/mysql/log/mysql-bin.log
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

datadir=/usr/local/mysql/data
socket=/usr/local/mysql/mysql.sock

log-error=/var/log/mysqld.log
slow_query_log_file             = /usr/local/mysql/log/slow.log
slow_query_log                  = 1
long_query_time                 = 10
log_slow_admin_statements       = ON
pid-file=/usr/local/mysql/mysqld.pid
#TIMESTAMP如果没有显示声明NOT NULL，允许NULL值
explicit_defaults_for_timestamp = true

# 允许最大连接数
max_connections=2000
# 允许连接失败的次数。
max_connect_errors=10
```

```shell script
cd /usr/local
mkdir -p /usr/local/mysql/data
mkdir -p /usr/local/mysql/log
mkdir -p /usr/local/mysql/tmp
chown -R mysql:mysql mysql
bin/mysql_ssl_rsa_setup

# Next command is optional
cp support-files/mysql.server /etc/init.d/mysqld


bin/mysqld --initialize --user=mysql --datadir=/usr/local/mysql/data
```

启动mysql

```
service mysqld start
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
ALTER USER 'root'@'localhost' IDENTIFIED BY 'TOONAN2019)#@*t123';

```
添加账号：

```
CREATE USER 'gc'@'%' IDENTIFIED BY '019#^(963adhCsw';

GRANT ALL PRIVILEGES ON nsmjpt2.* TO 'gc'@'%' ;

##兼容旧客户端
ALTER USER 'gc'@'%' IDENTIFIED WITH mysql_native_password BY '019#^(963adhCsw';
##ALTER USER 'root'@'%' IDENTIFIED BY 'password' PASSWORD EXPIRE NEVER; #修改加密规则 
##ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'password'; #更新一下用户的密码 
FLUSH PRIVILEGES; #刷新权限 


create database grafana default character set utf8mb4 collate utf8mb4_unicode_ci;
create database nsmjpt2 default character set utf8mb4 collate utf8mb4_unicode_ci;
CREATE USER 'gc'@'%' IDENTIFIED BY '019#^(963adhCsw';
GRANT ALL PRIVILEGES ON nsmjpt2.* TO 'gc'@'%' WITH GRANT OPTION;

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