

## 启动mgmd

登陆`172.29.205.13`, 使用root账号，执行：

```
ndb_mgmd -f /var/lib/mysql-cluster/config.ini
```


## 启动
登陆`172.29.205.14` 、 `172.29.205.15`、`172.29.205.16`, 使用root账号，执行：
执行：

```
ndbd
```

## 启动mysqld

登陆`172.29.205.17` 、 `172.29.205.18`, 使用root账号，执行：
执行：

```
systemctl start mysqld
```

未出现任何信息，这个说明运行成功。


`172.29.205.13`上显示节点消息：


```
ndb_mgm

>show

```