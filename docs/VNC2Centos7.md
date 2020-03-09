## GNOME

```
yum groupinstall "GNOME Desktop" "Graphical Administration Tools" -y
yum groupinstall "X Window System" "Desktop" -y
yum install tigervnc tigervnc-server -y
```
 
Tigervnc server
Tigervnc - is a high-performance, platform-neutral implementation of VNC

   #yum install tigervnc tigervnc-server -y
Lets copy example config and edit it:

    #cp /lib/systemd/system/vncserver@.service /etc/systemd/system/vncserver@:1.service

    #vi /etc/systemd/system/vncserver@:1.service
Example of config:

[Unit]
    Description=Remote desktop service (VNC)
    After=syslog.target network.target
    [Service]
    Type=forking
    ExecStartPre=/bin/sh -c '/usr/bin/vncserver -kill %i > /dev/null 2>&1 || :'
    ExecStart=/usr/sbin/runuser -l YOUR_NAME -c "/usr/bin/vncserver
    PIDFile=/home/YOUR_NAME/.vnc/%H%i.pid
    ExecStop=/bin/sh -c '/usr/bin/vncserver -kill %i > /dev/null 2>&1 || :'
    [Install]
    WantedBy=multi-user.target
Now enable autostart service and restart it:
    systemctl enable vncserver@:1.service
    systemctl restart vncserver@:1.service


configure: 

```
>vncserver
```
 
Reset passwd code :

```
#>vncpasswd

```

And finally connect to it:

    #vncviewer 192.168.1.236:5901