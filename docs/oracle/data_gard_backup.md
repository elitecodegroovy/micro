

强制记录日志
```
SQL> alter database force logging;
Database altered.

```

检查状态(YES为强制)
```
SQL> select name,force_logging from v$database;
NAME                        FORCE_LOG
--------------------------- ---------
SRMCLOUD                    YES

```

查询redo日志大小
```
SQL>select group#,bytes/1024/1024 as M from v$log;

    GROUP#	    M
---------- ----------
	 1	   50
	 2	   50
	 3	   50

```
这里三个是50M.


```
orcl226 =
   (DESCRIPTION =
     (ADDRESS_LIST =
       (ADDRESS = (PROTOCOL = TCP)(HOST = 192.168.1.226)(PORT = 1521))
     )
     (CONNECT_DATA =
       (SERVICE_NAME = orcl226)
     )
   )
orcl225 =
   (DESCRIPTION =
     (ADDRESS_LIST =
       (ADDRESS = (PROTOCOL = TCP)(HOST = 192.168.1.225)(PORT = 1521))
     )
     (CONNECT_DATA =
       (SERVICE_NAME = orcl225)
     )
   )
   
```