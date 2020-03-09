
Download the current stable node version.
>https://nodejs.org/dist/v8.11.3/node-v8.11.3-linux-x64.tar.xz

Configure the file:

```
tar xvf  node-v8.11.3-linux-x64.tar.xz

mv node-v8.11.3-linux-x64 node
```

Configure the env:

```
export NODE_HOME=/opt/node

export PATH=$PATH:$NODE_HOME/bin
```

Enable the configuration variables:
```
. /etc/profile
```

