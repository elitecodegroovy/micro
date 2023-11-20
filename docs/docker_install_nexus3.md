
##install nexus3

.拉取镜像
```
docker pull sonatype/nexus3 

``
2.启动容器
```
docker run -d --name nexus3  -p 8081:8081 -v /opt/nexus3/data:/var/nexus-data sonatype/nexus3
```

3.修改maven settings.xml 这里使用默认用户名 admin 密码 admin123

```
        <?xml version="1.0" encoding="UTF-8"?>
        <settings xmlns="http://maven.apache.org/SETTINGS/1.0.0"
                  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                  xsi:schemaLocation="http://maven.apache.org/SETTINGS/1.0.0 http://maven.apache.org/xsd/settings-1.0.0.xsd">
        <localRepository>H:\localRepository</localRepository>
        <pluginGroups></pluginGroups>    
          <proxies></proxies>    
          <servers>
             <server>
              <id>nexus</id>
              <username>admin</username>
              <password>admin123</password>
            </server>
        
            <server>
              <id>maven-public</id>
              <username>admin</username>
              <password>admin123</password>
            </server>
        
            <server>
              <id>maven-releases</id>
              <username>admin</username>
              <password>admin123</password>
            </server>
        
            <server>
              <id>maven-snapshots</id>
              <username>admin</username>
              <password>admin123</password>
            </server>
        
            <server>
              <id>3rd-part</id>
              <username>admin</username>
              <password>admin123</password>
            </server>
        
            <server>
              <id>maven-central</id>
              <username>admin</username>
              <password>admin123</password>
            </server>
          </servers>
         <mirrors>
           <mirror>
              <id>maven-public</id>
              <name>public</name>
              <mirrorOf>central</mirrorOf>
              <url>http://192.168.91.137:8081/repository/maven-public/</url>
            </mirror>
        
            <mirror>
              <id>maven-releases</id>
              <name>releases</name>
              <mirrorOf>central</mirrorOf>
              <url>http://192.168.91.137:8081/repository/maven-releases/</url>
            </mirror>
        
            <mirror>
              <id>maven-snapshots</id>
              <name>snapshots</name>
              <mirrorOf>central</mirrorOf>
              <url>http://192.168.91.137:8081/repository/maven-snapshots/</url>
            </mirror>
        
           <mirror>
              <id>3rd-part</id>
              <name>3rd-part</name>
              <mirrorOf>central</mirrorOf>
              <url>http://192.168.91.137:8081/repository/3rd-part/</url>
            </mirror>
        
               <mirror>
              <id>maven-central</id>
              <name>maven-central</name>
              <mirrorOf>central</mirrorOf>
              <url>http://192.168.91.137:8081/repository/maven-central/</url>
            </mirror>
         </mirrors>
        
        <profiles>
          <profile>
              <id>nexus</id>
              <repositories>
                <repository>
                  <id>nexus</id>
                  <url>http://192.168.91.137:8081/repository/maven-public/</url>
                  <releases>
                    <enabled>true</enabled>
                    <updatePolicy>always</updatePolicy>
                  </releases>
                  <snapshots>
                    <enabled>true</enabled>
                    <updatePolicy>always</updatePolicy>
                  </snapshots>
                </repository>
              </repositories>
              <pluginRepositories>
                <pluginRepository>
                  <id>nexus</id>
                  <url>http://192.168.91.137:8081/repository/maven-public/</url>
                  <releases>
                    <enabled>true</enabled>
                    <updatePolicy>always</updatePolicy>
                  </releases>
                  <snapshots>
                    <enabled>true</enabled>
                    <updatePolicy>always</updatePolicy>
                  </snapshots>
                </pluginRepository>
              </pluginRepositories>
            </profile>
          </profiles>
        
        <activeProfiles>
            <activeProfile>nexus</activeProfile>
         </activeProfiles>
        
        </settings>
```

4.如果项目需要发布到nexus,修改pom 添加以下 distributionManagement 内容

```
    <?xml version="1.0" encoding="UTF-8"?>
    <project xmlns="http://maven.apache.org/POM/4.0.0"
             xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
             xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
        <modelVersion>4.0.0</modelVersion>
    
        <groupId>org.lzw</groupId>
        <artifactId>idworker</artifactId>
        <version>1.0-SNAPSHOT</version>
    
        <distributionManagement>
    
            <repository>
                <id>maven-releases</id>
                <name>maven-releases</name>
                <url>http://192.168.91.137:8081/repository/maven-releases/</url>
            </repository>
    
            <snapshotRepository>
                <id>maven-snapshots</id>
                <name>maven-snapshots</name>
                <url>http://192.168.91.137:8081/repository/maven-snapshots/</url>
            </snapshotRepository>
    
        </distributionManagement>
    
    </project>
```
5.发布


docker run -d --name nexus3  -p 8081:8081 -v /opt/nexus3/data:/var/nexus-data sonatype/nexus3

docker run -d --name fe-doris  --net=host --restart=always apache-doris:1.2.7.1-fe