
## Usergroup

Create a usergroup

```
groupadd groupname
```

e.g.
```
groupadd app1
```

Create username with home directory.

```
sudo useradd username -m -s /sbin/nologin -d /home/username -g groupname

```

e.g.

```
useradd app1 -m -s /sbin/nologin -d /home/app1 -g app1

```

## account

在增加了-s /sbin/nologin 参数后，那么这个帐号就不能登陆了，如果想要恢复登陆使用

sudo usermod -s /bin/bash username
禁用用户登录权限

sudo usermod -s /sbin/nologin username