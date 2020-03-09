
QQ
```
SET PASSWORD FOR 'toonan'@'%' = PASSWORD('Toonan123$');
```

修改当前登录用户密码：

```
mysql>SET PASSWORD = PASSWORD('mypass');

```



##支持MySQL8.0

mysql_native_password 密码设置：

```

CREATE USER 'gaxt'@'%' IDENTIFIED BY "gaxt369#TGoodG";


ALTER USER 'gaxt'@'%' IDENTIFIED WITH mysql_native_password BY 'gaxt369#TGoodG';
grant select  on nsmjpt2.bs_bindgis  to backup@'%' ;
```


撤销权限：

```
revoke all on nsmjpt2.* from gaxt@'%' ;

```