#!/bin/bash

# 0. Установка git через wget (если его нет)
if ! command -v git &> /dev/null; then
    echo "Установка git..."
    GIT_VERSION="2.45.1"
    wget https://mirrors.edge.kernel.org/pub/software/scm/git/git-${GIT_VERSION}.tar.gz
    tar -xf git-${GIT_VERSION}.tar.gz
    cd git-${GIT_VERSION}
    ./configure --prefix=/usr/local
    make -j$(nproc)
    sudo make install
    cd ..
fi

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