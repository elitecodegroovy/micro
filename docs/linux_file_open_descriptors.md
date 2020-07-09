1、文件中添加如下

```
vi /etc/sysctl.conf

fs.file-max = 202400
fs.nr_open = 102400
vm.swappiness = 0
net.core.somaxconn = 32768
net.ipv4.tcp_syncookies = 0

```

保存并退出。

立即生效
>sysctl -p 


文件中添加如下：

```
vi /etc/security/limits.conf

*   		soft     nofile  	 102400
*   		hard     nofile  	 102400

```


命令行（立即生效）：

```
ulimit -n 102400 

```