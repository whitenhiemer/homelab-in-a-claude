# home-assistant

Wire Home Assistant OS (HAOS) into your homelab's Traefik reverse proxy so it's reachable at `https://home.<domain>`.

HAOS manages itself — this playbook does not touch the VM directly. It only deploys the Traefik dynamic route and prints the manual step required inside HA.

## What it does

1. Deploys `/etc/traefik/dynamic/homeassistant.yml` on the Traefik LXC
2. Prints the `trusted_proxies` block you must add to HA's `configuration.yaml`

## Usage

```bash
ansible-playbook playbook.yml \
  -i inventory.yml \
  --extra-vars "domain=example.com"
```

## Variables

| Variable | Default | Description |
|---|---|---|
| `ha_ip` | `192.168.86.41` | Home Assistant VM IP |
| `domain` | **required** | Base domain (e.g. `example.com`) |
| `traefik_ip` | `192.168.86.20` | Traefik LXC IP (printed in instructions) |

## Prerequisites

1. HAOS VM created and running (initial setup wizard completed)
2. Traefik LXC deployed and running
3. Cloudflare DNS record for `home.<domain>` pointing to your WAN IP (or let Traefik handle it via DDNS)

## Required manual step

After the playbook runs, add this to `/config/configuration.yaml` in Home Assistant:

```yaml
http:
  use_x_forwarded_for: true
  trusted_proxies:
    - <traefik_ip>
```

Then restart Home Assistant: **Settings → System → Restart Home Assistant**.

How to edit `configuration.yaml`:
- **Option 1:** Install the *Studio Code Server* add-on (Settings → Add-ons → Add-on Store)
- **Option 2:** Install the *Terminal & SSH* add-on and use the built-in terminal
- **Option 3:** HA UI → Settings → System → Edit configuration.yaml

## Files

```
home-assistant/
├── playbook.yml            # Ansible playbook
├── homeassistant.yml.j2    # Traefik dynamic route template
└── README.md
```
