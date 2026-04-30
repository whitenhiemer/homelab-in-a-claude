# authelia (Authentik)

Deploy [Authentik](https://goauthentik.io) identity provider via Docker Compose on a Debian LXC. Provides SSO for all homelab services via Traefik's `forwardAuth` middleware.

## Stack

| Container | Purpose |
|---|---|
| `postgres:16-alpine` | Database |
| `redis:alpine` | Session cache + task queue |
| `authentik-server` | Web UI + API (port 9000/9443) |
| `authentik-worker` | Background tasks (email, notifications) |

The server uses host networking so other containers and LXCs can reach it without hairpin NAT issues.

## Usage

```bash
ansible-playbook playbook.yml \
  -i inventory.yml \
  --extra-vars "domain=example.com"
```

PostgreSQL password and Authentik secret key are auto-generated on first run and stored in `/opt/authentik/.env`. Re-running the playbook is safe — the `.env` is only written once.

## Post-deploy setup (required)

1. Browse to `https://auth.<domain>/if/flow/initial-setup/`
2. Create your admin account (email + password)
3. **Admin UI → Applications → Providers → Create → Proxy Provider**
   - Mode: Forward auth (domain level)
   - External host: `https://auth.<domain>`
4. **Admin UI → Applications → Create**
   - Bind the provider you just created
5. **Admin UI → Outposts → Edit the default outpost**
   - Add the forward auth application
6. Add to Traefik dynamic config to protect services:

```yaml
http:
  middlewares:
    authentik:
      forwardAuth:
        address: "http://<authentik-ip>:9000/outpost.goauthentik.io/auth/nginx"
        trustForwardHeader: true
        authResponseHeaders:
          - X-authentik-username
          - X-authentik-groups
```

## Files

```
authelia/
├── playbook.yml              # Ansible playbook
├── docker-compose.yml.j2     # Stack definition template
└── README.md
```

## Requirements

- LXC with at least 2 GB RAM (Authentik is Java-based, takes ~1.2 GB at idle)
- Docker nesting enabled on the LXC (`features: nesting=1` in Proxmox config)
