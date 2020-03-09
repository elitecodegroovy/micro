
Steps to take hot backup

Method 1:

```
1.>alter tablespace hrms_ts begin backup;

2.$cp hrms *.* /opt/hotbkp

3.>alter tablespace hrms_TS end backup;

4.Repeat step 1 to step 3 for all databases

5.>alter database backup control file to ‘/opt/hotbkp/control01.ctl’

6.>alter system switch logfile
```
 

Method 2:

```
1.>alter database begin backup

2.$cp *.dbf  /opt/hotbkp

3.>alter database end backup

4.>alter database backup control file to ‘/opt/hotbkp/control01.ctl’

5.>alter system switch logfile
```
