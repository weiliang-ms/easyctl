#!/usr/bin/env bash

if [ ! -f ~/.ssh/id_rsa ]; then
    ssh-keygen -t rsa -N '' -f ~/.ssh/id_rsa -q
    cat ~/.ssh/id_rsa.pub > ~/.ssh/authorized_keys
fi

chmod 600 ~/.ssh -R