#!/bin/bash

# 1. Установка Docker (через официальный скрипт)
wget -qO- https://get.docker.com/ | sh

# 2. Добавляем пользователя в группу docker (чтобы не нужен был sudo)
sudo usermod -aG docker $USER
newgrp docker  # Обновляем группу без перезагрузки

# 3. Устанавливаем Docker Compose (бинарник)
DOCKER_COMPOSE_VERSION="v2.27.0"
sudo wget -O /usr/local/bin/docker-compose \
    "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-linux-$(uname -m)"
sudo chmod +x /usr/local/bin/docker-compose

# 4. Запускаем Docker (если нет systemd)
if ! command -v systemctl &> /dev/null; then
    sudo dockerd &  # Запуск в фоне (для систем без systemd)
else
    sudo systemctl start docker
    sudo systemctl enable docker
fi

# 5. Скачиваем репозиторий (без git, через архив GitHub)
wget https://github.com/vgbhj/minecraftServerAutoDepoy/archive/refs/heads/main.zip -O /tmp/webApp.zip
sudo unzip /tmp/webApp.zip -d /opt/
sudo mv /opt/minecraftServerAutoDepoy-main /opt/webApp

# 6. Запускаем приложение (если main.go — это веб-сервер)
cd /opt/webApp/webApp
if command -v go &> /dev/null; then
    sudo go run main.go  # Если Go установлен
else
    echo "Go не установлен. Установите его вручную или используйте Docker."
fi