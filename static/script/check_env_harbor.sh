#!/bin/bash
command_exists() {
    command -v "$@" > /dev/null 2>&1
}

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