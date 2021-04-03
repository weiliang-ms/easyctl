#!/bin/bash

set -e
dns_address=${address:-$1}

sed -i "/$dns_address/d" /etc/resolv.conf
echo "nameserver $dns_address" >> /etc/resolv.conf

echo "[dns] 配置成功..."

echo "[dns] 当前dns列表: "
cat /etc/resolv.conf
