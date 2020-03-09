

##卸载MySQL

查看是否还有mysql软件：

```
>rpm -qa|grep mysql
mysql-cluster-community-libs-7.6.9-1.el7.x86_64
mysql-cluster-community-devel-7.6.9-1.el7.x86_64
mysql-cluster-community-common-7.6.9-1.el7.x86_64
mysql-cluster-community-client-7.6.9-1.el7.x86_64
mysql-cluster-community-libs-compat-7.6.9-1.el7.x86_64
mysql-cluster-community-server-7.6.9-1.el7.x86_64

```

卸载程序：

```
yum remove mysql
```

卸载依赖包：

```
>yum remove mysql-cluster-community-libs
>yum remove mysql-cluster-community-common
```


查看系统中是否以rpm包安装的mysql：

 
```
> rpm -qa | grep -i mysql
```

删除mysql服务
```
[root@localhost local]# chkconfig --list | grep -i mysql
[root@localhost local]# chkconfig --del mysql
```

查找mysql目录：

```
find / -name mysql
/home/data/mysqlcluster/data/mysql
/etc/selinux/targeted/active/modules/100/mysql
/etc/selinux/targeted/tmp/modules/100/mysql
/var/lib/mysql
/usr/share/mysql
```
删除mysql目录：

```
rm -rf /home/data/mysqlcluster/data/mysql
rm -rf /etc/selinux/targeted/active/modules/100/mysql
rm -rf /var/lib/mysql
rm -rf /usr/share/mysql
rm -rf /usr/my.cnf
```


