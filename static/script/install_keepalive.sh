#!/bin/bash

SOFT_NAME=keepalive

do_install(){

echo "[$SOFT_NAME] check $SOFT_NAME.tar.gz exist ..."
rpm -qa|grep $SOFT_NAME &> /dev/null
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
    echo "[$SOFT_NAME] $SOFT_NAME had been installed..."
    exit 0
fi

echo "[$SOFT_NAME] check file ${filepath}..."

if [ ! -f "${filepath}" ];then
    echo "[kernel] check offline file failed, exit process..." >&2
    exit 2
fi

echo "[$SOFT_NAME] set up $SOFT_NAME.repo repo..."
# 生成repo
# shellcheck disable=SC2086
cat > /etc/yum.repos.d/$SOFT_NAME.repo <<EOF
[$SOFT_NAME]
name=$SOFT_NAME-repo
baseurl=file:///yum/data/$SOFT_NAME
gpgcheck=0
enabled=1
EOF

echo "[$SOFT_NAME] decompress file..."
# 解压缩
mkdir -p /yum
tar zxvf "${filepath}" -C /yum/
echo "[$SOFT_NAME] install $SOFT_NAME..."

# 安装
# shellcheck disable=SC2086
yum install -y $SOFT_NAME \
    --disablerepo=\* --enablerepo=$SOFT_NAME


# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
    echo "[$SOFT_NAME] install $SOFT_NAME failed..." >&2
fi


}

do_install

