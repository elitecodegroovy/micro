
## ElasticSearch（单节点安装）

使用curl命令下载tar包
```
# curl -L -O https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.6.12.tar.gz
也可以使用wget命令下载
# wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.6.12.tar.gz

```
解压至指定目录(安装目录)/usr/local/下，并将其重命名为elasticsearch (完整的目录应该为/usr/local/elasticsearch)

```
# tar -zxf elasticsearch-5.6.12.tar.gz -C /usr/local/
# cd /usr/local/
# mv elasticsearch-5.6.12 elasticsearch
```

创建一个用于运行elasticsearch的普通用户，随后修改elasticsearch家目录的所属权限为该用户；创建elasticsearch数据存储目录/data/elasticsearch

```
# groupadd elasticsearch
# useradd -g elasticsearch elasticsearch -m
# chown -R elasticsearch. /usr/local/elasticsearch
# mkdir /data/elasticsearch
# chown -R elasticsearch. /data/elasticsearch
```

修改elasticsearch.yml配置文件
```
# vim config/elasticsearch.yml
cluster.name: my-application        #ELK集群名称
path.data: /data/elasticsearch      #elasticsearch 数据存储目录
path.logs: /usr/local/elasticsearch/logs   #elasticsearch 日志存储路径
network.host: 10.17.83.42           #elasticsearch 监听地址，默认为localhost
http.port: 9200                     #elasticsearch 监听端口，默认问9200

```

修改相关的内核参数
```
vim /etc/security/limits.conf
```
# 打开/etc/security/limits.conf文件，添加以下参数

*        soft    nproc           2048
*        hard    nproc           16384
*        soft    nofile          65536
*        hard    nofile          65536

修改vm.max_map_count=262144

```
# echo "vm.max_map_count=262144" >> /etc/sysctl.conf
```

启动设置：

```
sysctl -p
```

运行elasticsearch (注意：要切换到普通用户运行)


```
# su - elasticsearch

$ ./bin/elasticsearch
```


运行elasticsearch
```
[elasticsearch@elk-server ~]$ cd /usr/local/elasticsearch/
[elasticsearch@elk-server elasticsearch]$ ./bin/elasticsearch

```

一般情况我们要求elasticsearch在后台运行，使用命令如下：
```
$ ./bin/elasticsearch -d
```

检查elasticsearch状态，如下则表示正常运行
```
[root@elk-server elasticsearch]# curl http://10.17.83.42:9200
{
  "name" : "6eN59Pf",
  "cluster_name" : "my-application",
  "cluster_uuid" : "cKopwE1iTciIQAiFI6d8Gw",
  "version" : {
    "number" : "5.5.2",
    "build_hash" : "b2f0c09",
    "build_date" : "2017-08-14T12:33:14.154Z",
    "build_snapshot" : false,
    "lucene_version" : "6.6.0"
  },
  "tagline" : "You Know, for Search"
}
```

## Logstash安装

下载logstash软件包
```
# wget https://artifacts.elastic.co/downloads/logstash/logstash-5.6.12.tar.gz

```
2、解压至指定安装目录
```
# tar -zxf logstash-5.5.12.tar.gz -C /usr/local/
# cd /usr/local/
 
# mv logstash-5.6.12 logstash
```


## Kibana安装



1、下载kibana
```
# wget https://artifacts.elastic.co/downloads/kibana/kibana-5.6.12-linux-x86_64.tar.gz
```

2、解压至安装目录

```
# tar -zxf kibana-5.6.12-linux-x86_64.tar.gz -C /usr/local/
# cd /usr/local/
# mv kibana-5.6.12-linux-x86_64 kibana
```

3、修改配置
```
# cd kibana/
 
# vim config/kibana.yml
 
server.port: 5601            # 监听端口
server.host: "10.17.83.42"   # 指定后端服务器
elasticsearch.url: "http://10.17.83.42:9200"  # 指定elasticsearch实例地址
```



4、运行kibana


```

[root@elk-server kibana]# ./bin/kibana &

```

5 验证kibana

在客户端浏览器打开http://10.17.83.42:5601

## 收集syslog日志

1、创建配置文件


```
# cd logstash/

# vim config/logstash.conf
input {
    file {
        path => ["/var/log/messages"]
        type => "syslog"
    }
}
 
filter {
    grok {
        match => [ "message", "%{SYSLOGBASE} %{GREEDYDATA:content}" ]
    }
}
 
output {
    elasticsearch {
        hosts => ["10.17.83.42:9200"]
        index => "syslog-%{+YYY.MM.dd}"
    }
    stdout { codec => rubydebug }
}




```

其中match => [ "message", "%{SYSLOGBASE} %{GREEDYDATA:content}" ]这一行也可以具体写成如下：

 

match => [ "message", "%{SYSLOGTIMESTAMP:timestamp} (?:%{SYSLOGFACILITY} )?%{SYSLOGHOST:logsource} %{SYSLOGPROG}: %{GREEDYDATA:content}" ]


参考官方文档logstash配置：

https://www.elastic.co/guide/en/logstash/current/plugins-filters-grok.html


2、指定配置文件运行logstash


```

[root@elk-server logstash]# ./bin/logstash -f ./config/logstash.conf &

````
3、登录kibana页面

 

点击Management --> Index Patterns --> Create index

在Index name or pattern处输入在logstash中指定的index，后面的日期直接用*号代替即可

```
syslog-*
```


4、验证是否正常收集syslog日志，执行以下命令手动生成日志
```
[root@elk-server kibana]# logger "helloooooooo22"
[root@elk-server kibana]# yum install httpd

```

5 Kibana的dicover浏览日志图表。
例如：

http://192.168.1.229:5601/app/kibana#/discover?_g=()&_a=(columns:!(_source),index:AWYz2M_7GWTsCeNG0DPP,interval:auto,query:(match_all:()),sort:!(_score,desc))


