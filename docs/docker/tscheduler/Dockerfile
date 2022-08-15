FROM harbor-api.service.consul/k8stxbp/infra/centos:7.9.2009.2
MAINTAINER liujigang "liujigang@tisson.cn"
# script for creating tscheduler instance
RUN yum -y install java-1.8.0-openjdk java-1.8.0-openjdk-devel sudo openssh  \
    && groupadd tscheduler \
    && useradd tscheduler -g tscheduler \
    && echo "tscheduler" | passwd --stdin tscheduler \
    && sed -i '$atscheduler  ALL=(ALL)  NOPASSWD: NOPASSWD: ALL' /etc/sudoers \
    && sed -i 's/Defaults    requirett/#Defaults    requirett/g' /etc/sudoers \
    && mkdir -p /opt/tscheduler \
    && chown -R tscheduler:tscheduler /opt/tscheduler
RUN su tscheduler \
    && ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa \
    && cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys \
    && chmod 600 ~/.ssh/authorized_keys
ADD ./tscheduler /opt/tscheduler