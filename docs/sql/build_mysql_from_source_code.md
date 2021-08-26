
1 Run the sql with root account

MySQL version: mysql-8.0.26-linux-glibc2.12-x86_64
``` 
# Licensed to Tisson. Write by pgy.
#!/usr/bin/env sh
source /etc/profile
source $HOME/.bash_profile

INSTPKG="$0"
INSTPKG_DIR=`dirname "$INSTPKG"`
MY_CNF=/etc/my.cnf
ROOT_PASSWD=6Rcf27.0BdulDnpNs8W+

if [ $USER != "root" ]
then
    echo -e "\033[5;31mERROR:\033[0m Please use \033[33mroot\033[0m account to execute the script."
    exit
fi
echo -e ""
echo -e "MYSQL_INSTALL=/opt If you want to change the installation location, you can change it."
read -p "Enter MYSQL_INSTALL: "  MYSQL_INSTALL

if [ -z "$MYSQL_INSTALL" ]
then
    MYSQL_INSTALL=/opt
fi

echo -e "MYSQL_INSTALL=$MYSQL_INSTALL"

CONFIRM_FLAG=1
while [ $CONFIRM_FLAG -eq 1 ]
do
    read -p "Do you want to install mysql-server on $MYSQL_INSTALL [Y/N]: "  YES_NO
    if  [ "$YES_NO" == "Y" ] || [ "$YES_NO" == "y" ] || [ "$YES_NO" == "N" ] || [ "$YES_NO" == "n" ]
    then
        CONFIRM_FLAG=0
    else
        CONFIRM_FLAG=1
    fi
done

if [ "$YES_NO" != "Y" ] && [ "$YES_NO" != "y" ]
then
    echo -e "\033[5;31mSorry,Bye!\033[0m"
    exit
fi

echo -e ""
echo -e "MYSQL_DATA_PATH=/opt/mysql_data If you want to change the installation location, you can change it."
read -p "Enter MYSQL_DATA_PATH: "  MYSQL_DATA_PATH

if [ -z "$MYSQL_DATA_PATH" ]
then
    MYSQL_DATA_PATH=/opt/mysql_data
fi

echo -e "MYSQL_DATA_PATH=$MYSQL_DATA_PATH"

CONFIRM_FLAG=1
while [ $CONFIRM_FLAG -eq 1 ]
do
    read -p "Do you want to set mysql data on $MYSQL_DATA_PATH [Y/N]: "  YES_NO
    if  [ "$YES_NO" == "Y" ] || [ "$YES_NO" == "y" ] || [ "$YES_NO" == "N" ] || [ "$YES_NO" == "n" ]
    then
        CONFIRM_FLAG=0
    else
        CONFIRM_FLAG=1
    fi
done
if [ "$YES_NO" != "Y" ] && [ "$YES_NO" != "y" ]
then
    echo -e "\033[5;31mSorry,Bye!\033[0m"
    echo ""
    exit
fi

echo -e ""
echo "====================================="
echo -e "\033[44;37mPlease confirm again whether it is installed in the following location.\033[0m"
echo -e "\033[44;37mMYSQL_INSTALL=$MYSQL_INSTALL\033[0m"
echo -e "\033[44;37mMYSQL_DATA_PATH=$MYSQL_DATA_PATH\033[0m"
echo "====================================="

CONFIRM_FLAG=1
while [ $CONFIRM_FLAG -eq 1 ]
do
    read -p "Do you want to continue? [Y/N]: "  YES_NO
    if  [ "$YES_NO" == "Y" ] || [ "$YES_NO" == "y" ] || [ "$YES_NO" == "N" ] || [ "$YES_NO" == "n" ]
    then
        CONFIRM_FLAG=0
    else
        CONFIRM_FLAG=1
    fi
done

if [ "$YES_NO" != "Y" ] && [ "$YES_NO" != "y" ]
then
    echo -e "\033[5;31mSorry,Bye!\033[0m"
    echo ""
    exit
fi


echo -e ""
export MYSQL_INSTALL=$MYSQL_INSTALL
export MYSQL_DATA_PATH=$MYSQL_DATA_PATH

if [ -d "$MYSQL_DATA_PATH" ]; then
  echo -e "\033[5;31mERROR:\033[0mThe path [\033[33m$MYSQL_DATA_PATH\033[0m] already exists, please check it ."
  echo ""
  exit
fi

mkdir -p $MYSQL_DATA_PATH
if [  -d  "$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64" ]
then
    echo -e "\033[5;31mSorry!\033[0m MySQL has been installed in $MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64"
    echo ""
    exit
fi



MYSGROUP=`cat /etc/group | grep -w mysql | awk -F':' '{print $1}'`

if [ -z $MYSGROUP ]
then
    groupadd mysql
fi

MYSGROUP=`cat /etc/passwd  | grep -w mysql | awk -F':' '{print $1}'`

if [ -z $MYSGROUP ]
then
   useradd -r -s /sbin/nologin -g mysql mysql
fi

tar xvzf $INSTPKG_DIR/mysql-8.0.26-linux-glibc2.12-x86_64.tar.gz -C  $MYSQL_INSTALL
echo "" 
echo "" > $MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/mysql.pid
chown -R mysql:mysql $MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64
tee $MY_CNF <<-'EOF'
[mysqld]
# 设置4417端口
port=4417
# 设置mysql的安装目录
basedir=$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64
# 设置mysql数据库的数据的存放目录
datadir=$MYSQL_DATA_PATH
# 允许最大连接数
max_connections=300
# 允许连接失败的次数。
max_connect_errors=10
# 服务端使用的字符集默认为utf8mb4
character-set-server=utf8
# 创建新表时将使用的默认存储引擎
default-storage-engine=INNODB
# 默认使用“mysql_native_password”插件认证
#mysql_native_password
default_authentication_plugin=mysql_native_password
#设置临时目录
tmpdir=/tmp
#设置socke文件所在目录
socket=$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/mysql.sock
#只能用IP地址检查客户端的登录，不用主机名
skip_name_resolve=1
#是否对sql语句大小写敏感，1表示不敏感
lower_case_table_names=1
log-error=$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/mysqld.log
#pid
pid-file=$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/mysql.pid
# 用户
user=mysql
# 允许访问的IP网段
bind-address=0.0.0.0
#binlog 失效时间
expire_logs_days=90

[mysql]
# 设置mysql客户端默认字符集
default-character-set=utf8
#设置socke文件所在目录
socket=$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/mysql.sock

[client]
# 设置mysql客户端连接服务端时默认使用的端口
port=4417
default-character-set=utf8
EOF

sed -i "s%\$MYSQL_INSTALL%$MYSQL_INSTALL%g"  $MY_CNF
sed -i "s%\$MYSQL_DATA_PATH%$MYSQL_DATA_PATH%g"  $MY_CNF

echo ""
echo "MySQL Configure: $MY_CNF"

if [ -f "/etc/rc.d/init.d/mysqld" ]
then
   rm -f /etc/rc.d/init.d/mysqld
fi

if [ ! -f "/usr/lib64/libtinfo.so.5" ]
then
  ln -s /usr/lib64/libtinfo.so.6.1 /usr/lib64/libtinfo.so.5
fi



cp $MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/support-files/mysql.server /etc/rc.d/init.d/mysqld
chown -R mysql:mysql $MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64
chkconfig --add mysqld
#chkconfig
chkconfig mysqld on
if [ -L "/usr/local/mysql" ]
then
   rm -f /usr/local/mysql
fi

ln -s $MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64  /usr/local/mysql


echo ""
export PATH=$PATH:$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/bin
echo -e "\033[44;37mInitial MySql database, It will take a few minutes, please wait. DO NOT CTRL+C.\033[0m"

mysqld --defaults-file=/etc/my.cnf  --initialize --console

_RUNOPT=`cat $MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/mysqld.log | grep 'A temporary password is generated for' | awk '{ print " -uroot -p\""$13"\""}'`

service mysqld start
#/bin/systemctl stop mysql.service
sleep 10
eval exec "\"$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/bin/mysql\" $_RUNOPT --connect-expired-password  -e \"ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '$ROOT_PASSWD';\nexit;\"  > /tmp/install_mysql.log 2>&1 &"
if [ $? -eq 0 ]; then
   sleep 1
   #echo "MySQL install successfully!"
   sleep 100
else
   echo "MySQL install failure!"
   exit
fi

echo "MySQL installed successfully. [default port: 4417]"
echo "==================================================="
echo "$MYSQL_INSTALL/mysql-8.0.26-linux-glibc2.12-x86_64/bin/mysql -uroot -p\"$ROOT_PASSWD\""
echo "==================================================="
#echo "================  MODIFY PASSWORD ================"
#echo "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'yourpassword';"



```