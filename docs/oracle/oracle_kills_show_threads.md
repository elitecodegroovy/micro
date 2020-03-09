

杀死所有oracle进程
 ```
$ ps -ef |grep $ORACLE_SID|grep -v grep|awk '{print $2}' | xargs kill -9
$ ipcs -m | grep oracle | awk '{print $2}' | xargs ipcrm shm

```


修改最大线程数

```
SQL> show parameter sessions
```


修改最大连接数
```
SQL> show parameter processe

```

修改最大连接数
```

SQL> alter system set processes=500 scope=spfile;

````