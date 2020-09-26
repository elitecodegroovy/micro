
# Install PostgreSQL

## 前置条件

root下操作:

``` 
# 停止postgresql,此处为注册了服务的停止方式，12版本
systemctl stop postgresql-12
# 卸载 包含postgresql名的程序包
yum remove postgresql*
# 删除postgres用户及其对应的用户目录（/home/postgres）
userdel -r postgres
# 检查是否存在 /usr/psqlXXX文件夹，有的话删除，例如
rm -rf /usr/pgsql-12

```

创建用户及用户组:

``` 
# 创建postgres用户组
groupadd postgres
# 创建postgres用户，用户位于postgres组内
useradd -g postgres postgres
# 为postgres用户设置密码
passwd postgres
```

安装部分依赖:
``` 
yum install gcc readline-devel.x86_64 zlib-devel.x86_64
```

## 安装PostgreSQL12

``` 
# 切换用户
su postgres
# 回到用户目录
cd ~
```

下载postgres12源码包(/home/postgres/postgresql-12.1.tar.gz):
``` 
wget -c https://ftp.postgresql.org/pub/source/v12.3/postgresql-12.3.tar.gz

tar -zxvf postgresql-12.3.tar.gz
mkdir postgresql
cd postgresql-12.3
./configure --prefix=/home/postgres/postgresql
make && make install

```

配置环境变量:
``` 
vi  ~/.bash_profile

export PGHOME=/home/postgres/postgresql
export PATH=$PATH:$HOME/.local/bin:$HOME/bin:$PGHOME/bin

# 输入： wq!  Enter(回车)
```

重载.bash_profile
```  
source  ~/.bash_profile
```

初始化postgreSQL (/home/postgres/):

``` 
# /home/postgres/PGDATA为postgres数据文件的存放路径
initdb -D /home/postgres/PGDATA
```
启动:
``` 
# /home/postgres/PGDATA为数据文件存放路径，里面也有一些配置文件
pg_ctl -D /home/postgres/PGDATA -l logfile start

# 停止
pg_ctl -D /home/postgres/PGDATA -l logfile stop

# 重启
pg_ctl -D /home/postgres/PGDATA -l logfile restart

# 查看状态
pg_ctl -D /home/postgres/PGDATA -l logfile status

```

## 开机启动

切换到root用户，　执行命令：

```  
cp /home/postgres/postgresql-12.3/contrib/start-scripts/linux /etc/init.d/postgresql

```
修改脚本中的参数项 `vi /etc/init.d/postgresql:
```
#pg主目录
prefix=/home/postgres/postgresql
#Pg data目录
PGDATA="/home/postgres/PGDATA"
```
可执行权限:

``` 
chmod a+x /etc/init.d/postgresql
```


``` 
chkconfig --add postgresql
chkconfig postgresql on
```

start|stop|reload service:

```  
启动服务：service postgresql start
停止服务：service postgresql stop
重启服务：service postgresql reload
```

## 修改配置PostgreSQL
配置允许访问的IP及postgres用户，允许任意IP通过密码访问postgres的任意用户。

``` 
vi /home/postgres/PGDATA/pg_hba.conf

#add setting
host    all     all     0.0.0.0/0   md5
```

配置监听的IP范围等,监听任意IP。

``` 
vi　/home/postgres/PGDATA/postgresql.conf
```

修改
``` 
#
listen_addresses = '*'
```

## 修改账户信息

切换postgres，执行：
```
psql -Upostgres -Ppostgres

postgres=# ALTER USER postgres PASSWORD 'postgres123';
```