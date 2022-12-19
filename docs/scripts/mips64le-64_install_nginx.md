
# install

```
FROM mips64le/nginx:1.22.1
COPY data /usr/share/nginx/

```

run the image:
``` 
sudo docker run -p 80:80 -p 443:443 \
-v /home/tisson/nginx/data:/usr/share/nginx/data  \
-v /home/tisson/nginx/conf/nginx.conf:etc/nginx/nginx.conf:ro  \
--restart=always --name nginx1.22.1 -d \
nginx:1.22.1
```


