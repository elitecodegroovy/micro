
在《Rethinking the Inception Architecture for Computer Vision》中认为：基于inception v1进行结构的改进是inception v2；在inception v2上加上BN是inception v3；

##Install Python 3.5.5

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
 

在《Inception-v4, Inception-ResNet and the Impact of Residual Connections on Learning》
中将《Batch Normalization: Accelerating Deep Network Training by Reducing Internal Covariate Shift》
认为是inception v2（即inception v1 上进行小改动再加上BN）；
《Rethinking the Inception Architecture for Computer Vision》
认为是inception v3


```
sudo yum -y install update
sudo yum -y install epel-release
sudo yum -y install gcc gcc-c++ python-pip python-devel atlas atlas-devel gcc-gfortran openssl-devel libffi-devel
# install lastest pip version


pip install --upgrade pip
# use pip or pip3 as you prefer for python or python3
pip install --upgrade virtualenv
virtualenv --system-site-packages /opt/venvs/tensorflow
source /opt/venvs/tensorflow/bin/activate
#pip3 install --upgrade numpy scipy wheel cryptography #optional
pip3 install --upgrade tensorflow
```

## Install pip

```
curl "https://bootstrap.pypa.io/get-pip.py" -o "get-pip.py"

python get-pip.py
```

## gpu support(Centos 7)
```
yum install epel-release
yum install --enablerepo=epel dkms

sudo yum install kernel-devel-$(uname -r) kernel-headers-$(uname -r)

./***.run
./***1.run
./***2.run


tar xvf cudnn-9.0-linux-x64-v7.1.tgz
cp cuda/include/cudnn.h /usr/local/cuda-9.0/include
cp cuda/lib64/libcudnn* /usr/local/cuda-9.0/lib64
chmod a+r /usr/local/cuda-9.0/include/cudnn.h /usr/local/cuda-9.0/lib64/libcudnn*

./NVIDIA-Linux-x86_64-390.59.run
pip3 install --upgrade tensorflow-gpu

```


##set environment


```
export PATH=$PATH:/usr/local/cuda-9.0/bin
export LD_LIBRARY_PATH=/usr/local/cuda-9.0/lib64
```



