
## pull image

```
docker pull redis:latest
```

## run container

```

docker run --name alone2-redis -p 8379:6379 -d --restart=always redis:latest redis-server --appendonly yes --requirepass "TOONAN"

```

-p 6379:6379 :将容器内端口映射到宿主机端口(右边映射到左边) 
redis-server –appendonly yes : 在容器执行redis-server启动命令，并打开redis持久化配置 
requirepass “your passwd” :设置认证密码 
–restart=always : 随docker启动而启动