在centos中增加用户使用adduser命令，而创建用户组使用groupadd命令，这个是不是非常的方便呀，其实复杂点的就是用户的组与组权限的命令了，下面来给各位介绍一下吧。 
1、建用户：

adduser phpq                         //新建phpq用户
passwd phpq                          //给phpq用户设置密码

2、建工作组

groupadd test                        //新建test工作组

3、新建用户同时增加工作组

useradd -g test phpq               //新建phpq用户并增加到test工作组
注：-g 所属组 -d 家目录 -s 所用的SHELL

4、给已有的用户增加工作组

usermod -G groupname username
或者：gpasswd -a username groupname 

（注意：添加用户到某一个组 可以使用usermod -G groupname username这个命令可以添加一个用户到指定的组，但是以前添加的组就会清空掉。 
所以想要添加一个用户到一个组，同时保留以前添加的组时，请使用gpasswd这个命令来添加操作用户）

5、临时关闭

在/etc/shadow文件中属于该用户的行的第二个字段（密码）前面加上星号就可以了。想恢复该用户，去掉星号即可。

或者使用如下命令关闭用户账号：

passwd peter –l
重新释放：
passwd peter –u

6、永久性删除用户账号

userdel peter
groupdel peter
usermod –G peter peter   （强制删除该用户的主目录和主目录下的所有文件和子目录）
1
2
3
7、从组中删除用户

编辑/etc/group 找到GROUP1那一行，删除 A 或者用命令 gpasswd -d A GROUP

8、显示用户信息

id user
cat /etc/passwd
1
2
补充:查看用户和用户组的方法

用户列表文件：/etc/passwd
用户组列表文件：/etc/group

查看系统中有哪些用户：cut -d : -f 1 /etc/passwd
查看可以登录系统的用户：cat /etc/passwd | grep -v /sbin/nologin | cut -d : -f 1
查看某一用户：w 用户名
查看登录用户：who
查看用户登录历史记录：last