
###MySQL setting
```sql
SET GLOBAL storage_engine = 'InnoDB';
CREATE DATABASE gogs CHARACTER SET utf8 COLLATE utf8_bin;
GRANT ALL PRIVILEGES ON gogs.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
QUIT;
```

### Install Gogs
```shell script
wget https://dl.gogs.io/0.11.34/linux_amd64.tar.gz
tar -xzf linux_amd64.tar.gz -C /
```
### Create Git User

```shell script
adduser git  
passwd git  
groupadd git
usermod -G git git

```

```shell script

mv /gogs /home/git/gogs
```


```shell script
cd /lib/systemd/system

vim gogs.service
```

```shell script

[Unit]
Description=Gogs
After=syslog.target
After=network.target
After=mysqld.service

[Service]
# Modify these two values and uncomment them if you have
# repos with lots of files and get an HTTP error 500 because
# of that
###
#LimitMEMLOCK=infinity
#LimitNOFILE=65535
Type=simple
User=git
Group=git
WorkingDirectory=/home/git/gogs
ExecStart=/home/git/gogs/gogs web
Restart=always
Environment=USER=git HOME=/home/git

[Install]
WantedBy=multi-user.target
```

```shell script
systemctl enable gogs.service
```

```shell script
systemctl start gogs
```


View the website: http://127.0.0.1:3000