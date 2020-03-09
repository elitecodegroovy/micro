
### Node Info

UI: 192.168.1.229:8500/ui
### Starting the Agent

```
consul agent -dev
```

附加配置文件
```
consul agent -dev -config-dir=/etc/consul.d
```

获取节点信息

>>http://localhost:8500/v1/catalog/nodes

返回：
```
[
  {
    "ID": "3422c716-04a3-2c88-db5c-5f7e6339a06e",
    "Node": "DESKTOP-10R9D2H",
    "Address": "127.0.0.1",
    "Datacenter": "dc1",
    "TaggedAddresses": {
      "lan": "127.0.0.1",
      "wan": "127.0.0.1"
    },
    "Meta": {
      "consul-network-segment": ""
    },
    "CreateIndex": 9,
    "ModifyIndex": 10
  }
]
```
查询web服务
>http://localhost:8500/v1/catalog/service/web

返回：

```
// 20180701125243
// http://localhost:8500/v1/catalog/service/web

[
  {
    "ID": "03123e86-f57c-8473-4743-3ec512c1c168",
    "Node": "DESKTOP-10R9D2H",
    "Address": "127.0.0.1",
    "Datacenter": "dc1",
    "TaggedAddresses": {
      "lan": "127.0.0.1",
      "wan": "127.0.0.1"
    },
    "NodeMeta": {
      "consul-network-segment": ""
    },
    "ServiceKind": "",
    "ServiceID": "web",
    "ServiceName": "web",
    "ServiceTags": [
      "rails"
    ],
    "ServiceAddress": "",
    "ServiceMeta": {
      
    },
    "ServicePort": 80,
    "ServiceEnableTagOverride": false,
    "ServiceProxyDestination": "",
    "ServiceConnect": {
      "Native": false,
      "Proxy": null
    },
    "CreateIndex": 10,
    "ModifyIndex": 10
  }
]
```

health check 
>http://localhost:8500/v1/health/service/web?passing

返回：

```

[
  {
    "Node": {
      "ID": "03123e86-f57c-8473-4743-3ec512c1c168",
      "Node": "DESKTOP-10R9D2H",
      "Address": "127.0.0.1",
      "Datacenter": "dc1",
      "TaggedAddresses": {
        "lan": "127.0.0.1",
        "wan": "127.0.0.1"
      },
      "Meta": {
        "consul-network-segment": ""
      },
      "CreateIndex": 9,
      "ModifyIndex": 10
    },
    "Service": {
      "ID": "web",
      "Service": "web",
      "Tags": [
        "rails"
      ],
      "Address": "",
      "Meta": null,
      "Port": 80,
      "EnableTagOverride": false,
      "ProxyDestination": "",
      "Connect": {
        "Native": false,
        "Proxy": null
      },
      "CreateIndex": 10,
      "ModifyIndex": 10
    },
    "Checks": [
      {
        "Node": "DESKTOP-10R9D2H",
        "CheckID": "serfHealth",
        "Name": "Serf Health Status",
        "Status": "passing",
        "Notes": "",
        "Output": "Agent alive and reachable",
        "ServiceID": "",
        "ServiceName": "",
        "ServiceTags": [
          
        ],
        "Definition": {
          
        },
        "CreateIndex": 9,
        "ModifyIndex": 9
      }
    ]
  }
]
```

### Consul Clustering With cmd

````
consul agent -server -bootstrap-expect=1 \
    -data-dir=/tmp/consul -node=agent-one -bind=192.168.1.225 \
    -enable-script-checks=true -config-dir=/etc/consul.d
	
	consul agent -server -bootstrap-expect=1 \
    -data-dir=/tmp/consul -node=agent-two—server -bind=192.168.1.224 \
    -enable-script-checks=true -config-dir=/etc/consul.d	

	consul agent -server  \
    -data-dir=/tmp/consul -node=agent-three-server -bind=192.168.1.226 \
    -enable-script-checks=true -config-dir=/etc/consul.d
	
consul agent -data-dir=/tmp/consul -node=agent-two \
    -bind=192.168.1.224 -enable-script-checks=true -config-dir=/etc/consul.d
	
consul join 172.20.20.11
```
### Consul Clustering With Config File
configuration file:

```
root@server1:~#vim /etc/consul.d/server/config.json
```

编辑内容：
```
{
"bind_addr": "192.168.1.225", 
"datacenter": "dc1",
"data_dir": "/opt/consul",
"encrypt": "HYn7j/ytTq37y9MvqRqqxw==",
"log_level": "INFO",
"enable_syslog": true,
"enable_debug": true,
"node_name": "consul-server1",
"server": true,
"bootstrap_expect": 3,
"leave_on_terminate": false,
"skip_leave_on_interrupt": true,
"rejoin_after_leave": true,
"retry_join": [ 
  "192.168.1.225:8301",
  "192.168.1.226:8301",
  "192.168.1.224:8301"
  ]
 }
 
 
 
 {
"bind_addr": "192.168.1.224", 
"datacenter": "dc1",
"data_dir": "/opt/consul",
"encrypt": "HYn7j/ytTq37y9MvqRqqxw==",
"log_level": "INFO",
"enable_syslog": true,
"enable_debug": true,
"node_name": "consul-server2",
"server": true,
"bootstrap_expect": 3,
"leave_on_terminate": false,
"skip_leave_on_interrupt": true,
"rejoin_after_leave": true,
"retry_join": [ 
  "192.168.1.225:8301",
  "192.168.1.226:8301",
  "192.168.1.224:8301"
  ]
 }
 
 
  {
"bind_addr": "192.168.1.226", 
"datacenter": "dc1",
"data_dir": "/opt/consul",
"encrypt": "HYn7j/ytTq37y9MvqRqqxw==",
"log_level": "INFO",
"enable_syslog": true,
"enable_debug": true,
"node_name": "consul-server3",
"server": true,
"bootstrap_expect": 3,
"leave_on_terminate": false,
"skip_leave_on_interrupt": true,
"rejoin_after_leave": true,
"retry_join": [ 
  "192.168.1.225:8301",
  "192.168.1.226:8301",
  "192.168.1.224:8301"
  ]
 }
 ```
 启动：
 ```
 consul agent -config-dir=/etc/consul.d/server/config.json
 ```
 
 客户端：
 
 vi /etc/consule.d/client/config.json
```
 {
 "bind_addr": "192.168.1.229",
 "client_addr": "192.168.1.229",
 "datacenter": "dc1",
 "data_dir": "/opt/consul",
 "encrypt": "HYn7j/ytTq37y9MvqRqqxw==",
 "log_level": "INFO",
 "enable_syslog": true,
 "enable_debug": true,
 "node_name": "consul-client1",
 "enable_script_checks":true,
 "server": false,
 "rejoin_after_leave": true,
 "retry_join": [
  "192.168.1.225:8301",
  "192.168.1.226:8301",
  "192.168.1.224:8301"
    ]
 }
 ```
  启动：
```
 consul agent -config-dir=/etc/consul.d/client/
```

获取所有节点消息
>http://192.168.1.229:8500/v1/catalog/nodes


### Consul K/V

Get

```
consul kv get redis/config/minconns
```

Put

```
$ consul kv put redis/config/minconns 1
Success! Data written to: redis/config/minconns

$ consul kv put redis/config/maxconns 25
Success! Data written to: redis/config/maxconns

$ consul kv put -flags=42 redis/config/users/admin abcd1234
Success! Data written to: redis/config/users/admin
```

GetALL

```
consul kv get -recurse
```

Delete
```
consul kv delete redis/config/minconns
```

 Check-And-Set operation
 
```
 consul kv put -cas -modify-index=123 foo bar
 Success! Data written to: foo

 $ consul kv put -cas -modify-index=123 foo bar
 Error! Did not write to foo: CAS failed
```
​

CONSUL_HTTP_ADDR

```
session := engine.NewSession()
defer session.Close()

//TODO ...事务操作开始区域
// 开启 Begin()事务
err := session.Begin()  

//执行语句操作步骤1
//user2 := Userinfo{Username: "yyy"}
//_, err = session.Where("id = ?", 2).Update(&user2)
//if err != nil {
//    session.Rollback()
//    return
//}

//执行语句操作步骤2
//_, err = session.Exec("delete from userinfo where username = ?", user2.Username)
//if err != nil {
//    session.Rollback()
//    return
//}

//TODO ...事务操作结束区域
// 提交所有事务
err = session.Commit()
if err != nil {
    return
}
```