# traefik

Install and configure Traefik v3 as the reverse proxy / TLS terminator for all homelab services.

## What it does

- Installs Traefik binary from GitHub releases
- Creates `traefik` system user with minimal privileges
- Deploys static config (entrypoints, ACME resolver, file provider)
- Deploys dynamic route configs from `dynamic/*.yml`
- Configures Let's Encrypt DNS-01 via Cloudflare
- Runs as a systemd service with `CAP_NET_BIND_SERVICE` so it can bind to 80/443 without root

## Files

```
traefik/
├── playbook.yml          # Ansible playbook
├── traefik.yml.j2        # Jinja2 template for Traefik static config
├── dynamic/              # Drop per-service route configs here (*.yml)
└── README.md
```

## Usage

```bash
ansible-playbook playbook.yml \
  -i inventory.yml \
  --extra-vars "cf_api_token=your_token acme_email=you@example.com domain=example.com"
```

## Inventory entry

```yaml
traefik:
  hosts:
    traefik:
      ansible_host: 192.168.86.20
      ansible_user: root
```

## Dynamic routes

Add files to `dynamic/` that follow this pattern:

```yaml
http:
  routers:
    myservice:
      rule: "Host(`myservice.example.com`)"
      entryPoints: [websecure]
      service: myservice
      tls:
        certResolver: cloudflare

  services:
    myservice:
      loadBalancer:
        servers:
          - url: "http://192.168.86.50:8080"
```

## Cloudflare API token

Create at dash.cloudflare.com → My Profile → API Tokens → Create Token → Edit zone DNS. Scope to your zone only.
