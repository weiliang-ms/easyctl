#!/bin/bash
harborDir=/usr/local/harbor

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

echo "############ settings ############"
echo "version -> $version"
echo "ssl -> $ssl"
echo "domain -> $domain"
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

echo "check file resources/repo/harbor/harbor-offline-installer-${version}.tgz..."

if [ ! -f resources/repo/harbor/harbor-offline-installer-"${version}".tgz ];then
    echo "check file 'resources/repo/harbor/harbor-offline-installer-${version}.tgz'failed..."
    exit 1
fi

tar zxvf resources/repo/harbor/harbor-offline-installer-"${version}".tgz -C /usr/local

don_config(){
  # shellcheck disable=SC2225
\cp "$harborDir"/harbor.yml.tmpl "$harborDir"/harbor.yml

# shellcheck disable=SC2154
sed -i "s#reg.mydomain.com#$domain#g" "$harborDir"/harbor.yml

if [ -d "$dataDir" ]; then
    sed -i "s#data_volume: /data#$dataDir#g" "$harborDir"/harbor.yml
fi

# load image

if [ -f "$harborDir"/harbor."${version}".tar.gz ]; then
    echo "load image..."
    # shellcheck disable=SC2086
    docker load -i $harborDir/harbor.${version}.tar.gz
fi

# config ssl
if [ "$ssl" == "false" ]; then
    sed -i "s;https:;#https:;g" "$harborDir"/harbor.yml
    sed -i "s;port: 443;#port: 443;g" "$harborDir"/harbor.yml
    sed -i "s;certificate:;#certificate:;g" "$harborDir"/harbor.yml
    sed -i "s;private_key:;#private_key:;g" "$harborDir"/harbor.yml
fi

sed -i '/docker-compose/d' /etc/rc.local
cat <<EOF >>/etc/rc.local
docker-compose -f $harborDir/docker-compose.yml down
docker-compose -f $harborDir/docker-compose.yml up -d
EOF

sh -c "$harborDir"/prepare

# install
sh -c "$harborDir"/install.sh

}

don_config