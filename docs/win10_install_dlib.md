安装dlib 需要先安装boost和cmake，而且cmake需要使用c编译器，所以需要装VS，版本最好VS2015以上



1 安装cmake

官网下载安装包：https://cmake.org/download/



我下载了win64的msi安装包，注意电脑是64位的就安装64位的安装包。



就一直next，最后finish

（我在安装结束后，运行还是有问题，发现是VS的问题，VS2008不支持C++11，所以安装了VS2015，但又因为没有卸载VS2008，导致VS2015安装不全）



2 安装boost

下载：http://www.boost.org/

我下载的是：boost_1_63_0

下载到C:\local目录下

1）首先：



双击bootstrap.bat

生成了

2）再在命令中输入b2 install

3）利用b2编译库文件
 --     b2 -a --with-python address-model=64 toolset=msvc runtime-link=static

之前你cmake下载的64位这里写64，如果是32位的就把之前的64改成32

4）设置变量
    --     set BOOST_ROOT=C:\local\boost_1_63_0
    --     set BOOST_LIBRARYDIR=C:\local\boost_1_63_0\stage\lib



参考：http://blog.csdn.net/worrydog/article/details/53947214



3安装dlib

下载：http://dlib.net/

直接输入python setup.py install 

```
pip install sklearn
```