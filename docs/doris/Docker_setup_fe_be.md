
docker network create --driver bridge --subnet=172.20.80.0/24 doris-network


# FE

sudo docker run -itd \
--name=fe \
--env FE_SERVERS="fe1:172.20.80.2:9010" \
--env FE_ID=1 \
-p 8030:8030 \
-p 9030:9030 \
-v /data/doris/fe/doris-meta:/opt/apache-doris/fe/doris-meta \
-v /data/doris/fe/conf:/opt/apache-doris/fe/conf \
-v /data/doris/fe/log:/opt/apache-doris/fe/log \
--network=doris-network \
--ip=172.20.80.2 \
apache-doris:1.2.7.1-fe


docker run -itd \
--name=fe \
--env FE_SERVERS="fe1:172.20.80.2:9010" \
--env FE_ID=1 \
-p 8030:8030 \
-p 9030:9030 \
-v /data/fe/doris-meta:/opt/apache-doris/fe/doris-meta \
-v /data/fe/conf:/opt/apache-doris/fe/conf \
-v /data/fe/log:/opt/apache-doris/fe/log \
--network=doris-network \
--ip=172.20.80.2 \
apache/doris:1.2.1-fe-x86_64
# BE



``` 
#!/bin/bash
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

# choose a base image
FROM openjdk:8u342-jdk

# set environment variables
ENV JAVA_HOME="/usr/local/openjdk-8/" \
	PATH="/opt/apache-doris/be/bin:${PATH}"

# download the software to the mirror and replace it as needed
ADD ./resource/apache-doris-be-x.x.x-bin-x86_64.tar.gz /opt/

# deploy software
RUN apt-get update && \
	apt-get install -y default-mysql-client && \
	apt-get clean && \
	mkdir /opt/apache-doris && \
	cd /opt && \
	mv apache-doris-be-x.x.x-bin-x86_64 /opt/apache-doris/be

ADD resource/init_be.sh /opt/apache-doris/be/bin
RUN chmod 755 /opt/apache-doris/be/bin/init_be.sh

ENTRYPOINT ["/opt/apache-doris/be/bin/init_be.sh"]

```

sudo docker run -itd \
--name=be \
--env FE_SERVERS="fe1:172.16.27.34:9010" \
--env BE_ADDR="172.16.27.34:9050" \
-p 8040:8040 \
-v /data/doris/be/storage:/opt/apache-doris/be/storage \
-v /data/doris/be/log:/opt/apache-doris/be/log \
--net=host \
apache-doris:1.2.7.1-be

sudo docker run -itd \
--name=be \
--env FE_SERVERS="fe1:172.20.80.2:9010" \
--env BE_ADDR="172.20.80.3:9050" \
-p 8040:8040 \
-v /data/doris/be/storage:/opt/apache-doris/be/storage \
-v /data/doris/be/conf:/opt/apache-doris/be/conf \
-v /data/doris/be/log:/opt/apache-doris/be/log \
--network=doris-network \
--ip=172.20.80.3 \
apache-doris:1.2.7.1-be
