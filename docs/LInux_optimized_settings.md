
- 

- Configure open file number.

- Configure the CPUfreq governor mode.


### Operating system parameters

| Configuration | Description |
| :-- | :-------------------------- |
| Nice Limits | For system users, set the default value of `nice` in TiDB to `0` |
| min_free_kbytes | The setting for `vm.min_free_kbytes` in `sysctl.conf` needs to be high enough |
| User Open Files Limit | For database administrators, set the number of TiDB open files to `1000000` |
| System Open File Limits | Set the number of system open files to `1000000` |
| User Process Limits | For TiDB users, set the `nproc` value to `4096` in `limits.conf` |
| Address Space Limits | For TiDB users, set the space to `unlimited` in `limits.conf` |
| File Size Limits | For TiDB users, set the `fsize` value to `unlimited` in `limits.conf` |
| Disk Readahead | Set the value of the `readahead` data disk to `4096` at a minimum |
| NTP service | Configure the NTP time synchronization service for each node |
| SELinux  | Turn off the SELinux service for each node |
| CPU Frequency Scaling | It is recommended to turn on CPU overclocking |
| Transparent Hugepages | For Red Hat 7+ and CentOS 7+ systems, it is required to set the Transparent Hugepages to `always` |
| I/O Scheduler | Set the I/O Scheduler of data disks to the `deadline` mode |
| vm.swappiness | Set `vm.swappiness = 0` in `sysctl.conf` |
| net.core.somaxconn | Set `net.core.somaxconn = 32768` in `sysctl.conf` |
| net.ipv4.tcp_syncookies | Set `net.ipv4.tcp_syncookies = 0` in `sysctl.conf` |


```
#vi /etc/sysctl.conf
fs.file-max = 202400
fs.nr_open = 102400
vm.swappiness = 0
net.core.somaxconn = 32768
net.ipv4.tcp_syncookies = 0

```

