# Operations Reference

All commands run from `~/Workspace/proxmox_kubernetes_cluster/` unless noted.

---

## Makefile Targets

### Infrastructure provisioning

| Target             | What it does                                                    |
|--------------------|-----------------------------------------------------------------|
| `make setup`       | Verify/configure Proxmox hosts base (run once after install)   |
| `make prepare`     | Download Talos ISO to Proxmox hosts                            |
| `make prepare-truenas` | Download TrueNAS Scale ISO to Proxmox hosts               |
| `make ddns`        | Deploy Cloudflare DDNS updater (required before TLS resolves)  |
| `make init`        | `terraform init` — initialize providers                        |
| `make plan`        | `terraform plan` — preview all VM/LXC changes                  |
| `make apply`       | `terraform apply` — create/update all VMs and LXCs             |
| `make apply-truenas` | Create/update TrueNAS VM only                                |
| `make apply-homeassistant` | Create/update Home Assistant VM only                   |
| `make plan-lxc`    | Preview LXC changes only                                       |
| `make apply-lxc`   | Create/update LXC containers only                              |
| `make harden`      | Apply security hardening to Proxmox hosts                      |
| `make destroy`     | Tear down all VMs (prompts for confirmation — destructive)     |
| `make clean`       | Remove `talos/_out/` generated configs (does not destroy VMs)  |

### Service deployment (Ansible)

| Target                   | Playbook                            | Required env vars / notes                                     |
|--------------------------|-------------------------------------|---------------------------------------------------------------|
| `make traefik`           | `setup-traefik.yml`                 | Reads CF_API_TOKEN from `scripts/ddns/cloudflare.env`; run `make ddns` first |
| `make recipe-site`       | `setup-recipe-site.yml`             | —                                                             |
| `make arr-stack`         | `setup-arr-stack.yml`               | `WG_PRIVATE_KEY=<key>` required                              |
| `make plex`              | `setup-plex.yml`                    | —                                                             |
| `make jellyfin`          | `setup-jellyfin.yml`                | —                                                             |
| `make monitoring`        | `setup-monitoring.yml`              | Optional: `DISCORD_WEBHOOK`, `GRAFANA_PASSWORD`, `PVE_USER`, `PVE_TOKEN_NAME`, `PVE_TOKEN_VALUE`, `DEXCOM_USERNAME`, `DEXCOM_PASSWORD`, `HA_GLUCOSE_WEBHOOK`, `HA_TOKEN` |
| `make openclaw`          | `setup-openclaw.yml`                | —                                                             |
| `make authentik`         | `setup-authentik.yml`               | `AUTHELIA_ADMIN_PASSWORD=...` (if still applicable)           |
| `make wireguard`         | `setup-wireguard.yml`               | —                                                             |
| `make homeassistant`     | `setup-homeassistant.yml`           | Optional: `HA_TOKEN=<long-lived-token>`                      |
| `make beardie`           | `setup-homeassistant.yml --tags packages` | Gutgrinda enclosure automation; optional: `GUTGRINDA_DISCORD_WEBHOOK`, `HA_TOKEN` |
| `make truenas`           | `setup-truenas.yml`                 | `TRUENAS_PASSWORD=<root-password>` required                  |
| `make libby-alert`       | `setup-libby-alert.yml`             | `DISCORD_WEBHOOK` and/or Twilio vars required: `TWILIO_SID`, `TWILIO_TOKEN`, `TWILIO_FROM`, `ALERT_PHONES` |
| `make kanboard`          | `setup-kanboard.yml`                | —                                                             |
| `make sdr`               | `setup-sdr.yml`                     | —                                                             |
| `make pxe`               | `setup-pxe.yml`                     | —                                                             |
| `make zigbee2mqtt`       | `setup-zigbee2mqtt.yml`             | —                                                             |
| `make pwnagotchi`        | `setup-pwnagotchi.yml`              | Optional: `WEB_PASSWORD=<password>`                          |
| `make claude-os`         | `setup-claude-os.yml`               | Optional: `OPENAI_API_KEY=sk-...` or `INSTALL_OLLAMA=true`  |
| `make mailserver`        | `setup-mailserver.yml`              | Optional: `MAILGUN_USER`, `MAILGUN_PASSWORD`                 |

### Service group lifecycle

| Target                          | What it does                                              |
|---------------------------------|-----------------------------------------------------------|
| `make group-status`             | Show running status of all service groups                 |
| `make group-start GROUP=<name>` | Start all members of a named group                        |
| `make group-stop GROUP=<name>`  | Stop all members of a named group                         |

Valid group names: `core`, `storage`, `security`, `home`, `media`, `observability`, `apps`, `infra`, `sdr`, `special`, `k8s`

### Kubernetes

| Target                   | What it does                                                     |
|--------------------------|------------------------------------------------------------------|
| `make bootstrap`         | Generate Talos configs and bootstrap cluster (destructive, once) |
| `make kubeconfig`        | Fetch kubeconfig from running cluster → `talos/_out/kubeconfig`  |
| `make health`            | `talosctl health` — cluster node + component health check        |
| `make k8s-base`          | Apply base manifests (namespaces)                                |
| `make k8s-base-metallb`  | Apply base manifests + MetalLB IP pool                          |

Bootstrap requires env vars:
```bash
export CLUSTER_VIP="192.168.86.100"
export CONTROLPLANE_IPS="192.168.86.101"
export WORKER_IPS="192.168.86.111,192.168.86.112,192.168.86.113"
```

### Patching

| Target               | What it does                                              |
|----------------------|-----------------------------------------------------------|
| `make patch-proxmox` | apt upgrade on all Proxmox nodes, serial (one at a time) |
| `make patch-lxc`     | apt upgrade on all LXC containers                        |
| `make patch-docker`  | `docker pull` + restart all Docker Compose stacks        |
| `make patch-pi`      | apt upgrade on Raspberry Pi devices (piboard, etc.)      |

### Documentation / static sites

| Target              | What it does                                           |
|---------------------|--------------------------------------------------------|
| `make docs-build`   | Build Docusaurus docs site → `docs-site/build/`       |
| `make docs-dev`     | Start Docusaurus dev server (hot reload)              |
| `make resume-build` | Build Hugo resume site → `resume-site/public/`        |

---

## SSH Access Patterns

### Default: Proxmox nodes and LXC containers

```bash
ssh -i ~/.ssh/id_ansible -o StrictHostKeyChecking=accept-new root@<ip>
```

### Raspberry Pi and special LXCs (claude-os, pwnagotchi)

```bash
ssh -i ~/.ssh/id_ed25519 -o StrictHostKeyChecking=accept-new <user>@<ip>
```

User is `bwoodwar` for Pi devices, `root` for LXCs.

### Quick reference

```bash
# Proxmox nodes
ssh -i ~/.ssh/id_ansible root@192.168.86.29   # pve1 / thinkcentre1
ssh -i ~/.ssh/id_ansible root@192.168.86.30   # pve2 / thinkcentre2
ssh -i ~/.ssh/id_ansible root@192.168.86.31   # pve3 / thinkcentre3
ssh -i ~/.ssh/id_ansible root@192.168.86.130  # tower1
ssh -i ~/.ssh/id_ansible root@192.168.86.147  # zotac

# Key LXCs
ssh -i ~/.ssh/id_ansible root@192.168.86.20   # traefik
ssh -i ~/.ssh/id_ansible root@192.168.86.22   # arr-stack
ssh -i ~/.ssh/id_ansible root@192.168.86.25   # monitoring
ssh -i ~/.ssh/id_ansible root@192.168.86.28   # authentik
ssh -i ~/.ssh/id_ansible root@192.168.86.34   # mailserver

# Pi devices
ssh -i ~/.ssh/id_ed25519 bwoodwar@192.168.86.131  # piboard
ssh -i ~/.ssh/id_ed25519 bwoodwar@192.168.86.136  # ender5pro
ssh -i ~/.ssh/id_ed25519 bwoodwar@192.168.86.138  # ender3
```

---

## Kubernetes Admin Commands

Set environment first:
```bash
export KUBECONFIG=~/Workspace/proxmox_kubernetes_cluster/talos/_out/kubeconfig
export TALOSCONFIG=~/Workspace/proxmox_kubernetes_cluster/talos/_out/talosconfig
```

### kubectl patterns

```bash
kubectl get nodes -o wide                              # node status + IPs
kubectl get pods -A                                    # all pods all namespaces
kubectl get pods -A | grep -Ev 'Running|Completed'    # surface unhealthy pods
kubectl describe node <name>                           # node detail + conditions
kubectl describe pod -n <ns> <name>                   # pod events + state
kubectl logs -n <ns> <pod> -f                         # stream pod logs
kubectl logs -n <ns> -l app=<label> --all-containers  # multi-pod logs by label
kubectl rollout restart deployment -n <ns> <name>     # restart a deployment
kubectl get events -n <ns> --sort-by=.lastTimestamp   # recent events
kubectl top nodes                                      # resource usage (needs metrics-server)
kubectl top pods -A
```

### talosctl patterns

```bash
# Health and status
talosctl --talosconfig talos/_out/talosconfig health
talosctl --talosconfig talos/_out/talosconfig -n 192.168.86.101 version
talosctl --talosconfig talos/_out/talosconfig nodes

# Logs
talosctl --talosconfig talos/_out/talosconfig logs -n <ip> kubelet
talosctl --talosconfig talos/_out/talosconfig logs -n <ip> etcd
talosctl --talosconfig talos/_out/talosconfig dmesg -n <ip>

# Services
talosctl --talosconfig talos/_out/talosconfig services -n <ip>

# Config apply (after editing patches)
talosctl --talosconfig talos/_out/talosconfig apply-config -n <ip> -f talos/_out/<node>.yaml

# Upgrade
talosctl --talosconfig talos/_out/talosconfig upgrade -n <ip> --image ghcr.io/siderolabs/installer:<version>
```

Node IPs: `192.168.86.101` (cp-0), `192.168.86.111` (worker-0), `192.168.86.112` (worker-1), `192.168.86.113` (worker-2)

---

## Service Restart Procedures

### Docker Compose services (arr-stack, monitoring, openclaw, etc.)

```bash
ssh -i ~/.ssh/id_ansible root@<lxc-ip>
cd /opt/<service>/          # typical stack dir
docker compose pull         # update images
docker compose up -d        # restart with latest
docker compose ps           # verify all containers up
docker compose logs -f      # stream logs
```

Common stack directories:
- arr-stack: `/opt/arr/` (or similar; check `ls /opt/`)
- monitoring: `/opt/monitoring/`
- openclaw: `/opt/openclaw/`
- authentik: `/opt/authentik/`
- kanboard: `/opt/kanboard/`
- claude-os: `/opt/claude-os/`

### systemd services (traefik, piboard, etc.)

```bash
ssh -i ~/.ssh/id_ansible root@<ip>
systemctl status <service>
systemctl restart <service>
journalctl -u <service> -f --no-pager
```

Traefik runs as `systemd` service on `192.168.86.20`:
```bash
systemctl status traefik
systemctl restart traefik
journalctl -u traefik -f
```

### Proxmox VM/LXC (via node)

```bash
# Determine hosting node first via group-status or pvesh
ssh -i ~/.ssh/id_ansible root@<node-ip>

qm start <vmid>    # start VM
qm stop <vmid>     # graceful stop
qm reset <vmid>    # hard reset
qm status <vmid>

pct start <vmid>   # start LXC
pct stop <vmid>    # stop LXC
pct status <vmid>
pct enter <vmid>   # shell into LXC (from node)
```

---

## Log Access by Service Type

### Docker Compose services

```bash
ssh -i ~/.ssh/id_ansible root@<lxc-ip>
docker compose -f /opt/<stack>/docker-compose.yml logs -f
docker compose -f /opt/<stack>/docker-compose.yml logs --tail=100 <service-name>
```

### Proxmox node logs

```bash
ssh -i ~/.ssh/id_ansible root@<node-ip>
journalctl -f                        # system journal
journalctl -u pveproxy -f            # PVE web UI
journalctl -u pvedaemon -f           # PVE daemon
journalctl -u pve-cluster -f         # cluster service
tail -f /var/log/syslog
```

### Talos node logs

```bash
talosctl --talosconfig talos/_out/talosconfig logs -n <ip> kubelet
talosctl --talosconfig talos/_out/talosconfig logs -n <ip> etcd      # cp only
talosctl --talosconfig talos/_out/talosconfig logs -n <ip> containerd
talosctl --talosconfig talos/_out/talosconfig logs -n <ip> controller-runtime
```

### Kubernetes pod logs

```bash
kubectl logs -n <namespace> <pod-name> -f
kubectl logs -n <namespace> <pod-name> --previous    # previous container instance
kubectl logs -n <namespace> -l <label-selector> --all-containers=true
```

### Traefik access/error logs

```bash
ssh -i ~/.ssh/id_ansible root@192.168.86.20
journalctl -u traefik -f
# or if file-based logging is configured:
tail -f /var/log/traefik/access.log
```

### Monitoring stack

```bash
ssh -i ~/.ssh/id_ansible root@192.168.86.25
# Prometheus
docker compose logs -f prometheus
# Grafana
docker compose logs -f grafana
# Alertmanager
docker compose logs -f alertmanager
```

---

## Ansible Inventory and Vars

Inventory: `ansible/inventory/hosts.yml`
Service group vars: `ansible/vars/service_groups.yml`
Traefik static config: `ansible/files/traefik/traefik.yml`
Traefik dynamic configs: `ansible/files/traefik/dynamic/*.yml`

Run a playbook with limit:
```bash
cd ansible && ansible-playbook playbooks/<name>.yml --limit <host-or-group>
```

Run with check (dry-run):
```bash
ansible-playbook playbooks/<name>.yml --check --diff
```

Run with tags:
```bash
ansible-playbook playbooks/setup-homeassistant.yml --tags packages
```

---

## Deployment Order (fresh install)

1. `make setup` + `make prepare` — Proxmox base config + ISOs
2. `make ddns` — Cloudflare DDNS (DNS resolves before TLS certs)
3. `make init` → `make apply` — Terraform provisions all VMs + LXCs
4. `make traefik` — Ingress must be up before any HTTP service is reachable
5. Service playbooks: `arr-stack`, `plex`, `jellyfin`, `monitoring`, `openclaw`, `authentik`, `wireguard`, etc.
6. `make bootstrap` → `make kubeconfig` → `make k8s-base` — Talos K8s cluster

Traefik dependency: all web-reachable services require Traefik to be running first.
K8s dependency: cluster VMs must exist from `make apply` before `make bootstrap`.
