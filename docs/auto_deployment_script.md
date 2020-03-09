
```
scp /var/jenkins_home/workspace/gc/gc-api/target/gc-api-1.0.0-SNAPSHOT.jar tidb@192.168.1.226:/opt/apps/gc/ && ssh tidb@192.168.1.226 'cd /opt/apps/gc; sh startup.sh'

```

startup.sh

```
# -----------------------------------------------------------------------------
# Start Script for the  Server
# -----------------------------------------------------------------------------
# resolve links - $0 may be a softlink

SCRIPT_PATH="${BASH_SOURCE[0]}"
SCRIPT_DIR="$(cd "$(dirname "${SCRIPT_PATH}")" ; pwd)"

EXECUTABLE=run.sh

RUNNING_FILE="$SCRIPT_DIR"/"$EXECUTABLE"
if [ ! -x $RUNNING_FILE ]; then
    echo "Cannot find $SCRIPT_DIR/$EXECUTABLE"
    echo "The file is absent or does not have execute permission"
    echo "This file is needed to run this program"
    exit 1
fi

exec "$SCRIPT_DIR"/"$EXECUTABLE" restart "$@"

```

run.sh

```
#!/bin/bash
#CREATE BY liujigang
#CREAte DATE 2018-04-23

#############################################
#  args: start  -- startup rpc              #
#        stop   -- stop rpc                 #
#        restart   -- restart rpc              #
#        status -- check the rpc status     #
#############################################

PRG="${BASH_SOURCE[0]}"

#PRGDIR工程当前目录丿
PRGDIR="$(cd "$(dirname "${PRG}")" ; pwd)"
echo "Now your current running environment is $PRGDIR"

jvmops="-Druntime.env=dev"
fname=`ls $PRGDIR/*jar`

status() {
 if [ -f $fname".pid" ]
 then
  pid=`awk '{print $1}' $fname".pid"`
  echo $fname" is running and pid is "$pid
 else
  echo $fname" is not running"
 fi
}

start() {
 if [ -f $fname".pid" ] 
 then
  status
  echo "It can't be runned twice!"
 else
  echo "starting "$fname
  nohup  /opt/jdk1.8.0_161/bin/java  $jvmops -jar $fname start >> /dev/null 2>&1 &
  echo $! > $fname".pid"
  status
 fi
}

stop() {
 echo "stoping "$fname
 ps -ef | grep $fname |grep -v grep | awk '{print $2}' | xargs kill
 rm -f $fname".pid"
 echo $fname" stoped"
}

case $1 in

 status)
  status
 ;;

 start)
  start
 ;;

 stop)
  stop 
 ;;

 restart)
  stop
  start
 ;;

 *)
  echo "#########################################"
  echo "  args: start  -- startup app            "
  echo "        stop   -- stop app               "
  echo "        restart-- restart app            "
  echo "        status -- check the app status   "
  echo "#########################################"
 ;;

esac
```