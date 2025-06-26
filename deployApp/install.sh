#!/bin/bash

# Определение пакетного менеджера и команд
if which apt-get > /dev/null 2>&1; then
    pm=$(which apt-get)
    silent_inst="-yq install"
    check_pkgs="-yq update"
    wget_pkg="wget"
    git_pkg="git"
    dist="debian"
elif which dnf > /dev/null 2>&1; then
    pm=$(which dnf)
    silent_inst="-yq install"
    check_pkgs="-yq check-update"
    wget_pkg="wget"
    git_pkg="git"
    dist="fedora"
elif which yum > /dev/null 2>&1; then
    pm=$(which yum)
    silent_inst="-y -q install"
    check_pkgs="-y -q check-update"
    wget_pkg="wget"
    git_pkg="git"
    dist="centos"
elif which zypper > /dev/null 2>&1; then
    pm=$(which zypper)
    silent_inst="-nq install"
    check_pkgs="-nq refresh"
    wget_pkg="wget"
    git_pkg="git"
    dist="opensuse"
elif which pacman > /dev/null 2>&1; then
    pm=$(which pacman)
    silent_inst="-Syu --noconfirm --noprogressbar --quiet"
    check_pkgs="-Sup"
    wget_pkg="wget"
    git_pkg="git"
    dist="archlinux"
else
    echo "Packet manager not found"
    exit 1
fi

echo "Dist: $dist, Packet manager: $pm"

if [ "$dist" = "debian" ]; then export DEBIAN_FRONTEND=noninteractive; fi

if ! command -v sudo > /dev/null 2>&1; then $pm $check_pkgs; $pm $silent_inst sudo; fi
if ! command -v wget > /dev/null 2>&1; then sudo $pm $check_pkgs; sudo $pm $silent_inst $wget_pkg; fi
if ! command -v git > /dev/null 2>&1; then sudo $pm $check_pkgs; sudo $pm $silent_inst $git_pkg; fi

# 1. Установка Docker (через официальный скрипт)
if ! command -v docker &> /dev/null; then
    wget -qO- https://raw.githubusercontent.com/vgbhj/minecraftServerAutoDepoy/refs/heads/main/install_docker.sh | sh
    sudo usermod -aG docker $USER
    newgrp docker
fi

# 2. Установка Docker Compose (бинарник)
if ! command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE_VERSION="v2.27.0"
    sudo wget -O /usr/local/bin/docker-compose \
        "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-linux-$(uname -m)"
    sudo chmod +x /usr/local/bin/docker-compose
fi

# 3. Запуск Docker (с проверкой systemd)
if ! command -v systemctl &> /dev/null; then
    sudo dockerd &> /dev/null &
else
    sudo systemctl start docker
    sudo systemctl enable docker
fi

# 4. Клонирование репозитория через git
REPO="https://github.com/vgbhj/minecraftServerAutoDepoy.git"
TARGET_DIR="/opt/mcSAD"
sudo rm -rf "$TARGET_DIR" 2>/dev/null
sudo git clone "$REPO" "$TARGET_DIR"

# 5. Запуск через docker-compose
cd "/opt/mcSAD/webApp"
sudo docker-compose up -d --build

echo "Готово! Приложение запущено в Docker."