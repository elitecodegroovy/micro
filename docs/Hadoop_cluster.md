

## Hadoop Cluster

```
192.168.1.229   master
192.168.1.231   node1
192.168.1.232   node2
```


## Hadoop Env

```
export HADOOP_INSTALL=/opt/hadoop
export PATH=$PATH:$HADOOP_INSTALL/bin  
export PATH=$PATH:$HADOOP_INSTALL/sbin  
export HADOOP_MAPRED_HOME=$HADOOP_INSTALL  
export HADOOP_COMMON_HOME=$HADOOP_INSTALL  
export HADOOP_HDFS_HOME=$HADOOP_INSTALL  
export HADOOP_CONF_DIR=$HADOOP_INSTALL/etc/hadoop
export YARN_HOME=$HADOOP_INSTALL  

```

### core-site.xml

```
core-site.xml
<configuration>
<property>
      <name>fs.defaultFS</name>
      <value>hdfs://hadoopmaster:9000</value>
       <description>NameNode URI</description>
 </property>
</configuration>


```

### hdfs-site.xml

```
<configuration>
	<property>
		<name>dfs.permissions.enabled</name>
		<value>false</value>
	</property>
	<property>
 		<name>dfs.replication</name>
  		<value>2</value>
  	</property>
  	<property>
   		<name>dfs.namenode.name.dir</name>
    	<value>/opt/hadoop/hdfs/namenode</value>
    </property>
    <property>
     	<name>dfs.datanode.data.dir</name>
      	<value>/opt/hadoop/hdfs/datanode</value>
    </property>
</configuration>


```

###mapred.xml

```
<configuration>
	<property>
		<name>mapreduce.framework.name</name>
		<value>yarn</value>
		<description>Execution framework.</description>
	</property>
	<property>
		<name>mapreduce.jobtracker.address</name>
		<value>master:9001</value>
	</property>
</configuration>

```

### yarn-site.xml


```
<configuration>

<!-- Site specific YARN configuration properties -->
<property>
		  <name>yarn.resourcemanager.hostname</name>
		  <value>hadoop_master</value>
</property>
<property>
		<name>yarn.nodemanager.aux-services.mapreduce_shuffle.class</name>
		<value>org.apache.hadoop.mapred.ShuffleHandler</value>
</property>
<property>
	    <name>yarn.nodemanager.aux-services</name>
	    <value>mapreduce_shuffle</value>
</property>
<property>
	    <name>yarn.resourcemanager.address</name>
	    <value>master:8032</value>
</property>
<property>
	    <name>yarn.resourcemanager.scheduler.address</name>
	    <value>master:8030</value>
</property>
<property>
	    <name>yarn.resourcemanager.resource-tracker.address</name>
	    <value>master:8031</value>
</property>
<property>
	<name>yarn.nodemanager.vmem-pmem-ratio</name>
	<value>3</value>
</property>
</configuration>

```


##HBase Cluster


```
<property>
		<name>hbase.cluster.distributed</name>
		<value>true</value>
</property>
<property>
		<name>hbase.rootdir</name>
		<value>hdfs://master:9000/hbase</value>
</property>
<property>
          <name>hbase.zookeeper.quorum</name>
            <value>192.168.1.224:2222,192.168.1.225:2222,192.168.1.226:2222</value>
</property>
<property>
	     <name>hbase.zookeeper.property.clientPort</name>
		 <value>2222</value>
</property>
```