#!/usr/bin/env bash

#!/bin/sh

orphanedPods=`cat /var/log/messages|grep 'orphaned pod'|awk -F '"' '{print $2}'|uniq`;
orphanedPodsNum=`echo $orphanedPods|awk -F ' ' '{print NF}'`;
echo -e "orphanedPods: $orphanedPodsNum \n$orphanedPods";

for i in $orphanedPods
do
echo "Deleting Orphaned pod id: $i";
rm -rf /var/lib/kubelet/pods/$i;
done