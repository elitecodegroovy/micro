#!/bin/bash
#CREATE BY ljg
#CREATE DATE 2021-12-17

# configure nginx version
nginxVersion="1.21.4"
NGINX_STARTUP_FILE=/etc/init.d/nginx
NGINX_CONFIG_FILE=/opt/nginx/conf/nginx.conf

yum -y install gcc gcc-c++ make zlib-devel pcre-devel openssl-devel

# create the user
groupadd nginx
useradd nginx -g nginx
mkdir -p /opt/nginx
rm -rf /opt/nginx/*
mkdir -p /opt/nginx/logs

mkdir -p /opt/src
rm -rf /opt/src/*
cd /opt/src

wget http://nginx.org/download/nginx-$nginxVersion.tar.gz
tar -xzf nginx-$nginxVersion.tar.gz
ln -sf nginx-$nginxVersion nginx

cd nginx
./configure \
    --user=nginx                                \
    --group=nginx                               \
    --prefix=/opt/nginx                         \
    --sbin-path=/opt/nginx/sbin/nginx           \
    --conf-path=$NGINX_CONFIG_FILE                \
    --pid-path=/opt/nginx/nginx.pid             \
    --lock-path=/opt/nginx/nginx.lock           \
    --error-log-path=/opt/nginx/logs/error.log  \
    --http-log-path=/opt/nginx/logs/access.log  \
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

tee $NGINX_STARTUP_FILE <<-'EOF'
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
EOF
make
make install
tee $NGINX_CONFIG_FILE <<-'EOF'
#user  nobody;
worker_processes  1;

#error_log  logs/error.log;
#error_log  logs/error.log  notice;
#error_log  logs/error.log  info;
#error_log  "pipe:rollback logs/error_log interval=1d baknum=7 maxsize=2G";

#pid        logs/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    #log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
    #                  '$status $body_bytes_sent "$http_referer" '
    #                  '"$http_user_agent" "$http_x_forwarded_for"';

    #access_log  logs/access.log  main;
    #access_log  "pipe:rollback logs/access_log interval=1d baknum=7 maxsize=2G"  main;

    sendfile        on;
    #tcp_nopush     on;

    #keepalive_timeout  0;
    keepalive_timeout  65;

    #gzip  on;
    # upstream - web
    upstream OPEN_PAAS {
        server 172.16.27.62:8001 max_fails=1  fail_timeout=30s;
    }
    upstream OPEN_PAAS_APPENGINE {
        server 172.16.27.62:8000 max_fails=1  fail_timeout=30s;
    }
    upstream OPEN_PAAS_ESB {
        server 172.16.27.62:8002 max_fails=1  fail_timeout=30s;
    }
    upstream OPEN_PAAS_LOGIN {
        server 172.16.27.62:8003 max_fails=1  fail_timeout=30s;
    }
    
    # upstream - paasagent
    upstream PAAS_AGENT_TEST {
        server 192.168.1.1:8085 max_fails=1  fail_timeout=30s;
    }
    upstream PAAS_AGENT_PROD {
        server 192.168.1.2:8085 max_fails=1  fail_timeout=30s;
    }
    # proxy_next_upstream  http_502 http_504 error timeout invalid_header;
    
    server {
        listen       80;
        server_name  www.bking.com;
    
        access_log /opt/nginx/access.log;
        error_log /opt/nginx/error.log;
    
        client_max_body_size    512m;
    
        # ============================ paas ============================
        # PAAS_SERVICE HOST/PORT
        location / {
#            proxy_pass http://OPEN_PAAS;
#            proxy_pass_header Server;
#            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
#            proxy_set_header X-Real-IP $remote_addr;
#            proxy_set_header X-Scheme $scheme;
#            proxy_set_header Host $http_host;
#            proxy_redirect off;
#            proxy_read_timeout 600;
            root   html;
            index  index.html index.htm;
        }
    
        # PAAS_SERVICE HOST/PORT, for doc
        location ~ ^/doc/(.*) {
            proxy_pass http://OPEN_PAAS/static/doc/$1$is_args$args;
            proxy_pass_header Server;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Scheme $scheme;
            proxy_set_header Host $http_host;
            proxy_redirect off;
            proxy_read_timeout 600;
    
        }
    
    
        # ============================ appengine ============================
        # ENGINE_SERVICE HOST/PORT
        location ~ ^/v1 {
            proxy_pass http://OPEN_PAAS_APPENGINE;
            proxy_pass_header Server;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header Host $http_host;
            proxy_set_header X-Scheme $scheme;
            proxy_read_timeout 600;
        }
    
        # ============================ esb ============================
        # ESB_SERVICE HOST/PORT
        location ~ ^/api/(.*) {
            proxy_pass http://OPEN_PAAS_ESB/$1$is_args$args;
            proxy_pass_header Server;
            proxy_set_header X-Request-Uri $request_uri;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Scheme $scheme;
            proxy_set_header Host $http_host;
            proxy_redirect off;
            proxy_read_timeout 600;
        }
    
    
        # ============================ login ============================
        # LOGIN_SERVICE HOST/PORT
        location ~ ^/login/(.*) {
            proxy_pass http://OPEN_PAAS_LOGIN/$1$is_args$args;
            proxy_pass_header Server;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Scheme $scheme;
            proxy_set_header Host $http_host;
            proxy_redirect off;
            proxy_read_timeout 600;
        }
    
    
        # ============================ paas_agent ============================
        # for apps test
        location ~ ^/t/ {
            proxy_pass http://PAAS_AGENT_TEST;
            proxy_pass_header Server;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Scheme $scheme;
            proxy_set_header Host $http_host;
            proxy_redirect off;
            proxy_read_timeout 600;
        }
    
        # for apps prod
        location ~ ^/o/ {
            proxy_pass http://PAAS_AGENT_PROD;
            proxy_pass_header Server;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Scheme $scheme;
            proxy_set_header Host $http_host;
            proxy_redirect off;
            proxy_read_timeout 600;
        }
    
    }
}
EOF
chown -R nginx:nginx /opt/nginx
chmod +x /etc/init.d/nginx
chmod +x $NGINX_CONFIG_FILE
systemctl enable nginx
# startup nginx service
service nginx start
echo 'Done!'