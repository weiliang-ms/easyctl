#!/bin/bash

do_install(){

echo "[kernel] check kernel-${version} exist ..."
rpm -qa|grep kernel-lt &> /dev/null
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
    echo "[kernel] kernel-${version} had been installed..."
    awk -F\' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg
    exit 0
fi

echo "[kernel] check file ${filepath}..."

if [ ! -f "${filepath}" ];then
    echo "[kernel] check offline file failed, exit process..." >&2
    exit 2
fi

echo "[kernel] set up kernel-${version}.repo repo..."
# 生成repo
# shellcheck disable=SC2086
cat > /etc/yum.repos.d/kernel-${version}.repo <<EOF
[kernel-${version}]
name=[kernel-${version}]-repo
baseurl=file:///yum/data/kernel-${version}
gpgcheck=0
enabled=1
EOF

echo "[kernel] decompress file..."
# 解压缩
mkdir -p /yum
tar zxvf "${filepath}" -C /yum/
echo "[kernel] install kernel-${version}..."

# 安装
# shellcheck disable=SC2086
yum install -y kernel-${version} \
    --disablerepo=\* --enablerepo=kernel-${version}

# shellcheck disable=SC2181
if [ $? -eq 0 ];then
    grub2-set-default 0 && grub2-mkconfig -o /etc/grub2.cfg
    grubby --args="user_namespace.enable=1" --update-kernel="$(grubby --default-kernel)"
fi

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
    echo "[kernel] upgrade kernel failed..." >&2
fi
echo "[kernel] upgrade kernel successful..."
awk -F\' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg
}

do_install

