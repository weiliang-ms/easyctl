#!/bin/bash
package_name=docker-ce

cat > /etc/yum.repos.d/"$package_name".repo <<EOF
[$package_name]
name=[$package_name]-repo
baseurl=file:///yum/data/"$package_name"
gpgcheck=0
enabled=1
EOF

# 解压缩
mkdir -p /yum
tar zxvf /tmp/"$package_name"-*.tar.gz -C /yum/

# 安装
yum install -y "$package_name" --disablerepo=\* --enablerepo="$package_name"

sed -i "/net.ipv4.ip_forward/d" /etc/sysctl.conf
cat >> /etc/sysctl.conf<<EOF
net.ipv4.ip_forward=1
EOF
sysctl -p

mkdir -p /etc/docker
cat > /etc/docker/daemon.json <<EOF
{
  "log-opts": {
    "max-size": "5m",
    "max-file":"3"
  }
}
EOF

systemctl daemon-reload
systemctl enable docker --now
