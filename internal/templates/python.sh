#!/bin/bash

# Auto-Deployer Cloud-Init Script for Python
# This script runs when the VM boots up

set -e  # Exit on any error

# Log everything
exec > >(tee -a /var/log/deploy.log) 2>&1
echo "Starting deployment at $(date)"

# Update system
apt-get update -y

# Install Python, pip, and other essentials
apt-get install -y python3 python3-pip python3-venv git

# Install PM2 (we'll use it to manage Python processes too)
apt-get install -y nodejs npm
npm install -g pm2

# Clone repository
echo "Cloning repository: {{REPO_URL}}"
cd /home
git clone {{REPO_URL}} app
cd app

# Create virtual environment
python3 -m venv venv
source venv/bin/activate

# Install dependencies
echo "Installing Python dependencies..."
{{INSTALL_COMMAND}}

# Create PM2 ecosystem file for Python
cat > ecosystem.config.js << EOF
module.exports = {
  apps: [{
    name: 'python-app',
    script: '{{START_COMMAND}}',
    interpreter: '/home/app/venv/bin/python3',
    cwd: '/home/app',
    env: {
      PATH: '/home/app/venv/bin:' + process.env.PATH
    }
  }]
};
EOF

# Start application with PM2
echo "Starting Python application..."
pm2 start ecosystem.config.js
pm2 startup
pm2 save

echo "Deployment completed successfully at $(date)"

# Output final status
pm2 list