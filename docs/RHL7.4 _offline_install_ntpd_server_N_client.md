

# RHL 7.5 Installs NTP Server

Install rpm files.
```
rpm -ivh autogen-libopts-5.18-5.el7.x86_64.rpm ntp-4.2.6p5-25.el7_3.2.ns7.01.x86_64.rpm

```

Configure the file.




## client NTP

Install rpm files.
```
rpm -ivh autogen-libopts-5.18-5.el7.x86_64.rpm ntp-4.2.6p5-25.el7_3.2.ns7.01.x86_64.rpm

```

Configure the file.

```
# vi /etc/ntp.conf

restrict -6 default kod nomodify notrap nopeer noquery
restrict 127.0.0.1
restrict -6 ::1
server 172.29.205.21


```

Start the client server:

```
systemctl start ntpd

ntpdate -u 172.29.205.21
systemctl enable ntpd
```

Verifiy the configuration:

```

#ntpq -p
remote           refid      st t when poll reach   delay   offset  jitter
==============================================================================
*172.29.205.21   211.19.59.28     3 u   27   64   17    0.147    1.718   5.659

```

It succeeds.