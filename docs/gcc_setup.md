
## Install gcc9.3.0

### Download File

Download the source code from the site https://ftp.gnu.org/gnu/gcc/gcc-9.3.0/gcc-9.3.0.tar.gz.
```
wget https://ftp.gnu.org/gnu/gcc/gcc-9.3.0/gcc-9.3.0.tar.gz
tar xvf gcc-9.3.0.tar.gz
cd gcc-9.3.0

```

### Install Dependencies

```shell script
yum install gmp-devel mpfr-devel libmpc-devel
```

### Build

```shell script
# It will take for an half of one hour.
mkdir build
cd build
../build/configure --enable-languages=c,c++ --disable-multilib
make -j$(nproc) 
make install
```

### Post-installation

You should add the install dir of GCC to your PATH and LD_LIBRARY_PATH in order to use the newer GCC.
 Add the following settings to /etc/profile:
 
```shell script
export PATH=/usr/local/bin:$PATH
export LD_LIBRARY_PATH=/usr/local/lib64:$LD_LIBRARY_PATH

```

Enable the profile.
```shell script
source /etc/profile
```