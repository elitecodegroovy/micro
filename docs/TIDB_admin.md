
### 创建新用新用户

- CREATE USER 'finley'@'%' IDENTIFIED BY 'some_pass';
- GRANT ALL PRIVILEGES ON *.* TO 'finley'@'%' WITH GRANT OPTION;

### 修改已经存在的密码

SET PASSWORD FOR 'root'@'%' = 'xxx';

或者

ALTER USER 'jeffrey'@'localhost' IDENTIFIED BY 'mypass';