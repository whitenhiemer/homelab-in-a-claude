# monitoring

Deploy Prometheus + Grafana + Alertmanager via Docker Compose on a Debian LXC. Optionally includes node-exporter and Proxmox PVE exporter.

## What it deploys

| Service | Port | Purpose |
|---|---|---|
| Prometheus | 9090 | Metrics scrape engine + TSDB |
| Grafana | 3000 | Dashboards |
| Alertmanager | 9093 | Alert routing (Discord, etc.) |
| node-exporter | 9100 | Host system metrics |
| pve-exporter | 9221 | Proxmox API metrics (optional) |

## Usage

```bash
ansible-playbook playbook.yml \
  -i inventory.yml \
  --extra-vars "grafana_password=secure_password"

# With Discord alerts and Proxmox scraping:
ansible-playbook playbook.yml \
  -i inventory.yml \
  --extra-vars "
    grafana_password=secure_password
    discord_webhook=https://discord.com/api/webhooks/...
    pve_host=192.168.86.29
    pve_user=monitoring@pve
    pve_token_name=prometheus
    pve_token_value=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  "
```

## Proxmox API token setup

```bash
# On Proxmox host (or via UI: Datacenter > Permissions > API Tokens)
pveum user add monitoring@pve
pveum aclmod / -user monitoring@pve -role PVEAuditor
pveum user token add monitoring@pve prometheus --privsep=0
```

## Files

```
monitoring/
├── playbook.yml                    # Ansible playbook
├── docker-compose.yml.j2           # Stack definition
├── prometheus.yml.j2               # Scrape targets
├── alertmanager.yml.j2             # Alert routing
└── README.md
```

## Post-deploy

1. Open Grafana at `http://<ip>:3000` (admin / your password)
2. Dashboards are not pre-loaded — import from grafana.com:
   - Node Exporter Full: ID 1860
   - Alertmanager: ID 9578
3. Add more scrape targets by editing `prometheus.yml.j2` and re-running
