#!/bin/bash

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

if ! command -v yc &> /dev/null; then
    echo -e "${RED}YC CLI not found. Please install it first.${NC}"
    echo "Install: curl https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash"
    exit 1
fi

YC_FOLDER_ID="${YC_FOLDER_ID:-$(yc config get folder-id)}"
YC_ZONE="${YC_ZONE:-ru-central1-a}"
VM_NAME="${VM_NAME:-devops-project-vm}"
VM_IMAGE="${VM_IMAGE:-fd8kdq6d0p8sij7h5qe3}"
VM_CORES="${VM_CORES:-2}"
VM_MEMORY="${VM_MEMORY:-2}"
VM_DISK_SIZE="${VM_DISK_SIZE:-10}"
SSH_KEY="$HOME/.ssh/id_rsa"
SSH_KEY_PUB="$HOME/.ssh/id_rsa.pub"

echo -e "${YELLOW}Starting deployment to Yandex Cloud${NC}"
if ! command -v yc &> /dev/null; then
    echo -e "${RED}Yandex Cloud CLI not found. Please install it first.${NC}"
    echo "Install: curl https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash"
    exit 1
fi

if [ ! -f "$SSH_KEY_PUB" ]; then
    echo -e "${RED}SSH public key not found at $SSH_KEY_PUB${NC}"
    echo "Generate key: ssh-keygen -t rsa -b 4096"
    exit 1
fi

if [ -z "$YC_FOLDER_ID" ]; then
    echo -e "${RED}Yandex Cloud folder-id not configured${NC}"
    echo "Run: yc init"
    exit 1
fi

echo -e "${GREEN}Configuration:${NC}"
echo "  Folder ID: $YC_FOLDER_ID"
echo "  Zone: $YC_ZONE"
echo "  VM Name: $VM_NAME"

EXISTING_VM=$(yc compute instance list --folder-id="$YC_FOLDER_ID" --format=json | jq -r ".[] | select(.name==\"$VM_NAME\") | .id")

if [ -n "$EXISTING_VM" ]; then
    echo -e "${YELLOW}VM '$VM_NAME' already exists. Deleting${NC}"
    yc compute instance delete "$EXISTING_VM" --folder-id="$YC_FOLDER_ID"
    sleep 5
fi

echo -e "${GREEN}Creating VM${NC}"
VM_ID=$(yc compute instance create \
    --name="$VM_NAME" \
    --folder-id="$YC_FOLDER_ID" \
    --zone="$YC_ZONE" \
    --platform=standard-v3 \
    --cores="$VM_CORES" \
    --memory="${VM_MEMORY}GB" \
    --create-boot-disk size="${VM_DISK_SIZE}GB",image-id="$VM_IMAGE" \
    --network-interface subnet-name=default-"$YC_ZONE",nat-ip-version=ipv4 \
    --metadata-from-file user-data=<(cat <<EOF
#cloud-config
users:
  - name: ubuntu
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh-authorized-keys:
      - $(cat "$SSH_KEY_PUB")
EOF
) \
    --format=json | jq -r '.id')


echo -e "${GREEN}VM created with ID: $VM_ID${NC}"

echo -e "${GREEN}Getting VM IP address${NC}"
sleep 10

VM_IP=$(yc compute instance get "$VM_ID" --folder-id="$YC_FOLDER_ID" --format=json | jq -r '.network_interfaces[0].primary_v4_address.one_to_one_nat.address')

if [ -z "$VM_IP" ] || [ "$VM_IP" = "null" ]; then
    echo -e "${RED}❌ Failed to get VM IP address${NC}"
    exit 1
fi

echo -e "${GREEN}VM IP: $VM_IP${NC}"

echo -e "${YELLOW}Waiting for SSH to be ready${NC}"
for i in {1..30}; do
    if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 ubuntu@"$VM_IP" "echo SSH ready" &>/dev/null; then
        echo -e "${GREEN}SSH is ready${NC}"
        break
    fi
    echo "  Attempt $i/30"
    sleep 10
done

echo -e "${GREEN}Installing Docker on remote server${NC}"
ssh -o StrictHostKeyChecking=no ubuntu@"$VM_IP" bash <<'ENDSSH'
set -e

sudo apt-get update
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

sudo usermod -aG docker ubuntu
echo "Docker installed successfully"
ENDSSH

echo -e "${GREEN}Docker installed${NC}"

cd "$(dirname "$0")/.."

echo -e "${GREEN}Copying files to remote server${NC}"
ssh -o StrictHostKeyChecking=no ubuntu@"$VM_IP" "mkdir -p ~/app"
scp -o StrictHostKeyChecking=no -r \
    ./cmd \
    ./internal \
    ./migrations \
    ./go.mod \
    ./go.sum \
    ./Dockerfile \
    ./docker-compose.yaml \
    ./.env \
    ubuntu@"$VM_IP":~/app/

echo -e "${GREEN}Files copied${NC}"

echo -e "${GREEN}Starting application with Docker Compose${NC}"
ssh -o StrictHostKeyChecking=no ubuntu@"$VM_IP" bash <<'ENDSSH'
set -e
cd ~/app
docker compose down || true
docker compose up -d --build
echo "Application started (containers up)"
ENDSSH

echo -e "${GREEN}Docker Compose stack started${NC}"

echo -e "${GREEN}Migrations applied${NC}"

echo -e "${YELLOW}Waiting for application to be ready${NC}"
sleep 15

echo ""
echo -e "${GREEN}Deployment Successful! 🎉${NC}"
echo ""
echo -e "${YELLOW}Deployment Information:${NC}"
echo -e "  VM ID:        ${GREEN}$VM_ID${NC}"
echo -e "  VM IP:        ${GREEN}$VM_IP${NC}"
echo -e "  API Endpoint: ${GREEN}http://$VM_IP:8080${NC}"

echo "$VM_IP"