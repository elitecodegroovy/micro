
##### Install re2c

##### Build re2c

Download the source code from the site https://codeload.github.com/skvadrik/re2c/tar.gz/1.3.

```shell script
wget https://codeload.github.com/skvadrik/re2c/tar.gz/1.3
tar xvf re2c-1.3.tar.gz
cd re2c-1.3

./autogen.sh
./configure
make
make install
```


##### Install Ninja

Python
```
./configure.py --bootstrap
```

This will generate the ninja binary and a build.ninja file you can now use to build Ninja with itself.

CMake
````shell script
cmake -Bbuild-cmake -H.
cmake --build build-cmake
````


Copy the binary to the PATH setting.

```shell script
cp ninja /usr/bin/
```

 ClickHouse source code :
 ```shell script
git clone --recursive git@github.com:ClickHouse/ClickHouse.git
# or: git clone --recursive https://github.com/ClickHouse/ClickHouse.git

cd ClickHouse
```

```
mkdir build && cd build
cmake ..
ninja
cd ..
```

若要创建一个执行文件， 执行 ninja clickhouse。 这个命令会使得 
dbms/programs/clickhouse 文件可执行，您可以使用 client or server 参数运行。