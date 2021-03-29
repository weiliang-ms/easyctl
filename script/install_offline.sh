#!/bin/bash

# 生成repo
cat > /etc/yum.repos.d/"$package_name".repo <<EOF
[$package_name]
name=[$package_name]-repo
baseurl=file:///yum/data/"$package_name"
gpgcheck=0
enabled=1
EOF

# 解压缩
mkdir -p /yum
tar zxvf ./resources/repo/"$package_name"/repo.tar.gz -C /yum/

# 安装
yum install -y "$package_name" --disablerepo=\* --enablerepo="$package_name"