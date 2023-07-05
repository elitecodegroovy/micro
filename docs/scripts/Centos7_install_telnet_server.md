
# 安装telnet-server
yum -y install telnet-server

# 启动并设置开机自启动
systemctl start telnet.socket && systemctl enable telnet.socket

# 如果有防火墙，则需要放行23端口
firewall-cmd --zone=public --add-port=23/tcp --permanent

# 在使用 telnet 连接服务器时，默认是不允许使用root登陆的，因此需要创建一个普通用户并赋予sudo权限

# 添加普通用户并设置密码
useradd ljg
echo Ljg@20230705 | passwd --stdin ljg

# 增加huge账号的sudo权限
# 在配置文件 /etc/sudoers 中添加配置，但该文件默认是没有写权限的，因此需要先增加写权限
chmod u+w /etc/sudoers

vi /etc/sudoers
ljg ALL=(ALL) ALL

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

sudo yum install -y gcc gcc-c++ glibc make automake autoconf zlib zlib-devel pcre-devel  perl perl-Test-Simple

tar -zxvf openssh-9.0p1.tar.gz
cd openssh-9.0p1

./configure --prefix=/usr/local/openssh --with-ssl-dir=/usr/local/openssl --with-zlib
sudo make && sudo make install


sudo cp contrib/redhat/sshd.init /etc/init.d/sshd

sudo ln -s /usr/local/openssh/etc /etc/ssh
sudo ln -s /usr/local/openssh/sbin/sshd /usr/sbin/
sudo ln -s /usr/local/openssh/bin/* /usr/bin/

sudo systemctl daemon-reload

sudo systemctl start sshd && sudo systemctl enable sshd
# 查看状态，已经是 running 状态了
# sudo systemctl status sshd

ssh -V

# ssh的默认配置文件是禁止root用户远程登录的
# 若需要root用户远程登录，则按修改如下配置文件，然后重启ssh服务即可
sudo vi /etc/ssh/etc/sshd_config
PermitRootLogin yes


sudo systemctl restart sshd
