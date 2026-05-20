# RecallFlow - Deployment Guide

## Production Deployment Architecture

```
┌─────────────────┐
│   Cloudflare    │ (DNS + CDN)
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼───────┐
│Vercel │ │AWS EC2   │
│(Next) │ │(Go API)  │
└───────┘ └────┬─────┘
                │
         ┌──────┴──────┐
         │             │
    ┌────▼────┐  ┌────▼────┐
    │AWS RDS  │  │ElastiCache│
    │(Postgres)│  │(Redis)  │
    └─────────┘  └─────────┘
```

## Backend Deployment (AWS EC2)

### 1. Provision Infrastructure

#### EC2 Instance
- Instance Type: `t3.small` (2 vCPU, 2GB RAM) for MVP
- OS: Ubuntu 22.04 LTS
- Storage: 20GB gp3
- Security Group:
  - Port 22 (SSH) - Your IP only
  - Port 8080 (API) - 0.0.0.0/0 (behind NGINX)
  - Port 443 (HTTPS) - 0.0.0.0/0

#### RDS PostgreSQL
- Instance Class: `db.t3.micro` for MVP
- Storage: 20GB gp3
- Multi-AZ: No (for MVP)
- Backup: 7 days retention

#### ElastiCache Redis
- Node Type: `cache.t3.micro`
- Number of nodes: 1

### 2. Server Setup

SSH into EC2:
```bash
ssh -i your-key.pem ubuntu@your-ec2-ip
```

Install dependencies:
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install NGINX
sudo apt install nginx -y

# Install certbot for SSL
sudo apt install certbot python3-certbot-nginx -y
```

### 3. Deploy Application

```bash
# Clone repository
git clone https://github.com/yourusername/recallflow.git
cd recallflow/backend

# Set up environment
cp .env.example .env
nano .env  # Edit with production values

# Build application
make build

# Create systemd service
sudo nano /etc/systemd/system/recallflow.service
```

`/etc/systemd/system/recallflow.service`:
```ini
[Unit]
Description=RecallFlow API Server
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/recallflow/backend
ExecStart=/home/ubuntu/recallflow/backend/bin/recallflow-api
Restart=on-failure
RestartSec=5s

Environment="PORT=8080"
EnvironmentFile=/home/ubuntu/recallflow/backend/.env

[Install]
WantedBy=multi-user.target
```

Enable and start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable recallflow
sudo systemctl start recallflow
sudo systemctl status recallflow
```

### 4. Configure NGINX

```bash
sudo nano /etc/nginx/sites-available/recallflow
```

```nginx
server {
    listen 80;
    server_name api.recallflow.ai;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/recallflow /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 5. Set Up SSL

```bash
sudo certbot --nginx -d api.recallflow.ai
```

### 6. Database Migration

```bash
cd /home/ubuntu/recallflow/backend
make migrate
```

## Frontend Deployment (Vercel)

### 1. Connect Repository

1. Go to [vercel.com](https://vercel.com)
2. Click "New Project"
3. Import your GitHub repository
4. Select `frontend` as root directory

### 2. Configure Environment Variables

In Vercel dashboard, add:
```
NEXT_PUBLIC_API_URL=https://api.recallflow.ai/api/v1
```

### 3. Deploy

Vercel automatically deploys on push to main branch.

### 4. Configure Custom Domain

1. Add domain in Vercel dashboard
2. Update DNS records:
   - `CNAME` record: `www` → `cname.vercel-dns.com`
   - `A` record: `@` → Vercel IP

## Twilio Configuration

### 1. Buy Phone Number

1. Go to Twilio Console
2. Buy a phone number with Voice + SMS capabilities

### 2. Configure Webhooks

Voice webhooks:
- When a call comes in: `https://api.recallflow.ai/api/v1/webhooks/twilio/voice` (HTTP POST)
- Status callback URL: `https://api.recallflow.ai/api/v1/webhooks/twilio/status` (HTTP POST)

SMS webhooks:
- When a message comes in: `https://api.recallflow.ai/api/v1/webhooks/twilio/sms` (HTTP POST)

## Monitoring & Maintenance

### Application Logs

```bash
# View live logs
sudo journalctl -u recallflow -f

# View last 100 lines
sudo journalctl -u recallflow -n 100
```

### Database Backups

Set up automated RDS snapshots in AWS Console (daily recommended).

### Health Checks

Set up CloudWatch alarms for:
- EC2 CPU > 80%
- RDS connections > 80%
- Application health endpoint

Health check endpoint: `https://api.recallflow.ai/health`

### Updates

```bash
# SSH into server
cd /home/ubuntu/recallflow

# Pull latest code
git pull origin main

# Rebuild and restart
cd backend
make build
sudo systemctl restart recallflow
```

## Cost Estimate (MVP)

Monthly AWS costs:
- EC2 t3.small: ~$15
- RDS db.t3.micro: ~$15
- ElastiCache cache.t3.micro: ~$12
- Data transfer: ~$5
- Total: **~$50/month**

Additional:
- Vercel: Free (Hobby plan)
- Twilio: Pay-as-you-go
- Domain: ~$12/year

## Security Checklist

- [ ] Enable AWS WAF on API
- [ ] Set up fail2ban on EC2
- [ ] Enable CloudWatch monitoring
- [ ] Rotate database credentials monthly
- [ ] Enable RDS encryption at rest
- [ ] Use AWS Secrets Manager for sensitive data
- [ ] Set up automated security updates
- [ ] Configure CORS properly
- [ ] Implement rate limiting

## Scaling Strategy

When you reach 50+ customers:

1. **Database**: Upgrade to `db.t3.small` with read replica
2. **API**: Add Application Load Balancer + Auto Scaling Group
3. **Cache**: Upgrade Redis to `cache.t3.small` with clustering
4. **CDN**: Add CloudFront for static assets
5. **Monitoring**: Add Datadog or New Relic

## Disaster Recovery

1. **Database**: Daily automated snapshots
2. **Application**: Store code in Git
3. **Secrets**: Backup to AWS Secrets Manager
4. **Documentation**: Keep runbooks updated

## Support

For deployment issues:
- Check logs: `sudo journalctl -u recallflow -f`
- Verify environment variables
- Test database connectivity
- Confirm Twilio webhooks

---

**Ready to deploy?** Start with backend setup, then frontend, then configure Twilio webhooks.
