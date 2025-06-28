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
    $pm $silent_install sudo || error_exit "Failed to install sudo"; 
fi
if ! command -v wget > /dev/null 2>&1; then 
    sudo $pm $check_pkgs || error_exit "Failed to update repositories"; 
    sudo $pm $silent_inst $wget_pkg || error_exit "Failed to install wget"; 
fi
if ! command -v git > /dev/null 2>&1; then 
    sudo $pm $check_pkgs || error_exit "Failed to update repositories"; 
    sudo $pm $silent_inst $git_pkg || error_exit "Failed to install git"; 
fi

# Install Podman if not present
if ! command -v podman &> /dev/null; then
    sudo $pm $check_pkgs || error_exit "Failed to update repositories"
    sudo $pm $silent_inst podman || error_exit "Failed to install podman"
fi

# Install podman-compose if not present
if ! command -v podman-compose &> /dev/null; then
    if [ "$dist" = "debian" ] || [ "$dist" = "ubuntu" ]; then
        sudo $pm $check_pkgs || error_exit "Failed to update repositories"
        sudo $pm $silent_inst podman-compose || error_exit "Failed to install podman-compose"
    else
        # Fallback to pip3 if no package available
        if command -v pip3 &> /dev/null; then
            sudo pip3 install podman-compose || error_exit "Failed to install podman-compose"
        else
            error_exit "pip3 not found and podman-compose package not available"
        fi
    fi
fi

# Start Podman service (for rootless, may not be needed)
if command -v systemctl &> /dev/null; then
    systemctl --user start podman.socket || true
fi

# Clone repository
REPO="https://github.com/vgbhj/minecraftServerAutoDepoy.git"
TARGET_DIR="/opt/mcSAD"
sudo rm -rf "$TARGET_DIR" 2>/dev/null
sudo git clone "$REPO" "$TARGET_DIR" || error_exit "Failed to clone repository"

# Deploy application
cd "/opt/mcSAD/webApp" || error_exit "Failed to enter project directory"
sudo podman-compose down || error_exit "Failed to stop and remove existing containers"
sudo podman-compose up -d --build || error_exit "Failed to run podman-compose"

# Get server IP
SERVER_IP=$(ip route get 8.8.8.8 | grep -oP 'src \K[\d.]+')
if [ -z "$SERVER_IP" ]; then
    SERVER_IP=$(hostname -I | awk '{print $1}')  
fi
if [ -z "$SERVER_IP" ]; then
    SERVER_IP="localhost"  
fi

echo "200"
echo "Done! The application has been successfully deployed in Podman."
echo "Admin panel is available at: http://${SERVER_IP}:8080"
exit 0