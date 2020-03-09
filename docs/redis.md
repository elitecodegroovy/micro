

## Redis

测试延时的性能测试，可以使用下面命令：


```
$ redis-benchmark -q -n 10000 -c 1 -d average_size_of_your_objects_in_bytes

```

优化参数设置：


```

tcp-backlog 1024
```

最大客户端连接数量

```
maxclients 30000
```

缓存提出策略：
```
maxmemory-policy volatile-ttl

```

内存设置：

```
maxmemory 21474836480
```

程序以守护进程运行：

```
daemonize yes
```

#表示是否需要密码登陆 注释掉就是不用密码
```
# requirepass toonan123
```

绑定固定的ip地址访问redis服务 如果注释掉就是 不做限制
```
# bind 127.0.0.1
```

redis的链接端口

```
port 6379
```