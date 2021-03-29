#!/bin/bash

#package_name=$1
#package_repo=$2

echo "导出 $package_name 及其依赖..."

if [ "$package_repo" != "" ]; then
  yum -y install yum-utils device-mapper-persistent-data lvm2
  yum-config-manager --add-repo "$package_repo"
fi

# 更新

yum update -y
yum install yum-plugin-downloadonly -y

mkdir -p ./resources/repo/"$package_name"
yum install --downloadonly --downloaddir=./resources/repo/"$package_name" "$package_name"

# 生成repo依赖关系

yum install -y createrepo
createrepo ./resources/repo/"$package_name"

# 压缩

tar zcvf ./resources/repo/"$package_name".tar.gz ./resources/repo/"$package_name"
rm -rf "$package_name"

# 生成repo

cat > ./resources/repo/"$1".repo <<EOF
[$1]
name=[$1]-repo
baseurl=file:///yum/$1
gpgcheck=0
enabled=1
EOF
