#!/bin/bash
harborDir=/usr/local/harbor
version=v2.1.4
ssl=false
command_exists() {
    command -v "$@" > /dev/null 2>&1
}

echo "############ settings ############"
echo "version -> $version"
echo "ssl -> $ssl"
echo "domain -> $domain"
echo "data_dir -> $data_dir"
echo "password -> $password"
echo "############ settings ############"

if command_exists docker && [ -e /var/run/docker.sock ]; then
    echo "check docker pass..."
else
    echo "check docker failed..."
    exit 1
fi

if command_exists docker-compose; then
    echo "check docker-compose pass..."
else
    echo "check docker-compose failed..."
    exit 1
fi

function do_clean() {
  # shellcheck disable=SC2006
  docker ps|grep goharbor &> /dev/null
  # shellcheck disable=SC2181
  if [ $? -eq 0 ]; then
    docker ps|grep goharbor |awk '{print $1}'|xargs -n1 docker rm -f
  fi

  dataDirs=(ca_download database job_logs redis registry secret)
  for dir in "${dataDirs[@]}"
  do
    if [ "$data_dir" != "" ] && [ -d "$data_dir"/"$dir" ];then
      echo "clean $data_dir/$dir"
      rm -rf "$data_dir"/"$dir"
    fi
  done
}
do_clean

echo "check file /tmp/harbor-offline-installer-${version}.tgz..."

if [ ! -f /tmp/harbor-offline-installer-"${version}".tgz ];then
    echo "check file '/tmp/harbor-offline-installer-${version}.tgz'failed..."
    exit 1
fi

tar zxvf /tmp/harbor-offline-installer-"${version}".tgz -C /usr/local

do_config(){
  # shellcheck disable=SC2225
\cp "$harborDir"/harbor.yml.tmpl "$harborDir"/harbor.yml

# harbor_admin_password: Harbor12345
# shellcheck disable=SC2154
sudo sed -i -e "s|reg.mydomain.com|$domain|g" \
            -e "s|harbor_admin_password: Harbor12345|harbor_admin_password: $password|g" \
            -e "s|port: 80|port: $http_port|g" \
            -e "s|data_volume: /data|data_volume: $data_dir|g" \
            "$harborDir"/harbor.yml

mkdir -p "$data_dir"

# load image

if [ -f "$harborDir"/harbor."${version}".tar.gz ]; then
    echo "load image..."
    # shellcheck disable=SC2086
    docker load -i $harborDir/harbor.${version}.tar.gz
fi

# config ssl
if [ "$ssl" == "false" ]; then
    sed -i \
        -e "s;https:;#https:;g" \
        -e "s;port: 443;#port: 443;g" \
        -e "s;certificate:;#certificate:;g" \
        -e "s;private_key:;#private_key:;g" \
        "$harborDir"/harbor.yml
fi

echo "config service boot..."
sed -i '/docker-compose/d' /etc/rc.local
cat <<EOF >>/etc/rc.local
docker-compose -f $harborDir/docker-compose.yml down
docker-compose -f $harborDir/docker-compose.yml up -d
EOF

chmod +x /etc/rc.local

echo "config domain resolve..."

sed -i "/$domain/d" /etc/hosts
echo "$resolv_ip $domain" >> /etc/hosts

sh -c "$harborDir"/prepare


# install
sh -c "$harborDir"/install.sh

sed -i "/volumes:/a\      - /etc/hosts:/etc/hosts" $harborDir/docker-compose.yml
cd $harborDir
docker-compose down
docker-compose up -d

firewall-cmd --zone=public --add-port="$http_port"/tcp --permanent
firewall-cmd --reload

}

do_config