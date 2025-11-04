#!/bin/bash
set -e

echo "🔐 Setting up Basic Auth for Sentinel Mesh"
echo "=========================================="
echo ""

# Prompt for username and password
read -p "Enter username for Sentinel Mesh access (default: admin): " USERNAME
USERNAME=${USERNAME:-admin}

read -sp "Enter password for $USERNAME: " PASSWORD
echo ""
read -sp "Confirm password: " PASSWORD2
echo ""

if [ "$PASSWORD" != "$PASSWORD2" ]; then
    echo "❌ Passwords don't match!"
    exit 1
fi

if [ -z "$PASSWORD" ]; then
    echo "❌ Password cannot be empty!"
    exit 1
fi

echo ""
echo "📝 Generating htpasswd credentials..."

# Check if htpasswd is available
if ! command -v htpasswd &> /dev/null; then
    echo "❌ htpasswd command not found. Please install apache2-utils (Debian/Ubuntu) or httpd-tools (RHEL/CentOS)"
    exit 1
fi

# Generate htpasswd credentials
HTPASSWD=$(htpasswd -nb "$USERNAME" "$PASSWORD")

# Create namespace if it doesn't exist
kubectl get namespace sentinel-mesh &>/dev/null || kubectl create namespace sentinel-mesh

# Create secret for Traefik basic auth
kubectl create secret generic sentinel-mesh-auth \
  --from-literal=users="$HTPASSWD" \
  --namespace=sentinel-mesh \
  --dry-run=client -o yaml | kubectl apply -f -

echo "✅ Basic auth secret created in sentinel-mesh namespace"
echo ""
echo "Username: $USERNAME"
echo "Password: [hidden]"
echo ""
echo "🔒 You can now access Sentinel Mesh UIs with these credentials"
echo ""
echo "Next steps:"
echo "1. Apply the middleware: kubectl apply -f middleware.yaml"
echo "2. Apply the IngressRoutes: kubectl apply -f ingressroutes.yaml"
echo "3. Update Cloudflare Tunnel config to include:"
echo "   - hostname: sentinel-mesh.georg-nikola.com"
echo "     service: http://traefik.traefik.svc.cluster.local:80"
echo "4. Add DNS records via Terraform"
