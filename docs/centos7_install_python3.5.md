

 ```
 ##env
 yum group install "Development Tools"
 yum install openssl-devel
 #
cd Python-3.5.1
./configure --enable-optimizations
make
make install

## pip3 upgrade
pip3 install --upgrade pip
 ```