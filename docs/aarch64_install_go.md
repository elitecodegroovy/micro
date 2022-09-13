
# 
飞腾1500A平台上构建Go语言环境指南

一、在X86平台上交叉编译出面向ARM64平台的Go语言自举编译工具链

所用工具 64位Linux操作系统的X86电脑

https://golang.org/dl/

下载 go1.8.1.linux-amd64.tar.gz

下载go1.8.1.src.tar.gz
解压压缩

在跟目录下建立两个文件夹
goarm64
gox64

把解压的内容拷贝两个目录

https://storage.googleapis.com/golang/go1.8.1.src.tar.gz 到 goarm64

https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz 到 gox64
下面我们用gox64目录中的go 去交叉编译 goarm64的go 生成安装环境包
在命令控制台执行
$ export GOROOT_BOOTSTRAP=/gox64
切换到goarm64的src

$ cd /goarm64/src
执行
$ GOOS=linux GOARCH=arm64 ./bootstrap.bash

龙芯CPU 架构：
1. 1.19版本 Go src支持龙心指令 （GOARCH=loong64）
$ GOOS=linux GOARCH=loong64 ./bootstrap.bash

2. mips64
mips64el 版本（uname -r）
$ GOOS=linux GOARCH=mips64le ./bootstrap.bash
执行完成根目录会生成
go-linux-arm64-bootstrap.tbz


二、在飞腾1500A平台上从源代码构建Go 1.8开发环境。
把go1.8.1.src.tar.gz和go-linux-arm64-bootstrap.tbz 拷贝到 arm64 平台的机器上
下面要用 go-linux-arm64-bootstrap.tbz 去编译go1.8.1.src.tar.gz
分别解压两个压缩包
假定：
go1.8.1.src.tar.gz解压到根目录的go1.8.1
go-linux-arm64-bootstrap.tbz 解压到根目录的 goarm64
验证
go-linux-arm64e-bootstrap.tbz里面的go可以运行

$ cd /goarm64/bin
$ export GOROOT=/goarm64
$ ./go version
如果观察到以下输出，则基本可以说明交叉编译没有问题。
go version go1.8.1 linux/arm64

cd  /go1.8.1
在当前目录创建名为env.sh的shell脚本，输入以下语句后保存退出。
#!/bin/bash
export GOROOT_BOOTSTRAP=/goarm64
export GOROOT=/go1.8.1
# Added for Loongson
export GO_TEST_TIMEOUT_SCALE=2
然后开始安装
$ source env.sh
$ cd src
$ ./all.bash

等吧。测试成功 ok！
最后
在~/.bashrc的末尾添加如下语句

export GOROOT=/root/workspace/go
export PATH=$PATH:$GOROOT/bin

然后执行
$ source ~/.bashrc

测试Go语言：

``` 
package main
import "fmt"
func main() {
    fmt.Printf("Welcome to Golang World!\n")
}
```

Compile and run the application.
``` 
go build hello.go
./hello
```
