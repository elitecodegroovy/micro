## remove old docker

```
sudo yum remove docker  docker-common docker-selinux docker-engine
```

## Docker Install

```
yum install -y http://mirror.centos.org/centos/7/extras/x86_64/Packages/container-selinux-2.42-1.gitad8f0f7.el7.noarch.rpm
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-18.03.0.ce-1.el7.centos.x86_64.rpm
```

最新：
``` 
yum install -y http://mirror.centos.org/centos/7/extras/x86_64/Packages/container-selinux-2.119.1-1.c57a6f9.el7.noarch.rpm
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-20.10.9-3.el7.x86_64.rpm  
```

文件： docker-ce-20.10.9-3.el7.x86_64.rpm  
Install Docker CE, changing the path below to the path where you downloaded the Docker package.
```
$ sudo yum install docker-ce-18.03.0.ce-1.el7.centos.x86_64.rpm

```

Start Docker. 

```
$ sudo systemctl start docker
Verify that docker is installed correctly by running the hello-world image.

$ sudo docker run hello-world

``` --

为了加快 pull image 的速度，可以使用国内的仓库镜像服务器，同时增加下载的并发数。(如果 dockerd 已经运行，则需要重启 dockerd 生效。)

```
$ cat > /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": ["https://docker.mirrors.ustc.edu.cn", "https://hub-mirror.c.163.com"],
  "max-concurrent-downloads": 10
}
EOF
```
e.g.安装jenkins

```
docker pull jenkins
mkdir -p /opt/jenkins_home
chown -R 1000 /opt/jenkins_home
docker run --name jenkinsd \
-u root \
-p 49001:8080 -p 50000:50000 \
-v /opt/jenkins_home:/var/jenkins_home \
-v /var/run/docker.sock:/var/run/docker.sock \
-t jenkins 
```

生产环境下添加参数--restart=always ，只要容器存在就自动启动。



管理员账号用户：admin/adminadmin
### Run Docker As Root Account

```
docker exec -u 0 -it mycontainer bash

```

### Change Docker Root Passwd

```
RUN echo "root:Docker!" | chpasswd
```

### Interactive Shell Inside Docker Container:

docker run -u 0 -it --rm centos/python-35-centos7 /bin/bash


### Check Out Docker Container's Logs

docker logs -f container_name


### Get Docker Image Info

docker inspect container_name  

### 停掉容器
```
docker stop $(docker ps -a -q)
```

### 删除容器名称

```
docker rm $(docker ps -a -q)

```
删除images，通过image的id来指定删除谁

```
docker rmi <image id>
```

想要删除untagged images，也就是那些id为的image的话可以用

```
docker rmi $(docker images | grep "^<none>" | awk "{print $3}")
```

要删除全部image的话

```
docker rmi $(docker images -q)
```

强制删除全部image的话
```
docker rmi -f $(docker images -q)
```

## 离线安装Docker

从已有文件中，解压缩docker.zipo得到docker的安装文件。执行项目命令：
```
rpm -ivh --nodeps --replacefiles --replacepkgs libseccomp-2.3.1-3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs libseccomp-devel-2.3.1-3.el7.x86_64.rpm 
rpm -ivh --nodeps --replacefiles --replacepkgs libltdl7-2.4.2-alt8.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-selinux-17.03.3.ce-1.el7.noarch.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-cli-18.09.2-3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs containerd.io-1.2.2-3.3.el7.x86_64.rpm
rpm -ivh --nodeps --replacefiles --replacepkgs docker-ce-18.09.2-3.el7.x86_64.rpm
```
