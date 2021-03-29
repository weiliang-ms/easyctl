#!/bin/bash
# https://github.com/kubesphere/kubekey/blob/master/scripts/docker-install.sh

DOCKER_REPO=${$1:http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo}

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

detect_kernel_version(){
    echo "get current kernel version..."
    # shellcheck disable=SC2006
    KERNEL_MAIN_VERSION=`uname -r |awk -F'.' '{print $1}'`
    # shellcheck disable=SC2006
    if [ "`uname -r`" -le 3 ]; then
        echo "detect kernel failed..."
        exit 1
    fi
}

config_daemon(){
    sudo mkdir -p /etc/docker
    sudo cat > /etc/docker/daemon.json <<EOF
{
  "log-opts": {
    "max-size": "5m",
    "max-file":"3"
  },
  "exec-opts": ["native.cgroupdriver=systemd"],
}
EOF
    sudo systemctl daemon-reload
    sudo systemctl restart docker
}

do_install(){

    detect_kernel_version

    if command_exists docker && [ -e /var/run/docker.sock ]; then
        version="$(docker -v | cut -d ' ' -f3 | cut -d ',' -f1)"
        echo "$version"

        cat >&2 <<EOF
        Warning: the "docker" command appears to already exist on this system.
        If you already have Docker installed, this script can cause trouble, which is
        why we're displaying this warning and exit the
        installation.
EOF
        ( set -x; sleep 5 )

    fi

    user="$(id -un 2>/dev/null || true)"

    sh_c='sh -c'
    if [ "$user" != 'root' ]; then
        if command_exists sudo; then
            sh_c='sudo -E sh -c'
        elif command_exists su; then
            sh_c='su -c'
        else
            cat >&2 <<EOF
            Error: this installer needs the ability to run commands as root.
            We are unable to find either "sudo" or "su" available to make this happen.
EOF
            exit 1
        fi
    fi

    pkg_manager="yum"
    config_manager="yum-config-manager"
    enable_channel_flag="--enable"
    pre_reqs="yum-utils"

    (
        set -x
        if [ "$lsb_dist" = "redhat" ]; then
            for rhel_repo in $rhel_repos ; do
                $sh_c "$config_manager $enable_channel_flag $rhel_repo"
            done
            fi
            $sh_c "yum install -y -q $pre_reqs"
            $sh_c "$config_manager --add-repo $DOCKER_REPO"
            $sh_c "$pkg_manager makecache fast"
            $sh_c "$pkg_manager install -y -q docker-ce"
            $sh_c 'systemctl enable docker --now'
    )

    config_daemon

}

do_install

