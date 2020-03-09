
## SSH key

```
# useradd -m -d /home/tidb tidb

# passwd tidb

# visudo
tidb ALL=(ALL) NOPASSWD: ALL

# su - tidb

# ssh-keygen -t rsa
...ENTER(always)
```

copy auth_key to targe machine. E.g. machine ip is 172.16.10.61.

```

ssh-copy-id -i ~/.ssh/id_rsa.pub 172.16.10.61
```