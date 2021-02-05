
```
 yum localinstall https://dev.mysql.com/get/mysql57-community-release-el7-11.noarch.rpm
```

Copy Install MySQL 5.7 package

Install MySQL as any other package using yum:

```
yum install -y mysql-community-server
```

Once the installation is completed, to enable and start the MySQL servic type:
```
sudo systemctl enable mysqld
sudo systemctl start mysqld

```

root passwd:

```
grep 'temporary password' /var/log/mysqld.log
```

change root passwd:

```
sudo mysql_secure_installation
```

authenticate user:


```
mysql -u root -p passwd

#create user
#CREATE USER 'toonan'@'%' IDENTIFIED BY "toonan123#T";
#grant all privileges  on *.* to toonan@'%' identified by "toonan123#T";
grant all privileges  on *.* to root@'%' identified by "password";
flush privileges;
```

root user privileges: 
``` 
GRANT ALL PRIVILEGES ON *.* TO 'myuser'@'192.168.1.3' IDENTIFIED BY 'mypassword' WITH GRANT OPTION;  
FLUSH   PRIVILEGES;
```

### 乱码

characterEncoding = utf8


### 备份

全局备份：
```
mysqldump –h 主机名 –u 用户名 –p  - -all-databases  > 备份文件名.sql

```

调用mysqldump带有- -databases选项备份指定的数据库:
```
mysqldump –u 用户名 –p  --databases db1 db2 db3 …  > 备份文件名.sql


mysqldump -alv -h 172.29.205.11 -u  backup -pbackupaaaaafsdasdafsw123#T --databases nsmjpt2 > nsmjpt2_db.sql
```

```
mysqldump -alv -h 172.29.205.17 -u  backup -pbackupaaaaafsdasdafsw123#T --databases nsmjpt2 --ignore-table=nsmjpt2.table1 --ignore-table=nsmjpt2.table2 > nsmjpt2_db.sql
```
调用mysqldump备份某个数据库中的某几张表：

```
mysqldump –u用户名 –p 数据库名 表名1 表名2 表名3… > 备份文件名.sql
```


过mysqldump备份的文件，如果用了- -all-databases或- -databases选项，则在备份文件中包含CREATE DATABASE和USE语句，故并不需要指定一个数据库名去恢复备份文件。

在Shell命令下：:
```
shell>  mysql –u 用户名 –p  < 备份文件.sql
```

果通过mysqldump备份的是单个数据库，且没有使用- -databases选项，则备份文件中不包含CREATE DATABASE和USE语句，那么在恢复的时候必须先创建数据库。

在shell命令下：
```

shell>  mysqladmin –u 用户名 –p create 数据库名     //创建数据库

shell>  mysql –u 用户名 –p 数据库名 < 备份文件.sql

```

在mysql命令下：

```

mysql>  CREATE DATABASE IF NOT EXIST 数据库名;

mysql>  USE 数据库名;

mysql>  source 备份文件.sql;
```

MySQL Opts:
```

GRANT ALL PRIVILEGES ON *.* TO 'gc'@'%' IDENTIFIED BY 'gc123456#G' WITH GRANT OPTION; 
FLUSH   PRIVILEGES ;

```

问题：
ERROR! The server quit without updating PID file (/data/mysql/data/mysq.pid).

权限问题，解决方案：
``` 
[root@localhost ~]# chown -R mysql:mysql /data/mysql/data
[root@localhost ~]# chmod -R 755 /data/mysql/data
```





mysqldump –h 192.168.1.229 -u gc –p  --all-databases  > all-5.7.sql