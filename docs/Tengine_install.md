
新建tengine用户组

```  
groupadd -r nginx
useradd -r -g nginx -M nginx
```


## Install Dependencies

```
yum install gcc openssl-devel zlib-devel pcre-devel
```

## Install jemalloc

```
cd jemalloc-3.6.0.tar.bz2
./configure --prefix=/opt/jemalloc
make 
```

## 安装geoip2 lib步骤

```
cd /usr/local/src
rm -f libmaxminddb-1.4.2.tar.gz
wget https://github.com/maxmind/libmaxminddb/releases/download/1.4.2/libmaxminddb-1.4.2.tar.gz
tar -xzf libmaxminddb-1.4.2.tar.gz
cd libmaxminddb-1.4.2
yum install gcc gcc-c++ make -y
./configure
make
make check
sudo make install

echo '/usr/local/lib' > /etc/ld.so.conf.d/geoip.conf
sudo ldconfig
```

## ngx_http_geoip2_module 模块

``` 

cd /usr/local/src
wget https://github.com/leev/ngx_http_geoip2_module/archive/3.3.tar.gz
tar -xzf 3.3.tar.gz
mv ngx_http_geoip2_module-3.3 ngx_http_geoip2_module 
```
## Install Teninge

```
./configure --prefix=/opt/nginx \
--with-jemalloc=/opt/jemalloc \
--with-http_gzip_static_module \
--with-http_realip_module \
--with-http_stub_status_module 

make && make install
```


or 

``` 
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
        --with-http_gzip_static_module                \
        --with-http_addition_module                   \
        --with-jemalloc=/opt/jemalloc                 \
        --with-http_stub_status_module                \
        --with-http_ssl_module                        \
        --with-openssl=                               \
        --with-http_realip_module                     \
        --with-http_v2_module                         \
        --with-pcre                                   \
        --with-file-aio                               \
        --with-http_realip_module                     \
        --without-http_scgi_module                    \
        --without-http_uwsgi_module                   \
        --without-http_fastcgi_module                 \
        --with-stream                                 \
        --add-module=/home/app/ngx_http_geoip2_module \

make && make install
```


更改tengine的权限，将该目录给Nginx用户

```
chown -R nginx:nginx /opt/nginx
chmod -R 755 /opt/nginx
```

## script

vi /etc/init.d/nginx

```
#!/bin/sh
#
# nginx - this script starts and stops the nginx daemon
#
# chkconfig:   - 85 15
# description:  Nginx is an HTTP(S) server, HTTP(S) reverse \
#               proxy and IMAP/POP3 proxy server
# processname: nginx
# config:      /etc/nginx/nginx.conf
# config:      /etc/sysconfig/nginx
# pidfile:     /var/run/nginx.pid

# Source function library.
. /etc/rc.d/init.d/functions

# Source networking configuration.
. /etc/sysconfig/network

# Check that networking is up.
[ "$NETWORKING" = "no" ] && exit 0

nginx="/opt/nginx/sbin/nginx"
prog=$(basename $nginx)

NGINX_CONF_FILE="/opt/nginx/conf/nginx.conf"

[ -f /etc/sysconfig/nginx ] && . /etc/sysconfig/nginx

lockfile=/var/lock/subsys/nginx

make_dirs() {
   # make required directories
   user=`nginx -V 2>&1 | grep "configure arguments:" | sed 's/[^*]*--user=\([^ ]*\).*/\1/g' -`
   options=`$nginx -V 2>&1 | grep 'configure arguments:'`
   for opt in $options; do
       if [ `echo $opt | grep '.*-temp-path'` ]; then
           value=`echo $opt | cut -d "=" -f 2`
           if [ ! -d "$value" ]; then
               # echo "creating" $value
               mkdir -p $value && chown -R $user $value
           fi
       fi
   done
}

start() {
    [ -x $nginx ] || exit 5
    [ -f $NGINX_CONF_FILE ] || exit 6
    make_dirs
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
    sleep 1
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


nginx文件添加权限

```
chmod 755 /etc/init.d/nginx
```

开机启动
```
systemctl enable nginx
```


```
#启动
service nginx start
#停止
service nginx stop

#重启
service nginx reload

```