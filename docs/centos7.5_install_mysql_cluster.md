为了简单，我后面简称mysql-cluster为mc。

1、mc已经包含了mysql，我下载的最新的mc7.5，官方说明包含的是mysql版本是5.7。所以不需要使用别的msyql的安装包安装数据库。同时注意mysql5.7的版本在安装的命令和配置上面和之前的版本有很大的不同，所以网上有很多mc7.5之前的版本，所包含的mysql版本不同，所以安装方法不同。

2、管理节点，mc管理节点负责管理、配置、监控整个集群。

3、数据节点，使用内存存放数据，保存进数据节点的数据都会自动复制并存储到其他数据节点。

4、mysql节点，也叫数据库节点，和我们平时使用的mysql相同，作为数据库使用。被数据节点访问。



注意关闭防火墙firewalld, 设置setenforcing 0

设置系统时钟：
```
timedatectl set-time '2019-01-26 19:44:38'
```

或者安装ntpd，并且启动ntpd服务
```
yum install ntp ntpdate
systemctl enable ntpd
# ntpd clock server
ntpdate -u -s 0.centos.pool.ntp.org 1.centos.pool.ntp.org 2.centos.pool.ntp.org
systemctl status ntpd
systemctl restart ntpd

```
关闭防火墙：
```
systemctl status firewalld
systemctl stop firewalld
systemctl disable firewalld
```
修改文件
```
#需要重启
vi /etc/selinux/config
```

命令语句：
```
#不需要重启，立刻生效
setenforce 0

```

设置

SELINUX=permissive


## 管理节点

运行：

```
./run_node_mgm.sh
```
也就是执行下面脚本：
```
rpm -Uvh perl/*.rpm --nodeps --force
yum -y remove mariadb-libs
rpm -Uvh mysql-cluster-community-common-7.6.9-1.el7.x86_64.rpm 
rpm -Uvh mysql-cluster-community-libs-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-libs-compat-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-client-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-server-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-management-server-7.6.9-1.el7.x86_64.rpm 
```

配置MySQL群集
为配置文件创建新目录。我将使用“/var/lib/mysql-cluster”目录.

```
mkdir -p /var/lib/mysql-cluster
vi /var/lib/mysql-cluster/config.ini

```
内容如下：
```
[ndb_mgmd default]
# Directory for MGM node log files
DataDir=/data1/mysql/mysql-cluster
 
[ndb_mgmd]
#Management Node db1
HostName=10.195.254.71
 
[ndbd default]
NoOfReplicas=2      # Number of replicas
DataMemory=4096M     # Memory allocate for data storage
IndexMemory=2048M    # Memory allocate for index storage
#Directory for Data Node
DataDir=/data1/mysql/data
 
[ndbd]
#Data Node db2
HostName=10.195.254.65
 
[ndbd]
#Data Node db3
HostName=10.195.254.66
 
[ndbd]
#Data Node db4
HostName=10.195.254.67

[mysqld]
#SQL Node db5
HostName=10.195.254.69
 
[mysqld]
#SQL Node db6
HostName=10.195.254.70
```


## 数据节点


在数据节点10.195.254.65上执行下面的一些列操作。

创建数据存储目录：

```
mkdir -p /data1/mysql/data
```

运行脚本：
```
./run_data_node.sh
```
即运行下面脚本：
```
rpm -Uvh perl/*.rpm --nodeps --force
yum -y remove mariadb-libs
rpm -Uvh mysql-cluster-community-common-7.6.9-1.el7.x86_64.rpm 
rpm -Uvh mysql-cluster-community-libs-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-libs-compat-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-client-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-server-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-data-node-7.6.9-1.el7.x86_64.rpm
```

使用vi编辑器在/etc目录中创建一个新配置文件：

```
vi /etc/my.cnf
```

添加的内容：
```
[mysqld]
ndbcluster
ndb-connectstring=10.195.254.71     # IP address of Management Node
 
[mysql_cluster]
ndb-connectstring=10.195.254.71     # IP address of Management Node
```

管理节点配置文件“config.ini”中定义的数据库数据创建新目录。
```
mkdir -p /var/lib/mysql-cluster

```

启动数据节点：
```
ndbd

```
在数据节点10.195.254.66、10.195.254.67上执行上面的操作。

## SQL节点

在10.195.254.69上执行下面操作。

执行脚本：
```
./run_sql_node.sh
```
即运行下面脚本：
```
rpm -Uvh perl/*.rpm --nodeps --force
yum -y remove mariadb-libs
rpm -Uvh mysql-cluster-community-common-7.6.9-1.el7.x86_64.rpm 
rpm -Uvh mysql-cluster-community-devel-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-libs-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-libs-compat-7.6.9-1.el7.x86_64.rpm
rpm -Uvh perl-Class-MethodMaker-2.24-1.el7.rf.x86_64.rpm
rpm -Uvh mysql-cluster-community-client-7.6.9-1.el7.x86_64.rpm
rpm -Uvh mysql-cluster-community-server-7.6.9-1.el7.x86_64.rpm
```

编辑配置文件：
```
vi /etc/my.cnf 
```

添加下面内容：

```
[mysqld]
ndbcluster
ndb-connectstring=10.195.254.71      # IP address for server management node
default_storage_engine=ndbcluster     # Define default Storage Engine used by MySQL
 
[mysql_cluster]
ndb-connectstring=10.195.254.71       # IP address for server management node
```
注意修改`datadir=/var/lib/mysql`的值，指定系统文件数据的存放点。权限是mysql:mysql.
修改指定路径的文件权限归属：
```
chown -R mysql:mysql '文件路径'
```
查看日志：
```
tail -n 100 /var/log/mysqld.log  
```

启动mysql集群。启动顺序为：管理节点→数据节点→SQL节点。

启动的命令上面都有，删去--initial即可

关闭时只需要关闭管理节点，后面的数据节点会同时被关闭，mysql就和原来一样即可

管理节点关闭命令：ndb_mgm -e shutdown

sql服务配置参数：
innodb_buffer_pool_size = 64M


## 验证MySQL

找到默认的root密码：
```
grep 'temporary password' /var/log/mysqld.log

```

现在使用以下命令更改密码：

```
mysql_secure_installation

```

创建用户：

```
mysql -u root -p

```
```
CREATE USER 'gc'@'%' IDENTIFIED BY '019#^(963';

GRANT ALL PRIVILEGES ON *.* TO 'gc'@'%' identified by '019#^(963';

CREATE USER 'root'@'%' IDENTIFIED BY '2019027Pig!88';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '2019027Pig!88' WITH GRANT OPTION;
```

创建数据库：

```
create database grafana default character set utf8mb4 collate utf8mb4_unicode_ci;
CREATE USER 'grafana'@'%' IDENTIFIED BY 'Toonan2019@!*';

GRANT ALL PRIVILEGES ON grafana.* TO 'grafana'@'%' identified by 'Toonan2019@!*';

flush privileges;
```
GRANT select  ON *.* TO 'grafana'@'%' identified by 'Toonan2019@!*';

create database `offline-oa` default character set utf8mb4 collate utf8mb4_unicode_ci;
CREATE USER 'gc'@'%' IDENTIFIED BY 'gc123456#G';

GRANT ALL PRIVILEGES ON *.* TO 'gc'@'%' identified by 'gc123456#G';

flush privileges;