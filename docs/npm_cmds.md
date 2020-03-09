
## npm基本命名

查看所有全局安装的模块:
```
npm list -g
```


看某个模块的版本号:
```
npm list grunt
```

命令来卸载 Node.js 模块:
```
npm uninstall express
```

更新模块:
```
npm update express
```

搜索模块:
```
npm search express
```


 npm 资源库中注册用户（使用邮箱注册）:
 
```
npm adduser
Username: mcmohd
Password:
Email: (this IS public) mcmohd@gmail.com
```

命令来发布模块：
```
$ npm publish

```

在package.json所在目录下使用npm install . -g可先在本地安装当前命令行程序，可用于发布前的本地测试。

使用npm update <package>可以把当前目录下node_modules子目录里边的对应模块更新至最新版本。

使用npm update <package> -g可以把全局安装的对应命令行程序更新至最新版。

使用npm cache clear可以清空NPM本地缓存，用于对付使用相同版本号发布新版本代码的人。

使用npm unpublish <package>@<version>可以撤销发布自己发布过的某个版本代码。

## 使用淘宝 NPM 镜像
淘宝 NPM 镜像是一个完整 npmjs.org 镜像，你可以用此代替官方版本(只读)，同步频率目前为 10分钟 一次以保证尽量与官方服务同步。
```
npm install -g cnpm --registry=https://registry.npm.taobao.org

```

这样就可以使用 cnpm 命令来安装模块了：

```
$ cnpm install [name]
```

## node.js事件

Node.js 是单进程单线程应用程序，但是因为 V8 引擎提供的异步执行回调接口，通过这些接口可以处理大量的并发，所以性能非常高。

Node.js 使用事件驱动模型，当web server接收到请求，就把它关闭然后进行处理，然后去服务下一个web请求。

当这个请求完成，它被放回处理队列，当到达队列开头，这个结果被返回给用户。

这个模型非常高效可扩展性非常强，因为webserver一直接受请求而不等待任何读写操作。（这也被称之为非阻塞式IO或者事件驱动IO）

在事件驱动模型中，会生成一个主循环来监听事件，当检测到事件时触发回调函数。

![](https://www.runoob.com/wp-content/uploads/2015/09/event_loop.jpg)

Node.js 有多个内置的事件，我们可以通过引入 events 模块，并通过实例化 EventEmitter 类来绑定和监听事件。