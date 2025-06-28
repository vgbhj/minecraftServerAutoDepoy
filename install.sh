#!/bin/bash

error_exit() {
    echo "500"
    echo "Error: $1" >&2
    exit 1
}

# Detect package manager
if which apt-get > /dev/null 2>&1; then
    pm=$(which apt-get)
    silent_inst="-yq install"
    check_pkgs="-yq update"
    wget_pkg="wget"
    git_pkg="git"
    dist="debian"
elif which yum > /dev/null 2>&1; then
    pm=$(which yum)
    silent_inst="-y -q install"
    check_pkgs="-y -q check-update"
    wget_pkg="wget"
    git_pkg="git"
    dist="centos"
elif which dnf > /dev/null 2>&1; then
    pm=$(which dnf)
    silent_inst="-yq install"
    check_pkgs="-yq check-update"
    wget_pkg="wget"
    git_pkg="git"
    dist="fedora"
elif which zypper > /dev/null 2>&1; then
    pm=$(which zypper)
    silent_inst="-nq install"
    check_pkgs="-nq refresh"
    wget_pkg="wget"
    git_pkg="git"
    dist="opensuse"
elif which pacman > /dev/null 2>&1; then
    pm=$(which pacman)
    silent_inst="-Sy --noconfirm --noprogressbar --quiet"
    check_pkgs="-Sup"
    wget_pkg="wget"
    git_pkg="git"
    dist="archlinux"
else
    error_exit "Package manager not found"
    exit 1
fi

echo "Dist: $dist, Packet manager: $pm"

if [ "$dist" = "debian" ]; then export DEBIAN_FRONTEND=noninteractive; fi

# Install required packages
if ! command -v sudo > /dev/null 2>&1; then 
    $pm $check_pkgs || error_exit "Failed to update repositories"; 
    $pm $silent_inst sudo || error_exit "Failed to install sudo"; 
fi
if ! command -v wget > /dev/null 2>&1; then 
    sudo $pm $check_pkgs || error_exit "Failed to update repositories"; 
    sudo $pm $silent_inst $wget_pkg || error_exit "Failed to install wget"; 
fi
if ! command -v git > /dev/null 2>&1; then 
    sudo $pm $check_pkgs || error_exit "Failed to update repositories"; 
    sudo $pm $silent_inst $git_pkg || error_exit "Failed to install git"; 
fi

# CentOS: always install Docker
if [ "$dist" = "centos" ]; then
    # Install Docker if not present
    if ! command -v docker &> /dev/null; then
        sudo $pm $check_pkgs || error_exit "Failed to update repositories"
        sudo $pm $silent_inst yum-utils device-mapper-persistent-data lvm2 || error_exit "Failed to install Docker dependencies"
        sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo || error_exit "Failed to add Docker repo"
        sudo $pm $silent_inst docker-ce docker-ce-cli containerd.io || error_exit "Failed to install Docker"
        sudo systemctl start docker || error_exit "Failed to start Docker"
        sudo systemctl enable docker || error_exit "Failed to enable Docker autostart"
        sudo usermod -aG docker $USER || error_exit "Failed to add user to docker group"
    fi
    # Start Docker service (again, just in case)
    sudo systemctl start docker || error_exit "Failed to start Docker"
    sudo systemctl enable docker || error_exit "Failed to enable Docker autostart"
else
    # For all other distros, keep your previous logic (install Docker as needed)
    if ! command -v docker &> /dev/null; then
        wget -qO- https://raw.githubusercontent.com/amnezia-vpn/amnezia-client/refs/heads/dev/client/server_scripts/install_docker.sh | sh || error_exit "Failed to install Docker"
        sudo usermod -aG docker $USER || error_exit "Failed to add user to docker group"
        newgrp docker
    fi
    # Start Docker service
    if ! command -v systemctl &> /dev/null; then
        sudo dockerd &> /dev/null &
    else
        sudo systemctl start docker || error_exit "Failed to start Docker"
        sudo systemctl enable docker || error_exit "Failed to enable Docker autostart"
    fi
fi

# Install Docker Compose for all distros the same way
if ! command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE_VERSION="v2.27.0"
    sudo wget -O /usr/local/bin/docker-compose \
        "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-linux-$(uname -m)" || error_exit "Failed to download Docker Compose"
    sudo chmod +x /usr/local/bin/docker-compose || error_exit "Failed to make Docker Compose executable"
fi

# Clone repository
REPO="https://github.com/vgbhj/minecraftServerAutoDepoy.git"
TARGET_DIR="/opt/mcSAD"
sudo rm -rf "$TARGET_DIR" 2>/dev/null
sudo git clone "$REPO" "$TARGET_DIR" || error_exit "Failed to clone repository"

# Generate admin password
ADMIN_PASSWORD=$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 16)
echo "ADMIN_PASSWORD=$ADMIN_PASSWORD" | sudo tee -a "$TARGET_DIR/webApp/.env" >/dev/null

# Deploy application
cd "/opt/mcSAD/webApp" || error_exit "Failed to enter project directory"
/usr/local/bin/docker-compose down || error_exit "Failed to stop and remove existing containers"
/usr/local/bin/docker-compose up -d --build || error_exit "Failed to run docker-compose"

# Get server IP
SERVER_IP=$(ip route get 8.8.8.8 | grep -oP 'src \K[\d.]+')
if [ -z "$SERVER_IP" ]; then
    SERVER_IP=$(hostname -I | awk '{print $1}')  
fi
if [ -z "$SERVER_IP" ]; then
    SERVER_IP="localhost"  
fi

echo "200"
echo "Done! The application has been successfully deployed in Docker."
echo "Admin panel is available at: http://${SERVER_IP}:8080"
echo "Admin password: $ADMIN_PASSWORD"
exit 0