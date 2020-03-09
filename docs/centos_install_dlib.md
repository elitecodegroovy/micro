We need g++ cmake , gcc and gcc-c++ to compile dlib and boost

```
sudo yum install python36* -y
sudo yum install gcc g++ gcc gcc-c++ cmake
```

download boost and dlib
```
wget https://dl.bintray.com/boostorg/release/1.65.1/source/boost_1_65_1.tar.gz
wget http://dlib.net/files/dlib-19.7.tar.bz2

```
extract the packages :

```
tar -xvf boost_1_65_1.tar.gz
tar -xvf dlib-19.7.tar.bz2
```

remove the compressed files :

```
sudo rm boost_1_65_1.tar.gz
sudo rm dlib-19.7.tar.bz2

```
Let’s begin with boost.

go in the boost folder and configure run bootstrap with following command. change it depending your python version and path.

```
cd ../boost_1_65_1/
sudo ./bootstrap.sh — with-python=python — with-libraries=python — prefix=/usr

#pthon3
./bootstrap.sh — with-python=python3.6 — with-libraries=python — prefix=/usr

```
We will need to edit the boost configuration to reflect your system and python. the final file should look like this.
```
import feature ;
# Compiler configuration. This definition will be used unless
# you already have defined some toolsets in your user-config.jam
# file.
if ! gcc in [ feature.values <toolset> ]
{
 using gcc ;
}
project : default-build <toolset>gcc ;
# Python configuration
import python ;
if ! [ python.configured ]
{
 #using python : 3.6 : /usr/bin/python3.6 ;
using python : 3.6 : /usr/bin/python3.6 : /usr/include/python3.6m : /usr/lib64 ;
}
# List of — with-<library> and — without-<library>
# options. If left empty, all libraries will be built.
# Options specified on the command line completely
# override this variable.
libraries = — with-python ;
# These settings are equivivalent to corresponding command-line
# options.
option.set prefix : /usr ;
option.set exec-prefix : /usr ;
option.set libdir : /usr/lib64 ;
option.set includedir : /usr/include ;
# Stop on first error
option.set keep-going : false ;
```
Then run the installer

```
sudo ./b2

```
export this values :

```
export BOOST_INCLUDEDIR=/home/ec2-user/boost_1_65_1
sudo -E python3 setup.py install
export BOOST_ROOT=/home/ec2-user/boost_1_65_1
sudo -E python3 setup.py install
export BOOST_LIBRARYDIR=/home/ec2-user/boost_1_65_1/stage/lib

```
Now let’s go with dlib

```
python3 setup.py install
```

you only need the python api.

### install pip3

 ```
cd Python-3.6.3
./configure --enable-optimizations
make
make install
 ```
 
 pip3 install 
```

pip3 install --upgrade pip setuptools wheel
```
