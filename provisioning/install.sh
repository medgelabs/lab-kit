#!/bin/sh
# Initial install script for new nodes

MYIP=${1:-"192.168.1.41"}
CIDR=${2:-"192.168.1.1/24"}

# Disable cloud-init because sanity
touch /etc/cloud/cloud-init.disabled

# Install required tools
apt update && apt install -y \
  curl \
  htop \
  vim \
  iftop \
  ufw

# Automatic Security Updates
echo "unattended-upgrades unattended-upgrades/enable_auto_updates boolean true" | debconf-set-selections
apt install unattended-upgrades -y

# Firewall setup
ufw default deny incoming
ufw default allow outgoing

# Allow localhost and local network
ufw allow from 127.0.0.1
ufw allow from ${CIDR}

# Allow SSH to me, only
ufw allow from ${MYIP} ssh

# Allow K3S API to me, only
ufw allow from ${MYIP} to any port 6443

# Allow Kubernetes interfaces
ufw allow in on cni0
ufw allow out on cni0
ufw allow in on flannel.1
ufw allow out on flannel.1

# HTTP/HTTPS traffic to services
ufw allow http
ufw allow https

ufw enable
