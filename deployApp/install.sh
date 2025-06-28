#!/bin/bash

error_exit() {
    echo "500"
    echo "Ошибка: $1" >&2
    exit 1
}

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
    silent_inst="-S --noconfirm --noprogressbar --quiet"
    check_pkgs="-Sup"
    wget_pkg="wget"
    git_pkg="git"
    dist="archlinux"
else
    error_exit "Пакетный менеджер не найден"
fi

echo "Dist: $dist, Packet manager: $pm"

if [ "$dist" = "debian" ]; then export DEBIAN_FRONTEND=noninteractive; fi

# Установка необходимых пакетов
if ! command -v sudo > /dev/null 2>&1; then $pm $check_pkgs || error_exit "Не удалось обновить репозитории"; $pm $silent_install sudo || error_exit "Не удалось установить sudo"; fi
if ! command -v wget > /dev/null 2>&1; then sudo $pm $check_pkgs || error_exit "Не удалось обновить репозитории"; sudo $pm $silent_inst $wget_pkg || error_exit "Не удалось установить wget"; fi
if ! command -v git > /dev/null 2>&1; then sudo $pm $check_pkgs || error_exit "Не удалось обновить репозитории"; sudo $pm $silent_inst $git_pkg || error_exit "Не удалось установить git"; fi

# 1. Установка Docker
if ! command -v docker &> /dev/null; then
    wget -qO- https://raw.githubusercontent.com/amnezia-vpn/amnezia-client/refs/heads/dev/client/server_scripts/install_docker.sh | sh || error_exit "Не удалось установить Docker"
    sudo usermod -aG docker $USER || error_exit "Не удалось добавить пользователя в группу docker"
    newgrp docker
fi

# 2. Установка Docker Compose
if ! command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE_VERSION="v2.27.0"
    sudo wget -O /usr/local/bin/docker-compose \
        "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-linux-$(uname -m)" || error_exit "Не удалось загрузить Docker Compose"
    sudo chmod +x /usr/local/bin/docker-compose || error_exit "Не удалось сделать Docker Compose исполняемым"
fi

# 3. Запуск Docker
if ! command -v systemctl &> /dev/null; then
    sudo dockerd &> /dev/null &
else
    sudo systemctl start docker || error_exit "Не удалось запустить Docker"
    sudo systemctl enable docker || error_exit "Не удалось добавить Docker в автозагрузку"
fi

REPO="https://github.com/vgbhj/minecraftServerAutoDepoy.git"
TARGET_DIR="/opt/mcSAD"
sudo rm -rf "$TARGET_DIR" 2>/dev/null
sudo git clone "$REPO" "$TARGET_DIR" || error_exit "Не удалось клонировать репозиторий"


cd "/opt/mcSAD/webApp" || error_exit "Не удалось перейти в директорию проекта"
sudo docker-compose up -d --build || error_exit "Не удалось запустить docker-compose"

# Получаем IP-адрес сервера
SERVER_IP=$(hostname -I | awk '{print $1}')
if [ -z "$SERVER_IP" ]; then
    SERVER_IP="localhost"
fi

echo "200"
echo "Готово! Приложение успешно запущено в Docker."
echo "Админ-панель доступна по адресу: http://${SERVER_IP}:8080"
exit 0