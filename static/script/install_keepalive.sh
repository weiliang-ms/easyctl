#!/bin/bash

interface_name=$1
master_ip=$2
slave_ip=$3
virtual_ip=$4
local_ip=$5

role=""
peer_ip=""


# shellcheck disable=SC2086
if [ $local_ip == $slave_ip ]; then
    role="BACKUP"
    peer_ip=$master_ip
else
  if [ $local_ip == $master_ip ]; then
    role="MASTER"
    peer_ip=$slave_ip
  fi
fi

echo "interface:$interface_name master:$master_ip slave:$slave_ip vip:$virtual_ip local:$local_ip"

SOFT_NAME=keepalived
filepath="/tmp/$SOFT_NAME.tar.gz"

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

do_install(){

if command_exists $SOFT_NAME;then
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

# 生成文件
echo "[$SOFT_NAME] config $SOFT_NAME..."
cat >/etc/keepalived/keepalived.conf <<EOF
global_defs {
  notification_email {
  }
  router_id LVS_DEVEL
  vrrp_skip_check_adv_addr
  vrrp_garp_interval 0
  vrrp_gna_interval 0
}

vrrp_instance haproxy-vip {
  state $role
  priority 100
  interface $interface_name
  virtual_router_id 60
  advert_int 1
  authentication {
    auth_type PASS
    auth_pass keepalive
  }
  unicast_src_ip $local_ip
  unicast_peer {
    $peer_ip
  }

  virtual_ipaddress {
    $virtual_ip/24                  #vip地址
  }
}
EOF

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
    echo "[$SOFT_NAME] install $SOFT_NAME failed..." >&2
fi

# 启动
echo "[$SOFT_NAME] boot $SOFT_NAME..."
systemctl enable $SOFT_NAME --now
# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
    echo "[$SOFT_NAME] boot $SOFT_NAME failed..." >&2
fi

}


do_install

