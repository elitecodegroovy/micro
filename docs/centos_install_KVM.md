
## 设置虚拟化
```
lscpu | grep Virtualization
Virtualization: VT-x
```

## 安装kvm
```
yum install -y qemu-kvm libvirt libvirt-python libguestfs-tools virt-install

```
Start the libvirtd service:
```
systemctl enable libvirtd

```

验证安装：
```
lsmod | grep -i kvm

```

## 配置网桥网络
```
cp ifcfg-eno49 ifcfg-br0

```

Edit the Interface file and set followings:
```
root@linuxtechi network-scripts]# vi ifcfg-eno49
TYPE=Ethernet
BOOTPROTO=static
DEVICE=eno49
ONBOOT=yes
BRIDGE=br0

```
Edit the Bridge file (ifcfg-br0) and set the followings:
```
[root@linuxtechi network-scripts]# vi ifcfg-br0
TYPE=Bridge
BOOTPROTO=static
DEVICE=br0
ONBOOT=yes
IPADDR=192.168.10.21
NETMASK=255.255.255.0
GATEWAY=192.168.10.1
DNS1=192.168.10.11
```

Restart the network Service to enable the bridge interface.

```
systemctl restart network
```

Check the Bridge interface using below command :

```
[root@linuxtechi ~]# ip addr show br0
```
## 安装kvm

```
virt-install \
--virt-type=kvm \
--name=centos7 \
--vcpus=2 \
--memory=4096 \
--location=/tmp/CentOS-7-x86_64-Minimal-1804.iso \
--disk path=/home/vm/centos/centos7.qcow2,size=40,format=qcow2 \
--network bridge=br0 \
--graphics none \
--extra-args='console=ttyS0' \
--force

```




## virsh 命令

```
virsh list                 # 查看在运行的虚拟机
virsh dumpxml vm-name      # 查看kvm虚拟机配置文件
virsh start vm-name        # 启动kvm虚拟机
virsh shutdown vm-name     # 正常关机

virsh destroy vm-name      # 非正常关机，强制关闭虚拟机（相当于物理机直接拔掉电源）
virsh undefine vm-name     # 删除vm的配置文件

ls  /etc/libvirt/qemu
# 查看删除结果，Centos-6.6的配置文件被删除，但磁盘文件不会被删除

virsh define file-name.xml # 根据配置文件定义虚拟机
virsh suspend vm-name      # 挂起，终止
virsh resumed vm-name      # 恢复被挂起的虚拟机
virsh autostart vm-name    # 开机自启动vm
virsh start vm-name -console #启动vm，并且connect
virsh console <虚拟机名称>   # 连接虚拟机

```
