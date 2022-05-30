
### Docker Commands

去mirror镜像：
>docker pull

端口映射：
>-p 49008:8080(母机：docker VM)

目录映射：
>-v $PWD/jenkins:/var/jenkins_data

app名称：
>-t 


查看docker name：

>sudo docker inspect -f='{{.Name}}' $(sudo docker ps -a -q)

查看dockers ip：

>sudo docker inspect -f='{{.NetworkSettings.IPAddress}}' $(sudo docker ps -a -q)

容器增：

docker create [OPTIONS] IMAGE [COMMAND] [ARG...]  #创建一个新的容器但不启动它

docker run [OPTIONS] IMAGE [COMMAND] [ARG...]  #创建一个新的容器并运行一个命令
 

进：

docker exec -ti <container name/id> #不会像attach方式因为退出，导致整个容器退出。          
docker attach  <container name/id> #进入虚拟机，如果从这个stdin中exit，会导致容器的停止。
开启/停止/重启

docker container  start/stop/restart   <hash> 
删：

docker container rm <hash> # 从此机器中移除指定的容器【删除容器时，容器必须是停止状态，否则会报如下错误】
docker container rm $(docker container ls -a -q) # 删除所有容器
docker container kill <hash> # 强制关闭指定的容器
查：

复制代码
docker container ls # 列出所有运行的容器
docker container ls -a # 列出所有的容器

docker ps # 查看我们正在运行的容器
docker ps -l # 查询最后一次创建的容器

docker logs <container id> # 查看容器内的标准输出
docker logs <container name> # 查看容器内的标准输出


docker port <container name/id> <port> # 查看容器端口的映射情况

docker inspect <container id/name> #查看Docker的底层信息,它会返回一个 JSON 文件记录着 Docker 容器的配置和状态信息
复制代码
回到顶部
二：镜像
增：

docker build -t friendlyname . # 使用此目录的“Dockerfile”创建镜像
docker push 192.168.1.52:5000/zabbix #提交镜像到本地私有
docker pull ubuntu:13.10 # 下载ubuntu:13.10镜像
删：

docker image rm <image id> # 从机器中移除指定镜像
docker image rm $(docker image ls -a -q) # 从机器上移除所有镜像
查：

docker image ls -a # 列出机器上所有镜像
docker search httpd # 通过 docker search 命令搜索 httpd 来寻找适合我们的镜像
运:

docker run httpd # 使用镜像仓库
#镜像标签
docker tag <image> username/repository:tag # 标签<image>上传到仓库【repository为镜像源名】
docker push username/repository:tag # 上传标记的图像到注册表
docker run username/repository:tag # 从注册表运行映像

回到顶部
三：其他
docker login # 使用您的Docker凭据登录此CLI会话
docker run -d -p 127.0.0.1:4000:80/udp friendlyname # 后台运行"friendlyname" 镜像并将4000 端口映射到80端口
# -P:将容器内部使用的网络端口映射到我们使用的主机上
# -d:让容器在后台运行
# 默认都是绑定 tcp 端口，如果要绑定 UDP 端口，可以在端口后面加上 /udp
docker run ubuntu:15.10 /bin/echo "Hello world" #Docker 以 ubuntu15.10 镜像创建一个新容器，然后在容器里执行 bin/echo "Hello world"，然后输出结果。
docker run -d -P --name runoob training/webapp python app.py #对容器镜像重新命名


### Define a custom mirror

docker pull centos/py-35-centos


### Committing Changes in a Container to a Docker Image

Then commit the changes to a new Docker image instance using the following command. The -m switch is for the commit message that helps you and others know what changes you made, while -a is used to specify the author. The container ID is the one you noted earlier in the tutorial when you started the interactive docker session. Unless you created additional repositories on Docker Hub, the repository is usually your Docker Hub username:

```
docker commit -m "What did you do to the image" -a "Author Name" container-id repository/new_image_name
```

e.g.

```
docker commit -m "added mariadb-server" -a "Sunday Ogwu-Chinuwa" 59839a1b7de2 finid/centos-mariadb
```


1、删除所有容器
```
docker rm `docker ps -a -q`

#强制关闭和删除容器
docker stop  `docker ps -a -q` && docker rm --force `docker ps -a -q`
```
2、删除所有镜像
```
docker rmi `docker images -q`
```
3、按条件删除镜像
　　没有打标签
```
docker rmi `docker images -q | awk '/^<none>/ { print $3 }'`

```
　　镜像名包含关键字
```
docker rmi --force `docker images | grep doss-api | awk '{print $3}'`    
//其中doss-api为关键字
```

列出退出的容器：
```
docker ps -a -f status=exited

```
删除所有退出的容器
```
docker rm $(docker ps -a -f status=exited -q)
```

## 制作CentOS镜像

mkdir -p ~/build/centos_7.9.2009

``` 
cat >  build/centos_7.9.2009/Dockerfile << EOF
# 指定基础镜像
FROM centos:7.9.2009

# 只在构建镜像的时候执行这些shell指令
RUN yum clean all && \\
    echo "export LC_ALL=zh_CN.UTF-8"  >>  /etc/profile && \\
    echo "export LC_CTYPE=zh_CN.UTF-8"  >>  /etc/profile && \\
    rm -rf /tmp/* && rm -rf /var/cache/yum/* && \\
    localedef -c -f UTF-8 -i zh_CN zh_CN.UTF-8 && \\
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
EOF

```

构建Docker镜像：

``` 
docker build -t  base-centos:7.9.2009 build/centos_7.9.2009/
```

运行：
``` 
docker run \
--name base-centos-7.9.2009 \
--privileged=true \
-dit \
base-centos:7.9.2009 \
/usr/sbin/init
```

进入Docker image中：
``` 
docker exec -it base-centos-7.9.2009 /bin/bash
```

install OpenJDK 17.2

```
yum install -y wget 
wget https://repo.huaweicloud.com/openjdk/17.0.2/openjdk-17.0.2_linux-x64_bin.tar.gz
mv -f openjdk-17.0.2_linux-x64_bin.tar.gz /opt/
cd /opt/
tar xvf openjdk-17.0.2_linux-x64_bin.tar.gz
rm -rf openjdk-17.0.2_linux-x64_bin.tar.gz
cat >>  /etc/profile << EOF
export PATH=/opt/jdk-17.0.2/bin:$PATH
EOF
source /etc/profile
java --version
yum remove -y wget
exit
```

docker commit <容器名称或者ID> <生成的镜像名>:<镜镜像版本号>

发布新的镜像：

```
docker commit base-centos-7.9.2009 base-centos:7.9.2009-openJDK17
```

## 将镜像推送到远程仓库

docker登陆远程仓库。格式为docker login --username=<用户名> <仓库地址>或者docker login -u <用户名> -p <密码> <仓库地址>。

``` 
docker login -u admin -p Harbor12345 harbor-api.service.consul

```

docker tag <镜像ID> <远程镜像仓库地址>:<镜像版本号>标记后。

```
docker tag 66b1bc81e1f2 10.111.38.54/exposure/base-centos:7.9.2009-openJDK17
```

docker push <远程镜像仓库地址>:<镜像版本号>推送到远程仓库。


```` 
docker push 10.111.38.54/exposure/base-centos:7.9.2009-openJDK17
````


e.g.

```

docker tag 625a227574b6 harbor-api.service.consul/public/platform/base-centos:7.9.2009-openJDK17
docker push harbor-api.service.consul/public/platform/base-centos:7.9.2009-openJDK17
```
exposure 项目名称，必须有才能够成功。

## 保存、加载镜像tar包
镜像保存为tar包，格式为docker save -o <文件名> <镜像名>:<镜像标签>。
``` 
docker save -o /data/docker/build/centos_7.9.2009/base-centos_7.9.2009-openJDK17.tar base-centos:7.9.2009-openJDK17
```

加载tar包生成镜像。

``` 
docker load --input base-centos_7.9.2009-openJDK17.tar 

```