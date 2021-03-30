#!/bin/bash

do_install(){

echo "check file resources/repo/kernel/kernel-${version}.tar.gz..."

if [ ! -f resources/repo/kernel/kernel-"${version}".tar.gz ];then
    echo "check file 'resources/repo/kernel/kernel-${version}.tar.gz'failed..."
    exit 1
fi

# 生成repo
cat > /etc/yum.repos.d/kernel-${version}.repo <<EOF
[kernel-${version}]
name=[kernel-${version}]-repo
baseurl=file:///yum/data/"kernel-${version}"
gpgcheck=0
enabled=1
EOF

# 解压缩
mkdir -p /yum
tar zxvf ./resources/repo/kernel/kernel-"${version}".tar.gz -C /yum/

# 安装
yum install -y kernel-"${version}" --disablerepo=\* --enablerepo="kernel-${version}"

# shellcheck disable=SC2181
if [ $? -eq 0 ];then
    grub2-set-default 0 && grub2-mkconfig -o /etc/grub2.cfg
    grubby --args="user_namespace.enable=1" --update-kernel="$(grubby --default-kernel)"
fi

}

do_install

