
# 安装telnet-server
yum -y install telnet-server

# 启动并设置开机自启动
systemctl start telnet.socket && systemctl enable telnet.socket

# 如果有防火墙，则需要放行23端口
firewall-cmd --zone=public --add-port=18081/tcp --permanent
firewall-cmd --reload
# 在使用 telnet 连接服务器时，默认是不允许使用root登陆的，因此需要创建一个普通用户并赋予sudo权限

# 添加普通用户并设置密码
useradd ljg
echo Ljg@20230705 | passwd --stdin ljg

# 增加huge账号的sudo权限
# 在配置文件 /etc/sudoers 中添加配置，但该文件默认是没有写权限的，因此需要先增加写权限
chmod u+w /etc/sudoers

vi /etc/sudoers
ljg     ALL=(ALL)       ALL


# 上面配置完成后就可以在windows下的终端中使用telnet命令来测试连接
telnet 192.168.0.31 23

## Upgrade SSH

wget https://cdn.openbsd.org/pub/OpenBSD/OpenSSH/portable/openssh-9.0p1.tar.gz
sudo systemctl stop sshd
sudo cp -a /etc/ssh /etc/ssh.bak
sudo cp -a /usr/sbin/sshd /usr/sbin/sshd.bak
sudo cp -a /usr/bin/ssh /usr/bin/ssh.bak
sudo rpm -qa | grep openssh
sudo rpm -e `rpm -qa | grep openssh` --nodeps

sudo yum install -y gcc gcc-c++ glibc make automake autoconf zlib zlib-devel pcre-devel  perl perl-Test-Simple openssl-devel
tar -zxvf openssh-9.0p1.tar.gz
cd openssh-9.0p1

./configure --prefix=/usr/local/openssh --with-ssl-dir=/usr/local/openssl --with-zlib
sudo make && sudo make install

# configure: error: OpenSSL library not found.
# sudo ln -s /usr/local/lib64/libssl.so.1.1 /usr/lib64/libssl.so.1.1
# sudo ln -s /usr/local/lib64/libcrypto.so.1.1 /usr/lib64/libcrypto.so.1.1
# CCFLAGS="-I/usr/local/include" LDFLAGS="-L/usr/local/lib" ./configure
# 或者 CCFLAGS="-I/usr/local/include" LDFLAGS="-L/usr/local/lib64" ./configure
注意： /usr/local/lib 根据系统版本选择 /usr/local/lib或/usr/local/lib64


sudo cp contrib/redhat/sshd.init /etc/init.d/sshd
sudo ln -s /usr/local/openssh/etc /etc/ssh
sudo ln -s /usr/local/openssh/sbin/sshd /usr/sbin/
sudo ln -s /usr/local/openssh/bin/* /usr/bin/
sudo systemctl daemon-reload
sudo systemctl start sshd && sudo systemctl enable sshd
sudo systemctl status sshd
# 查看状态，已经是 running 状态了



ssh -V


# ssh的默认配置文件是禁止root用户远程登录的
# 若需要root用户远程登录，则按修改如下配置文件，然后重启ssh服务即可
sudo vi /etc/ssh/etc/sshd_config
PermitRootLogin yes

注意重启：

sudo systemctl restart sshd

# ssh升级完成，测试ssh连接是否正常（有条件可以重启服务再测试ssh是否连接正常）
# 测试正常后就可以删除 telnet-server 和对应的用户了，删除步骤如下
systemctl stop telnet.socket
systemctl disable telnet.socket
rpm -e telnet-server
userdel -r ljg
# 删除用户
# 编辑sudo配置文件，去除对应配置段
vi /etc/sudoers

# 恢复sudo配置文件的权限
chmod u-w /etc/sudoers


ALTER SYSTEM ADD BACKEND "172.16.27.34:9050";


## change root account
``` 
SET PASSWORD FOR 'root' = PASSWORD('Ggjsptb@2023!');
```

