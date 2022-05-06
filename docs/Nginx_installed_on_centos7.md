
#Install dependencies for nginx
We have few pre-requisites to be installed to compile which include development libraries along with source code compilers.

```
yum -y install gcc gcc-c++ make zlib-devel pcre-devel openssl-devel
```
Lets first create a directory to store our source code:

```
mkdir -p src && cd src
```
## 编译最新的openssl

下载源码包
wget https://www.openssl.org/source/openssl-1.1.1i.tar.gz
编译安装
# tar xvf openssl-1.1.1i.tar.gz -C /data/
# ln -sv openssl-1.1.1i openssl
# cd openssl
# ./config
#make -j 4 && make install
报错
升级后如果执行 openssl version 命令出现openssl: error while loading shared libraries: libssl.so.1.1: cannot open shared object file: No such file or directory错误。执行以下命令即可。

# ln -sv /data/openssl/libssl.so.1.1 /usr/lib64/libssl.so.1.1
# ln -sv /data/openssl/libcrypto.so.1.1 /usr/lib64/libcrypto.so.1.1
更新软连接
# cp /usr/bin/openssl{,.bak}
# ln -sfv /data/openssl/apps/openssl /usr/bin/openssl


查看版本
#未编译前
# openssl version
OpenSSL 1.0.2k-fips  26 Jan 2017
#编译后
# openssl version
OpenSSL 1.1.1i  8 Dec 2020




# Compiling from Source

## Downloading the source code
Lets get the current nginx version number from http://nginx.org/en/download.html

Run the following commands to download the source.
```

nginxVersion="1.21.3"
wget http://nginx.org/download/nginx-$nginxVersion.tar.gz
tar -xzf nginx-$nginxVersion.tar.gz 
ln -sf nginx-$nginxVersion nginx
```
Preparing the nginx source

We first want to prepare nginx with necessary basic options.

For a full list of options you can look at ./configure --help

Options for basic file path names
These options are the basic variables which we override to use default system paths at /etc/ to ensure it works simliar when installed via rpm. The user and group option are used to run the nginx worker processes in non-privileged.

```
--user
--group
--prefix
--sbin-path
--conf-path
--pid-path
--lock-path
--error-log-path
--http-log-path

```
Other options


- --with-http_gzip_static_module option enables nginx to use gzip (Before serving a file from disk to a gzip-enabled client, this module will look for a precompressed file in the same location that ends in ".gz". The purpose is to avoid compressing the same file each time it is requested.).[recommended for reducing size of information sent]


- --with-http_stub_status_module option enables other plugins over nginx to allow us to get the status (This module provides the ability to get some status from nginx.). [recommended for getting stats]


- --with-http_ssl_module - required if you want to run a HTTPS server. See How To Create a SSL Certificate on nginx for CentOS 6


- --with-pcre option enables to match routes via Regular Expression Matching when defining routes. [recommended, you will find more use of this once you start adding and matching routes]


- --with-file-aio - enables asynchronous I/O, better than the default send file option (recommended if you are allowing users to download static files)


- --with-http_realip_module is used for getting the IP of the client when behind a load balancer. This is useful when serving content behind CloudFlare like services.


- --without-http_scgi_module - Disable SCGI module (normally used when running CGI scripts)


- --without-http_uwsgi_module - Disable UWSGI module (normally used when running CGI scripts)


- --without-http_fastcgi_module - Disable FastCGI module (normally used when running CGI scripts)
Our configuration options
```
        cd nginx
        ./configure \
        --user=nginx                                \
        --group=nginx                               \
        --prefix=/data/nginx                         \
        --sbin-path=/data/nginx/sbin/nginx           \
        --conf-path=/data/nginx/conf/nginx.conf      \
        --pid-path=/data/nginx/nginx.pid               \
        --lock-path=/data/nginx/nginx.lock             \
        --error-log-path=/data/nginx/logs/error.log  \
        --http-log-path=/data/nginx/logs/access.log  \
        --with-http_gzip_static_module        \
        --with-http_addition_module           \
        --with-http_stub_status_module        \
        --with-http_ssl_module                \
        --with-openssl=                       \
        --with-http_realip_module             \
        --with-http_v2_module                 \
        --with-pcre                           \
        --with-file-aio                       \
        --with-http_realip_module             \
        --without-http_scgi_module            \
        --without-http_uwsgi_module           \
        --without-http_fastcgi_module         \
        --with-stream                         
```
 wget https://github.com/openresty/headers-more-nginx-module/archive/refs/tags/v0.33.tar.gz

tar -xvf v0.33.tar.gz
mv  mv headers-more-nginx-module-0.33 headers-more-nginx-module

# --add-module=headers-more-nginx-module的路径

        ./configure \
        --user=nginx                                \
        --group=nginx                               \
        --prefix=/data/nginx                         \
        --sbin-path=/data/nginx/sbin/nginx           \
        --conf-path=/data/nginx/conf/nginx.conf      \
        --pid-path=/data/nginx/nginx.pid               \
        --lock-path=/data/nginx/nginx.lock             \
        --error-log-path=/data/nginx/logs/error.log  \
        --http-log-path=/data/nginx/logs/access.log  \
        --with-http_gzip_static_module         \
        --with-http_addition_module            \
        --with-http_stub_status_module         \
        --with-http_ssl_module                 \
        --with-openssl=/data/openssl           \
        --add-module=/data/headers-more-nginx-module \
        --with-http_realip_module                    \
        --with-http_v2_module                        \
        --with-pcre                                  \
        --with-file-aio                        \
        --with-http_realip_module              \
        --without-http_scgi_module             \
        --without-http_uwsgi_module            \
        --without-http_fastcgi_module          \
        --with-stream       


Compiling the nginx source
Once we are able to configure the source which even checks for additional requirements like the compiler(gcc, g++) which we installed in the pre-requisites step:

```
 make
 make install

```


## Running the VPS

1. Add the user nginx to the system. This is a one time command:

```
groupadd nginx
useradd nginx -g nginx
chown -R nginx:nginx /data/nginx
```

2. We need to setup the file /etc/init.d/nginx to run when system starts:

```
#!/bin/sh
#
# nginx - this script starts and stops the nginx daemin
#
# chkconfig:   - 85 15
# description:  Nginx is an HTTP(S) server, HTTP(S) reverse \
#               proxy and IMAP/POP3 proxy server
# processname: nginx
# config:      /data/nginx/nginx.conf
# pidfile:     /data/nginx/nginx.pid
# user:        nginx

# Source function library.
. /etc/rc.d/init.d/functions

# Source networking configuration.
. /etc/sysconfig/network

# Check that networking is up.
[ "$NETWORKING" = "no" ] && exit 0

nginx="/data/nginx/sbin/nginx"
prog=$(basename $nginx)

NGINX_CONF_FILE="/data/nginx/conf/nginx.conf"

lockfile=/data/nginx/nginx.lock

start() {
    [ -x $nginx ] || exit 5
    [ -f $NGINX_CONF_FILE ] || exit 6
    echo -n $"Starting $prog: "
    daemon $nginx -c $NGINX_CONF_FILE
    retval=$?
    echo
    [ $retval -eq 0 ] && touch $lockfile
    return $retval
}

stop() {
    echo -n $"Stopping $prog: "
    killproc $prog -QUIT
    retval=$?
    echo
    [ $retval -eq 0 ] && rm -f $lockfile
    return $retval
}

restart() {
    configtest || return $?
    stop
    start
}

reload() {
    configtest || return $?
    echo -n $"Reloading $prog: "
    killproc $nginx -HUP
    RETVAL=$?
    echo
}

force_reload() {
    restart
}

configtest() {
  $nginx -t -c $NGINX_CONF_FILE
}

rh_status() {
    status $prog
}

rh_status_q() {
    rh_status >/dev/null 2>&1
}

case "$1" in
    start)
        rh_status_q && exit 0
        $1
        ;;
    stop)
        rh_status_q || exit 0
        $1
        ;;
    restart|configtest)
        $1
        ;;
    reload)
        rh_status_q || exit 7
        $1
        ;;
    force-reload)
        force_reload
        ;;
    status)
        rh_status
        ;;
    condrestart|try-restart)
        rh_status_q || exit 0
            ;;
    *)
        echo $"Usage: $0 {start|stop|status|restart|condrestart|try-restart|reload|force-reload|configtest}"
        exit 2
esac
```

Optionally, you can obtain the source from:

```
wget -O /etc/init.d/nginx https://gist.github.com/sairam/5892520/raw/b8195a71e944d46271c8a49f2717f70bcd04bf1a/etc-init.d-nginx
```


This file should be made executable so that we can use it via 'service nginx ':

```
mv -f nginx /etc/init.d/nginx
chmod +x /etc/init.d/nginx
```

3. Set the service to start whenever the system boots:

```
chkconfig --add nginx
chkconfig --level 345 nginx on

```

4. Configure /etc/nginx/nginx.conf to set types_hash_bucket_size and server_names_hash_bucket_size which needs to be increased.

```
http {
    include       mime.types;
    default_type  application/octet-stream;
    # add the below 2 lines under http around line 20
    types_hash_bucket_size 64;
    server_names_hash_bucket_size 128;
```

5. Start the server. This will start the VPS on port 80.

```
service nginx start
```

优化Linux Open Files


limits.conf
修改/etc/security/limits.conf文件

在文件末尾添加

``` 
* soft nofile 204800
* hard nofile 204800
* soft nproc 204800
* hard nproc 204800
```
  说明：

*             代表针对所有用户 
noproc     是代表最大进程数
nofile     是代表最大文件打开数

光修改上面一个文件是不行的，还需要修改一个文件。
修改/etc/security/limits.d/20-nproc.conf文件

删掉默认配置，修改如下：
``` 
*          soft    nproc     204800
*          hard    nproc     204800
```

修改2个配置文件之后，重启后生效

reboot -f


fixed an issues:

The --add-icmp-block=<type> option can be used to block a certain type.
``` 
firewall-cmd --add-icmp-block=echo-reply --permanent
firewall-cmd --add-icmp-block=time-exceeded --permanent
firewall-cmd --add-icmp-block=destination-unreachable --permanent
firewall-cmd --add-icmp-block=timestamp-reply --permanent
firewall-cmd --add-icmp-block=timestamp-request --permanent
 
firewall-cmd --reload
```

The --remove-icmp-block=<type> option can be used to not block a certain type.

``` 
firewall-cmd --remove-icmp-block=echo-request --permanent
```

