
## Docker install mysql

```
docker pull centos/mysql-57-centos7:5.7
```

创建mysql：mysql的用户和用户组：

```
useradd -g mysql mysql

##
mdkir -p /opt/docker_mysql/data
```
需要设置数据文件目录的权限，假设目录文件为/opt/docker_mysql/data，
设置对应权限的命令：

```
chown -R mysql:mysql /opt/docker_mysql 
```

启动一个容器：
```
docker run -d --name mysql_database  \
-e MYSQL_ROOT_PASSWORD=TOONAN123456  \
-v /opt/docker_mysql/data:/var/lib/mysql/data \
-p 3309:3306 centos/mysql-57-centos7:5.7
```

生产环境下，可以添加启动加载项`--restart=always `。 


```
mysql -u root -p

```
创建数据库的gc用户
```
CREATE USER 'gc'@'%' IDENTIFIED BY '019#^(963';

GRANT ALL PRIVILEGES ON *.* TO 'gc'@'%' identified by '019#^(963';

```

## Environment variables and volumes
The image recognizes the following environment variables that you can set during initialization by passing `-e VAR=VALUE` to the Docker run command.

> MYSQL_USER

User name for MySQL account to be created

> MYSQL_PASSWORD

Password for the user account

>MYSQL_DATABASE

Database name

>MYSQL_ROOT_PASSWORD

Password for the root user (optional)

The following environment variables influence the MySQL configuration file. They are all optional.

>MYSQL_LOWER_CASE_TABLE_NAMES (default: 0)

Sets how the table names are stored and compared

>MYSQL_MAX_CONNECTIONS (default: 151)

The maximum permitted number of simultaneous client connections

>MYSQL_MAX_ALLOWED_PACKET (default: 200M)

The maximum size of one packet or any generated/intermediate string

>MYSQL_FT_MIN_WORD_LEN (default: 4)

The minimum length of the word to be included in a FULLTEXT index

>MYSQL_FT_MAX_WORD_LEN (default: 20)

The maximum length of the word to be included in a FULLTEXT index

>MYSQL_AIO (default: 1)

Controls the innodb_use_native_aio setting value in case the native AIO is broken. See http://help.directadmin.com/item.php?id=529

>MYSQL_TABLE_OPEN_CACHE (default: 400)

The number of open tables for all threads

>MYSQL_KEY_BUFFER_SIZE (default: 32M or 10% of available memory)

The size of the buffer used for index blocks

>MYSQL_SORT_BUFFER_SIZE (default: 256K)

The size of the buffer used for sorting

>MYSQL_READ_BUFFER_SIZE (default: 8M or 5% of available memory)

The size of the buffer used for a sequential scan

>MYSQL_INNODB_BUFFER_POOL_SIZE (default: 32M or 50% of available memory)

The size of the buffer pool where InnoDB caches table and index data

>MYSQL_INNODB_LOG_FILE_SIZE (default: 8M or 15% of available available)

The size of each log file in a log group

>MYSQL_INNODB_LOG_BUFFER_SIZE (default: 8M or 15% of available memory)
The size of the buffer that InnoDB uses to write to the log files on disk

>MYSQL_DEFAULTS_FILE (default: /etc/my.cnf)

Point to an alternative configuration file

>MYSQL_BINLOG_FORMAT (default: statement)

Set sets the binlog format, supported values are row and statement

>MYSQL_LOG_QUERIES_ENABLED (default: 0)

To enable query logging set this to 1

You can also set the following mount points by passing the -v /host:/container flag to Docker.

    /var/lib/mysql/data
MySQL data directory

Notice: When mouting a directory from the host into the container, ensure that the mounted directory has the appropriate permissions and that the owner and group of the directory matches the user UID or name which is running inside the container.

MySQL auto-tuning
When the MySQL image is run with the --memory parameter set and you didn't specify value for some parameters, their values will be automatically calculated based on the available memory.

>MYSQL_KEY_BUFFER_SIZE (default: 10%)

key_buffer_size

>MYSQL_READ_BUFFER_SIZE (default: 5%)

read_buffer_size

>MYSQL_INNODB_BUFFER_POOL_SIZE (default: 50%)

innodb_buffer_pool_size

>MYSQL_INNODB_LOG_FILE_SIZE (default: 15%)

innodb_log_file_size

>MYSQL_INNODB_LOG_BUFFER_SIZE (default: 15%)

innodb_log_buffer_size

MySQL root user
The root user has no password set by default, only allowing local connections. You can set it by setting the MYSQL_ROOT_PASSWORD environment variable. This will allow you to login to the root account remotely. Local connections will still not require a password.

To disable remote root access, simply unset `MYSQL_ROOT_PASSWORD `and restart the container.

Changing passwords
Since passwords are part of the image configuration, the only supported method to change passwords for the database user (MYSQL_USER) and root user is by changing the environment variables MYSQL_PASSWORD and MYSQL_ROOT_PASSWORD, respectively.

Changing database passwords through SQL statements or any way other than through the environment variables aforementioned will cause a mismatch between the values stored in the variables and the actual passwords. Whenever a database container starts it will reset the passwords to the values stored in the environment variables.

Default my.cnf file
With environment variables we are able to customize a lot of different parameters or configurations for the mysql bootstrap configurations. If you'd prefer to use your own configuration file, you can override the MYSQL_DEFAULTS_FILE env variable with the full path of the file you wish to use. For example, the default location is `/etc/my.cnf` but you can change it to `/etc/mysql/my.cnf` by setting MYSQL_DEFAULTS_FILE=/etc/mysql/my.cnf

Extending image
This image can be extended using source-to-image.

For example, to build a customized MariaDB database image my-mysql-rhel7 with a configuration in ~/image-configuration/ run:
```
$ s2i build ~/image-configuration/ rhscl/mysql-57-rhel7 my-mysql-rhel7
```

The directory passed to s2i build can contain these directories:

`mysql-cfg/` When starting the container, files from this directory will be used as a configuration for the mysqld daemon. envsubst command is run on this file to still allow customization of the image using environmental variables

`mysql-pre-init/` Shell scripts (*.sh) available in this directory are sourced before mysqld daemon is started.

`mysql-init/` Shell scripts (*.sh) available in this directory are sourced when mysqld daemon is started locally. In this phase, use ${mysql_flags} to connect to the locally running daemon, for example mysql $mysql_flags < dump.sql

Variables that can be used in the scripts provided to s2i:

`$mysql_flags `arguments for the mysql tool that will connect to the locally running mysqld during initialization

`$MYSQL_RUNNING_AS_MASTER` variable defined when the container is run with run-mysqld-master command

`$MYSQL_RUNNING_AS_SLAVE` variable defined when the container is run with run-mysqld-slave command

`$MYSQL_DATADIR_FIRST_INIT` variable defined when the container was initialized from the empty data dir

During s2i build all provided files are copied into /opt/app-root/src directory into the resulting image. If some configuration files are present in the destination directory, files with the same name are overwritten. Also only one file with the same name can be used for customization and user provided files are preferred over default files in /usr/share/container-scripts/mysql/- so it is possible to overwrite them.

Same configuration directory structure can be used to customize the image every time the image is started using docker run. The directory has to be mounted into /opt/app-root/src/ in the image (-v ./image-configuration/:/opt/app-root/src/). This overwrites customization built into the image.

Securing the connection with SSL
In order to secure the connection with SSL, use the extending feature described above. In particular, put the SSL certificates into a separate directory:

```
sslapp/mysql-certs/server-cert-selfsigned.pem
sslapp/mysql-certs/server-key.pem
```
And then put a separate configuration file into mysql-cfg:

```
$> cat sslapp/mysql-cfg/ssl.cnf
[mysqld]
ssl-key=${APP_DATA}/mysql-certs/server-key.pem
ssl-cert=${APP_DATA}/mysql-certs/server-cert-selfsigned.pem
```
Such a directory sslapp can then be mounted into the container with -v, or a new container image can be built using s2i.

Upgrading and data directory version checking
MySQL and MariaDB use versions that consist of three numbers X.Y.Z (e.g. 5.6.23). For version changes in Z part, the server's binary data format stays compatible and thus no special upgrade procedure is needed. For upgrades from X.Y to X.Y+1, consider doing manual steps as described at https://dev.mysql.com/doc/refman/5.7/en/upgrading-from-previous-series.html.

Skipping versions like from X.Y to X.Y+2 or downgrading to lower version is not supported; the only exception is ugrading from MariaDB 5.5 to MariaDB 10.0.

Important: Upgrading to a new version is always risky and users are expected to make a full back-up of all data before.

A safer solution to upgrade is to dump all data using `mysqldump` or `mysqldbexport` and then load the data using mysql or `mysqldbimport` into an empty (freshly initialized) database.

Another way of proceeding with the upgrade is starting the new version of the mysqld daemon and run mysql_upgrade right after the start. This so called in-place upgrade is generally faster for large data directory, but only possible if upgrading from the very previous version, so skipping versions is not supported.


There are also some other actions that you may want to run at the beginning of the container start, when the local daemon is running, no matter what version of the data is detected:



- optimize -- runs mysqlcheck --optimize. It optimizes all the tables.

- analyze -- runs mysqlcheck --analyze. It analyzes all the tables.

- disable -- nothing is done regarding data directory version.

Multiple values are separated by comma and run in-order, e.g. MYSQL_DATADIR_ACTION="optimize,analyze".

Changing the replication binlog_format
Some applications may wish to use row binlog_formats (for example, those built with change-data-capture in mind). The default replication/binlog format is statement but to change it you can set the MYSQL_BINLOG_FORMAT environment variable. For example` MYSQL_BINLOG_FORMAT=row`. Now when you run the database with master replication turned on (ie, set the Docker/container cmd to be run-mysqld-master) the binlog will emit the actual data for the rows that change as opposed to the statements (ie, DML like insert...) that caused the change.

Troubleshooting
The mysqld deamon in the container logs to the standard output, so the log is available in the container log. The log can be examined by running:

```
docker logs <container>
```

docker run \
-p 3306:3306 \
--name mysql \
--privileged=true \
-v /Users/superpowerai/docker/mysql:/etc/mysql \
-v /Users/superpowerai/docker/mysql/logs:/logs \
-v /Users/superpowerai/docker/mysql/data:/var/lib/mysql \
-v /etc/localtime:/etc/localtime \
-e MYSQL_ROOT_PASSWORD=Root20230702mysql8 \
-d mysql:8.0.30 \
--default-authentication-plugin=mysql_native_password 