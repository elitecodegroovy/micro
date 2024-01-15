

## Install packages

yum install -y gcc-c++  pcre pcre-devel  zlib zlib-devel openssl openssl-devel 

看执行情况，aarch64架构的常见编译环境还是没有问题的.

``` 
groupadd nginx
useradd nginx -g nginx
chown -R nginx:nginx /opt/nginx
```

### 下载源码编译
``` 
wget -g https://nginx.org/download/nginx-1.19.6.tar.gz
tar xvf nginx-1.19.6.tar.gz
mv nginx-1.19.6 nginx
cd nginx 
./configure --sbin-path=/opt/nginx/nginx --conf-path=/opt/nginx/nginx.conf --pid-path=/opt/nginx/nginx.pid --with-pcre --with-http_stub_status_module --with-http_gzip_static_module --with-http_ssl_module


        ./configure \
        --user=nginx                                \
        --group=nginx                               \
        --prefix=/opt/nginx                         \
        --sbin-path=/opt/nginx/sbin/nginx           \
        --conf-path=/opt/nginx/conf/nginx.conf      \
        --pid-path=/opt/nginx/nginx.pid               \
        --lock-path=/opt/nginx/nginx.lock             \
        --error-log-path=/opt/nginx/logs/error.log  \
        --http-log-path=/opt/nginx/logs/access.log  \
        --with-http_gzip_static_module         \
        --with-http_addition_module            \
        --with-http_stub_status_module         \
        --with-http_ssl_module                 \
        --with-http_ssl_module           \
        --with-http_realip_module                    \
        --with-http_v2_module                        \
        --with-pcre                                  \
        --with-file-aio                        \
        --with-http_realip_module              \
        --without-http_scgi_module             \
        --without-http_uwsgi_module            \
        --without-http_fastcgi_module          \
        --with-stream

make

sudo make install
```

## 开机启动

vi /etc/init.d/nginx


``` 
#!/bin/sh
#
# nginx - this script starts and stops the nginx daemin
#
# chkconfig:   - 85 15
# description:  Nginx is an HTTP(S) server, HTTP(S) reverse \
#               proxy and IMAP/POP3 proxy server
# processname: nginx
# config:      /opt/nginx/nginx.conf
# pidfile:     /opt/nginx/nginx.pid
# user:        nginx

# Source function library.
. /etc/rc.d/init.d/functions

# Source networking configuration.
. /etc/sysconfig/network

# Check that networking is up.
[ "$NETWORKING" = "no" ] && exit 0

nginx="/opt/nginx/sbin/nginx"
prog=$(basename $nginx)

NGINX_CONF_FILE="/opt/nginx/conf/nginx.conf"

lockfile=/opt/nginx/nginx.lock

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

chmod +x /etc/init.d/nginx

``` 
chkconfig --add nginx
chkconfig --level 345 nginx on
```

## 运维

####开启、关闭、重启nginx

/opt/nginx/nginx 

/opt/nginx/nginx -s stop

/opt/nginx/nginx -s reload

修改配置

```
vi /opt/nginx/nginx.conf
```

开启运行者：

``` 
# 启动者账号

user    root;
```

访问地址：

``` 
http://localhost
```

或者命令：

``` 
curl http://localhost
```