又遇到了Oracle数据库序列的问题，索性来个全库的备份和恢复。如下

imp/exp 方式

表模式备份： ­

oracle@sencloudServer: exp dhoffice/dhoffice rows=y indexes=n compress=n buffer=65536 file=dhoffice_090101.dmp log=exp_tianle_090101.log tables=david.tianle; ­

用户模式备份： ­

oracle@sencloudServer: exp dhoffice/dhoffice owner=david rows=y indexes=n compress=n buffer=65536 file=dhoffice__090101.dmp log=exp_david_090101.log; ­

完全模式备份： ­

oracle@sencloudServer: exp dhoffice/dhoffice rows=y indexes=n compress=n buffer=65536 full=y file=dhoffice_fulldatabase_090101.dmp log=exp_fulldatabase_090101.log; ­

先建立data_dir目录，在oracle中建立：

```
expdp  \'sys/PY2017_db AS SYSDBA\' directory=data_dir dumpfile=pyuser_20190926.dmp logfile=pyuser_20190926.log full=y  parallel=2 compression=all
```
表模式恢复：

oracle@sencloudServer: imp dhoffice/dhoffice fromuser=david touser=david rows=y indexes=n commit=y buffer=65536 file=exp_tianle_090101.dmp log=imp_tianle_090101.log; ­

用户模式恢复： ­

oracle@sencloudServer: imp dhoffice/dhoffice fromuser=david touser=david rows=y indexes=n commit=y buffer=65536 file=exp_tianle_090101.dmp log=exp_tianle_090101.log; ­

全库模式恢复： ­

oracle@sencloudServer: imp dhoffice/dhoffice rows=y indexes=n commit=y full=y ignore=y buffer=65536 file=/tmp/exp_fulldatabase_090101.dmp log=/tmp/imp.log;­


前期准备(注意执行的用户,root or oracle)
1.新建Oracle数据库备份存放目录

[root@637c4d4d4af4 /]# mkdir –p /home/oracle/backup 
2.修改目录权限（Oracle在oracle用户下运行）

[root@637c4d4d4af4 /]# chown oracle:oinstall /home/oracle/backup
创建目录对象
# 切换到oracle用户
[root@637c4d4d4af4 /]# su oracle

# 切换到SQLPLUS环境
[oracle@637c4d4d4af4 ~]$ sqlplus /nolog

# 管理员身份登录
SQL> conn /as sysdba

# 制定EXPDP输出目录
SQL> create directory expdp_bak_dir as '/home/oracle/backup';

# 授予system权限
SQL> grant create any directory to system;
全量备份
前期工作做好就可以开始备份了，下面是针对xm_uap这个用户进行全量备份的
# expdp_bak_dir 为上面创建的一个备份目录
# schemas 为要备份的用户对象
# compression=all 表示所有数据都进行压缩，可以大量减少备份文件大小
[oracle@637c4d4d4af4 ~]$ expdp system@helowin directory=expdp_bak_dir schemas=xm_uap dumpfile=backup.dmp logfile=backup.log compression=all
生成的备份文件就在expdp_bak_dir指向的目录下，backup.dmp为数据备份文件，backup.log为备份过程中的日志
全量恢复
恢复时记得先创建好对应的表空间,用户不需要创建
[oracle@637c4d4d4af4 ~]$ impdp system directory=expdp_bak_dir dumpfile=backup.dmp logfile=recover.log



