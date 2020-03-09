

ogg热备策略：
```
alter database add supplemental log data (primary key) columns;
alter database add supplemental log data (unique) columns;
alter system set DB_RECOVERY_FILE_DEST_SIZE=100g;
```

关闭archive

```

SQL>shutdown immediate;

SQL> startup mount;
SQL> alter database noarchivelog;



SQL> archive log list;

SQL>alter database open;


```