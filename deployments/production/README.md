# Sentinel Mesh - Production Deployment

This directory contains configuration files for deploying Sentinel Mesh to the production Kubernetes cluster (Talos).

## Architecture

The production cluster uses:
- **Traefik** as ingress controller
- **Cloudflare Tunnel** for secure external access (no open ports)
- **Basic Auth** for dashboard protection
- **Cloudflare DNS** managed via Terraform

## Prerequisites

1. Access to the production Talos Kubernetes cluster
2. `kubectl` configured with production cluster context
3. `htpasswd` utility installed (apache2-utils on Debian/Ubuntu)
4. Cloudflare account with tunnel configured
5. Docker images pushed to registry

## Deployment Steps

### 1. Set up Basic Authentication

Run the auth setup script to create credentials:

```bash
cd deployments/production
./auth-setup.sh
```

This will:
- Prompt for username/password
- Generate htpasswd credentials
- Create Kubernetes secret `sentinel-mesh-auth` in the `sentinel-mesh` namespace

### 2. Deploy Middleware and IngressRoutes

Apply the Traefik configurations:

```bash
kubectl apply -f namespace.yaml
kubectl apply -f middleware.yaml
kubectl apply -f ingressroutes.yaml
```

### 3. Deploy Sentinel Mesh Services

Deploy using Helm:

```bash
# From repository root
helm upgrade --install sentinel-mesh deployments/helm/sentinel-mesh \
  --namespace sentinel-mesh \
  --create-namespace \
  --set image.tag=<VERSION> \
  --set ingress.enabled=false  # We use IngressRoutes instead
```

### 4. Configure Cloudflare Tunnel

Update the Cloudflare Tunnel configuration to include Sentinel Mesh:

```yaml
# In your Cloudflare Tunnel configuration
ingress:
  - hostname: sentinel-mesh.example.com
    service: http://traefik.traefik.svc.cluster.local:80
  - hostname: sentinel-mesh-api.example.com
    service: http://traefik.traefik.svc.cluster.local:80
  # ... other services ...
  - service: http_status:404
```

Apply and restart:
```bash
# Apply updated tunnel configuration
kubectl apply -f /path/to/cloudflare-tunnel/config.yaml
kubectl rollout restart deployment/cloudflared -n cloudflare-tunnel
```

### 5. Add DNS Records

Add DNS records via Terraform:

```bash
cd ~/path/to/terraform

# Add to main.tf:
# resource "cloudflare_record" "sentinel_mesh" {
#   zone_id = data.cloudflare_zone.main.id
#   name    = "sentinel-mesh"
#   content = "${var.cloudflare_tunnel_id}.cfargotunnel.com"
#   type    = "CNAME"
#   proxied = true
#   ttl     = 1
# }
#
# resource "cloudflare_record" "sentinel_mesh_api" {
#   zone_id = data.cloudflare_zone.main.id
#   name    = "sentinel-mesh-api"
#   content = "${var.cloudflare_tunnel_id}.cfargotunnel.com"
#   type    = "CNAME"
#   proxied = true
#   ttl     = 1
# }

terraform apply
```

### 6. Verify Deployment

Wait ~2 minutes for DNS propagation, then access:

- **Frontend**: https://sentinel-mesh.example.com
- **API**: https://sentinel-mesh-api.example.com

You'll be prompted for the username/password you configured in step 1.

## Verification Commands

```bash
# Check deployment status
kubectl get pods -n sentinel-mesh

# Check services
kubectl get svc -n sentinel-mesh

# Check IngressRoutes
kubectl get ingressroute -n sentinel-mesh

# View logs
kubectl logs -n sentinel-mesh -l app=frontend
kubectl logs -n sentinel-mesh -l app=api

# Check auth secret
kubectl get secret sentinel-mesh-auth -n sentinel-mesh
```

## Updating Basic Auth Credentials

To update credentials, simply re-run the auth setup script:

```bash
./auth-setup.sh
```

The script will update the existing secret. No need to restart pods - Traefik will pick up the change automatically.

## Removing Basic Auth (Not Recommended)

If you need to remove basic auth protection:

```bash
# Remove middleware reference from IngressRoutes
kubectl edit ingressroute sentinel-mesh-frontend -n sentinel-mesh
# Delete the middlewares section

# Or apply IngressRoutes without middleware
kubectl apply -f ingressroutes-no-auth.yaml
```

## Troubleshooting

### 401 Unauthorized
- Verify secret exists: `kubectl get secret sentinel-mesh-auth -n sentinel-mesh`
- Check middleware is applied: `kubectl describe ingressroute sentinel-mesh-frontend -n sentinel-mesh`
- Try re-creating the secret with `./auth-setup.sh`

### 404 Not Found
- Check IngressRoute exists: `kubectl get ingressroute -n sentinel-mesh`
- Verify service is running: `kubectl get pods,svc -n sentinel-mesh`
- Check Traefik can reach the service:
  ```bash
  kubectl run test --rm -it --image=curlimages/curl -- \
    curl http://frontend.sentinel-mesh.svc.cluster.local
  ```

### DNS not resolving
- Check DNS record: `dig sentinel-mesh.example.com`
- Verify Terraform applied in your Terraform directory
- Check Cloudflare Dashboard for CNAME records

### Tunnel not routing traffic
- Check tunnel config: `kubectl get configmap cloudflared-config -n cloudflare-tunnel -o yaml`
- View tunnel logs: `kubectl logs -n cloudflare-tunnel -l app=cloudflared`
- Verify tunnel is registered: Look for "Registered tunnel connection" in logs

## Security Considerations

1. **HTTPS Only**: All traffic is encrypted via Cloudflare
2. **Basic Auth**: Dashboard protected with username/password
3. **No Open Ports**: Tunnel uses outbound-only connections
4. **DDoS Protection**: Built into Cloudflare
5. **Secret Management**: Store auth credentials securely

## Rollback

To rollback a deployment:

```bash
# Rollback to previous Helm release
helm rollback sentinel-mesh -n sentinel-mesh

# Or rollback to specific revision
helm rollback sentinel-mesh <REVISION> -n sentinel-mesh

# List revisions
helm history sentinel-mesh -n sentinel-mesh
```

## Monitoring

Check the monitoring stack for Sentinel Mesh metrics:

- **Grafana**: https://grafana.example.com
- **Prometheus**: https://prometheus.example.com

Sentinel Mesh services expose metrics on port 9090 at `/metrics` endpoint.
