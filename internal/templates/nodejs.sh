!#bin/bash 


# cloiud init ke liye 

echo "Nodejs setup started"

set -e 


exec > >(tee -a /var/log/deploy.log) 2>&1
echo "Starting deployment at $(date)"

# Update system
apt-get update -y

# Install Node.js 18.x
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
apt-get install -y nodejs git

# Install PM2 globally
npm install -g pm2

# Clone repository
echo "Cloning repository: {{REPO_URL}}"
cd /home
git clone {{REPO_URL}} app
cd app

# Install dependencies
echo "Installing dependencies..."
{{INSTALL_COMMAND}}

# Start application with PM2
echo "Starting application..."
pm2 start {{START_COMMAND}} --name "deployed-app"
pm2 startup
pm2 save

# Create a simple health check endpoint (optional)
echo "Deployment completed successfully at $(date)"
echo "Application should be running on port 3000 (or as configured)"

# Output final status
pm2 list