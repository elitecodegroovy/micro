# Centos 7 install standalone redis 6 from source code


## install dependencies
``` 
yum install gcc make wget tcl
```

## Install Redis

``` 
cd /usr/local/src

wget http://download.redis.io/releases/redis-6.2.6.tar.gz

tar xf redis-6.2.6.tar.gz
cd redis-6.2.6
make
cd src
make test
make install
```

## Configure Redis
The following commands will create a group and user named “redis”:
``` 
groupadd redis
adduser --system -g redis --no-create-home redis
```

Next create the /var/lib/redis directory which will store the redis databases:
``` 
mkdir -p /var/lib/redis
```
Give the redis user and group ownership over the directory:
``` 
chown redis: /var/lib/redis
```

Set the correct permissions with:
``` 
chmod 770 /var/lib/redis
```

Now that Redis is installed and intial setup has been completed on your CentOS 7 VPS, you’ll need to create the configuration file.

First, create the /etc/redis directory with:

``` 
mkdir /etc/redis
```

Next copy the default configuration file by running:
``` 
cp /usr/local/src/redis-6.2.6/redis.conf /etc/redis/
```
Open the file with your preferred text editor (we’ll be using nano):
``` 
vi /etc/redis/redis.conf
```
Search for the supervised directive and change it to “systemd”:

``` 
# If you run Redis from upstart or systemd, Redis can interact with your
# supervision tree. Options:
#   supervised no      - no supervision interaction
#   supervised upstart - signal upstart by putting Redis into SIGSTOP mode
#   supervised systemd - signal systemd by writing READY=1 to $NOTIFY_SOCKET
#   supervised auto    - detect upstart or systemd method based on
#                        UPSTART_JOB or NOTIFY_SOCKET environment variables
# Note: these supervision methods only signal "process is ready."
#       They do not enable continuous liveness pings back to your supervisor.

supervised auto
```
Search for the dir directive and set it to /var/lib/redis:

``` 
# The working directory.
#
# The DB will be written inside this directory, with the filename specified
# above using the 'dbfilename' configuration directive.
#
# The Append Only File will also be created inside this directory.
#
# Note that you must specify a directory here, not a file name.
dir /var/lib/redis
```

``` 
bind 0.0.0.0
```
Once done, save and close the file.

## Creating a Redis Systemd Unit File
Now that Redis is installed and configured, the last step is to create a systemd unit file so that you can manage the Redis servis as a regular systemd service.

Open your text editor and create the follwing file:

``` 
vi /etc/systemd/system/redis.service
```

Then add the following content to it:

``` 
[Unit]
Description=Redis Data Store
After=network.target

[Service]
User=redis
Group=redis
ExecStart=/usr/local/bin/redis-server /etc/redis/redis.conf
ExecStop=/usr/local/bin/redis-cli shutdown
Restart=always

[Install]
WantedBy=multi-user.target
```

Save and close the file. Enable and start the Redis service with:
``` 
systemctl enable redis
systemctl start redis
```
To check that the Redis service has no errors, run:
``` 
systemctl status redis
```

other setting for system (添加到文件/etc/sysctl.conf):
``` 
net.core.somaxconn=1024
vm.overcommit_memory=1
fs.file-max=102400
```

执行：
``` 
sysctl -p
```
