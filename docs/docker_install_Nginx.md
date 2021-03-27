
##install nexus3

.拉取镜像
```
# 从 Docker Hub 服务器下载镜像并启动
$ docker run -d -p 80:8080 colovu/nginx:1.18

# 从 Aliyun 服务器下载镜像并启动
$ docker run -d -p 80:8080 registry.cn-shenzhen.aliyuncs.com/colovu/nginx:1.18

```
2.启动容器
```
docker run -d -p 80:8080 -v /opt/nginx/conf:/srv/conf registry.cn-shenzhen.aliyuncs.com/colovu/nginx:1.18
```

镜像默认提供以下数据卷定义，默认数据分别存储在自动生成的应用名对应nginx子目录中：
```
 /srv/data				# 站点源文件
 /srv/conf				# nginx 配置文件
 /var/log					# 日志文件
 /var/run					# 进程运行PID文件

 ```