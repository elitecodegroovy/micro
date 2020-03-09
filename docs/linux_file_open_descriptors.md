1、文件中添加如下

```
vi /etc/sysctl.conf

#file-max是内核可分配的最大文件数
fs.file-max = 202400 
#nr_open是单个进程可分配的最大文件数
fs.nr_open = 102400 

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