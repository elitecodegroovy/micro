# 版本信息
FROM centos:7.9.2009
MAINTAINER liujigang "liujigang@tisson.cn"
#设置cmd
ENV nginxVersion 1.22.0
RUN mkdir -p /opt/src && cd /opt/src
RUN yum -y install wget gcc gcc-c++ make zlib-devel pcre-devel openssl-devel initscripts
RUN wget http://nginx.org/download/nginx-$nginxVersion.tar.gz && tar -xzf nginx-$nginxVersion.tar.gz \
&& ln -sf nginx-$nginxVersion nginx \
&& cd nginx \
&& ./configure \
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
--with-stream               \
&& make                     \
&& make install
RUN  groupadd nginx    && useradd nginx -g nginx
RUN chown -R nginx:nginx /opt/nginx
ADD ./nginx /etc/init.d/nginx
RUN  chmod +x /etc/init.d/nginx && chkconfig --add nginx  &&  chkconfig --level 345 nginx on
RUN  touch /etc/sysconfig/network
RUN yum -y remove wget unzip
ENV PATH="/opt/nginx/sbin:${PATH}"
ADD nginx.conf /opt/nginx/conf/nginx.conf
CMD ["nginx", "-g", "daemon off;"]
