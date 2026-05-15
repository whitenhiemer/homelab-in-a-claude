# ShopStack Recipes

Managed business infrastructure for small businesses — vet clinics, retail shops, salons.
Deploys a complete stack (website, email, file storage, database) on a single server.

## What gets deployed

| Service | Role |
|---------|------|
| PostgreSQL | Shared database backend |
| Traefik | Reverse proxy + automatic TLS via Cloudflare |
| WireGuard | Encrypted remote management tunnel |
| Mailcow | Email server (SMTP/IMAP, webmail, spam filter) |
| Nextcloud | File storage (Dropbox replacement) |
| Uptime Kuma | Internal uptime monitoring |

## Platforms

ShopStack runs identically on all three — the Ansible playbooks are platform-agnostic.

| Platform | Hardware | Setup |
|----------|----------|-------|
| **On-premises** | Beelink EQ12 mini PC (~$200, customer buys) | Debian 12, no Terraform needed |
| **AWS** | EC2 t3.large (~$60/mo on-demand) | `terraform/aws/` provisions instance + Elastic IP + Security Group |
| **GCP** | Compute e2-standard-2 (~$50/mo) | `terraform/gcp/` provisions instance + Static IP + Firewall Rules |

## Quick start

### On-premises (Beelink)

```bash
# 1. Customer installs Debian 12 on the Beelink, assigns static LAN IP
# 2. You port-forward 80, 443, 25, 465, 587, 993, 995, 51820/udp on their router

# 3. Create inventory
cat > inventory.ini <<EOF
[shopstack]
shopstack ansible_host=<beelink-ip> ansible_user=root ansible_ssh_private_key_file=~/.ssh/id_ansible
EOF

# 4. Run
ansible-playbook ansible/shopstack.yml \
  -i inventory.ini \
  --extra-vars "@profiles/on-prem.yml" \
  --extra-vars "domain=yourclinic.com" \
  --extra-vars "cf_api_token=<token>" \
  --extra-vars "acme_email=admin@yourclinic.com" \
  --extra-vars "postgres_password=<secret>" \
  --extra-vars "nextcloud_db_pass=<secret>" \
  --extra-vars "nextcloud_admin_pass=<secret>" \
  --extra-vars "mailcow_domain=yourclinic.com"
```

### AWS

```bash
# 1. Provision infrastructure
cd terraform/aws
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars with your values
terraform init && terraform apply

# 2. terraform outputs the inventory line — paste it into inventory.ini

# 3. Run Ansible with the aws profile
ansible-playbook ansible/shopstack.yml \
  -i inventory.ini \
  --extra-vars "@profiles/aws.yml" \
  --extra-vars "domain=yourclinic.com" \
  # ... same required vars as above
```

### GCP

```bash
cd terraform/gcp
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars
terraform init && terraform apply

ansible-playbook ansible/shopstack.yml \
  -i inventory.ini \
  --extra-vars "@profiles/gcp.yml" \
  --extra-vars "domain=yourclinic.com" \
  # ... same required vars
```

## Directory layout

```
shopstack/
  ansible/
    shopstack.yml          # Master playbook — runs all services in order
    mailcow/               # Email server
    nextcloud/             # File storage
    postgres/              # Database
    uptime-kuma/           # Monitoring
  terraform/
    aws/                   # EC2 + Elastic IP + Security Group
    gcp/                   # Compute + Static IP + Firewall Rules
  profiles/
    on-prem.yml            # Vars for Beelink mini PC
    aws.yml                # Vars for AWS EC2
    gcp.yml                # Vars for GCP Compute
```

## DNS records required

After deploying, create these Cloudflare A records pointing to the server's public IP.
All should be **DNS-only (grey cloud)** except the website:

| Subdomain | Purpose |
|-----------|---------|
| `mail.domain.com` | Mailcow webmail + autodiscover |
| `files.domain.com` | Nextcloud |
| `wg.domain.com` | WireGuard endpoint (DNS-only, never proxy) |
| `www.domain.com` | Customer website |

Also add MX, SPF, DKIM, and DMARC records — Mailcow's admin panel generates these after setup.

## Per-client pricing reference

| Deployment | Infrastructure cost | Managed fee | Margin |
|------------|--------------------|--------------|-|
| On-prem ($500 setup) | $0/mo (customer owns hardware) | $200/mo | ~$200/mo |
| AWS | ~$60/mo | $250/mo | ~$190/mo |
| GCP | ~$50/mo | $250/mo | ~$200/mo |
