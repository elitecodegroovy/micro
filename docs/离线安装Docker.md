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