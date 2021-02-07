

## Install packages

yum install -y gcc-c++  pcre pcre-devel  zlib zlib-devel openssl openssl-devel 

看执行情况，aarch64架构的常见编译环境还是没有问题的.

### 下载源码编译
``` 
wget -g https://nginx.org/download/nginx-1.19.6.tar.gz
tar xvf nginx-1.19.6.tar.gz
mv nginx-1.19.6 nginx
cd nginx 
./configure --sbin-path=/usr/local/nginx/nginx --conf-path=/usr/local/nginx/nginx.conf --pid-path=/usr/local/nginx/nginx.pid --with-pcre --with-http_stub_status_module --with-http_gzip_static_module --with-http_ssl_module


make

sudo make install
```

## 运维

####开启、关闭、重启nginx

/usr/local/nginx/nginx 

/usr/local/nginx/nginx -s stop

/usr/local/nginx/nginx -s reload

修改配置

```
vi /usr/local/nginx/nginx.conf
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