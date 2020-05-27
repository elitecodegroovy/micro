
#Install dependencies for nginx
We have few pre-requisites to be installed to compile which include development libraries along with source code compilers.

```
yum -y install gcc gcc-c++ make zlib-devel pcre-devel openssl-devel
```
Lets first create a directory to store our source code:

```
mkdir -p src && cd src
```

# Compiling from Source

## Downloading the source code
Lets get the current nginx version number from http://nginx.org/en/download.html

Run the following commands to download the source.
```
nginxVersion="1.14.2"
wget http://nginx.org/download/nginx-1.14.2.tar.gz
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
        --prefix=/opt/nginx                         \
        --sbin-path=/opt/nginx/sbin/nginx           \
        --conf-path=/opt/nginx/conf/nginx.conf      \
        --pid-path=/opt/nginx/nginx.pid               \
        --lock-path=/opt/nginx/nginx.lock             \
        --error-log-path=/opt/nginx/logs/error.log  \
        --http-log-path=/opt/nginx/logs/access.log  \
        --with-http_gzip_static_module        \
        --with-http_addition_module           \
        --with-jemalloc=/opt/jemalloc         \
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
        --with-stream                         \
        --add-module=/home/app/ngx_http_geoip2_module
```

Compiling the nginx source
Once we are able to configure the source which even checks for additional requirements like the compiler(gcc, g++) which we installed in the pre-requisites step:

```
 make
 make install

```


## Running the VPS

1. Add the user nginx to the system. This is a one time command:

```
useradd -r nginx
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
# config:      /etc/nginx/nginx.conf
# pidfile:     /var/run/nginx.pid
# user:        nginx

# Source function library.
. /etc/rc.d/init.d/functions

# Source networking configuration.
. /etc/sysconfig/network

# Check that networking is up.
[ "$NETWORKING" = "no" ] && exit 0

nginx="/usr/sbin/nginx"
prog=$(basename $nginx)

NGINX_CONF_FILE="/etc/nginx/nginx.conf"

lockfile=/var/run/nginx.lock

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

