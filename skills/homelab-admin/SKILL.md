---
name: homelab-admin
description: >
  Administers Brandon's Proxmox homelab at woodhead.tech. Activate when the user
  asks to "check service status", "restart the arr stack", "what's running on
  proxmox", "restart a service", "show cluster status", "check k8s",
  "start or stop a service group", "view logs for X", "what's the status of Y",
  "deploy a playbook", "SSH into a host", or any question about the state of
  infrastructure at woodhead.tech.
---

# homelab-admin

Skill for administering Brandon's Proxmox homelab running at woodhead.tech. This
covers: service group lifecycle, SSH access, Kubernetes (Talos) cluster operations,
Proxmox VM/LXC management, Ansible playbook execution, and log access.

All infrastructure data — hosts, IPs, ports, service groups, Makefile targets — is
in the reference files. Read them before acting:

- `references/infrastructure.md` — host inventory, service groups, K8s topology, URL/port map
- `references/operations.md` — Makefile targets, Ansible playbooks, SSH patterns, K8s commands, log access

---

## Working directory

All `make` and `ansible-playbook` commands run from the repo root:

```
~/Workspace/proxmox_kubernetes_cluster/
```

The Ansible working directory for direct playbook invocations is:

```
~/Workspace/proxmox_kubernetes_cluster/ansible/
```

---

## Service group management

Service groups are defined in `ansible/vars/service_groups.yml`. Groups are the
primary unit of lifecycle management — use them instead of touching individual VMs.

**Check status of all groups:**
```bash
cd ~/Workspace/proxmox_kubernetes_cluster && make group-status
```

**Start a group:**
```bash
make group-start GROUP=<name>
```

**Stop a group:**
```bash
make group-stop GROUP=<name>
```

Valid group names: `core`, `storage`, `security`, `home`, `media`, `observability`,
`apps`, `infra`, `sdr`, `special`, `k8s`. See `references/infrastructure.md` for
each group's members and dependency rules.

Safety rules enforced by the playbooks:
- `always_on` groups (`core`, `storage`) refuse stop requests
- `required_by` relationships block stop if dependent groups are running
- `hardware_bound` groups (`special`) are excluded from bulk operations

When a user asks to stop a group, check its `required_by` chain in
`references/infrastructure.md` first and warn if stopping it would break a
running dependent.

---

## SSH access

All Proxmox nodes and LXC containers use `root` with `~/.ssh/id_ansible`. Raspberry
Pi devices and Klipper hosts use `bwoodwar` with `~/.ssh/id_ed25519`.

**SSH to a Proxmox node:**
```bash
ssh -i ~/.ssh/id_ansible root@192.168.86.29   # pve1 / thinkcentre1
```

**SSH to an LXC:**
```bash
ssh -i ~/.ssh/id_ansible -o StrictHostKeyChecking=accept-new root@192.168.86.22   # arr-stack
```

**SSH to a Raspberry Pi:**
```bash
ssh -i ~/.ssh/id_ed25519 bwoodwar@192.168.86.131   # piboard
```

For all host IPs, consult `references/infrastructure.md`.

---

## Kubernetes cluster operations

The Talos cluster uses `talosctl` and `kubectl`. Generated configs live at
`talos/_out/` (gitignored — must be present locally to operate the cluster).

**Set environment before any K8s command:**
```bash
export KUBECONFIG=~/Workspace/proxmox_kubernetes_cluster/talos/_out/kubeconfig
export TALOSCONFIG=~/Workspace/proxmox_kubernetes_cluster/talos/_out/talosconfig
```

**Check cluster health:**
```bash
cd ~/Workspace/proxmox_kubernetes_cluster && make health
# or directly:
talosctl --talosconfig talos/_out/talosconfig health
```

**Get nodes:**
```bash
kubectl get nodes -o wide
```

**Get all pods across namespaces:**
```bash
kubectl get pods -A
```

**Apply base manifests (namespaces, MetalLB):**
```bash
make k8s-base
make k8s-base-metallb   # includes MetalLB IP pool
```

**Fetch a fresh kubeconfig:**
```bash
make kubeconfig
```

Cluster topology: VIP `.100`, control plane `.101` (tower1), workers `.111`
(thinkcentre2), `.112` (thinkcentre3), `.113` (zotac). Kubernetes v1.31.0, Talos
v1.12.5. MetalLB pool: `192.168.86.150–199`.

---

## Ansible playbook execution

Deploy or reconfigure a service by running its Makefile target. Most targets wrap
an Ansible playbook in `ansible/playbooks/setup-<service>.yml`.

**General pattern:**
```bash
cd ~/Workspace/proxmox_kubernetes_cluster && make <target>
```

Some targets require env vars — the Makefile will error and print usage if they
are missing. Check `references/operations.md` for required vars per target.

**Day-2 patching:**
```bash
make patch-proxmox    # Proxmox nodes, serial
make patch-lxc        # All LXC containers (apt)
make patch-docker     # Pull latest images + restart all stacks
make patch-pi         # Raspberry Pi devices
```

**Run a playbook directly** (when you need `--limit`, `--tags`, or `--check`):
```bash
cd ~/Workspace/proxmox_kubernetes_cluster/ansible && \
  ansible-playbook playbooks/setup-traefik.yml --check
```

---

## Proxmox VM/LXC management

For ad-hoc VM operations not covered by Ansible, use `pvesh` or `qm`/`pct` on the
Proxmox node that hosts the VM. To determine which node hosts a given VMID, run
`make group-status` or query the cluster API:

```bash
ssh -i ~/.ssh/id_ansible root@192.168.86.29 \
  'pvesh get /cluster/resources --type vm --output-format json' | \
  python3 -m json.tool | grep -A3 '"vmid": 202'
```

**Start/stop a specific VM or LXC:**
```bash
# On the hosting node:
qm start <vmid>     # VM
pct start <vmid>    # LXC
qm stop <vmid>
pct stop <vmid>
```

Prefer service group targets (`make group-start/stop`) for anything managed by
Ansible. Direct `qm`/`pct` is appropriate for one-off recovery or inspection.

---

## Traefik and routing

Traefik LXC is at `192.168.86.20`. All HTTPS traffic for `*.woodhead.tech` enters
here. Static config: `ansible/files/traefik/traefik.yml`. Per-service routes:
`ansible/files/traefik/dynamic/*.yml`.

**Redeploy Traefik config:**
```bash
make traefik    # requires scripts/ddns/cloudflare.env to exist
```

**SSH to Traefik to inspect logs or reload:**
```bash
ssh -i ~/.ssh/id_ansible root@192.168.86.20
# Traefik runs as a systemd service: systemctl status traefik
```

All services use Cloudflare DNS-01 for TLS. Services behind Authentik use the
`authentik@file` forwardAuth middleware.

---

## Monitoring and logs

Monitoring LXC is at `192.168.86.25`. Grafana: `grafana.woodhead.tech:3000`,
Prometheus: `prometheus.woodhead.tech:9090`, Alertmanager: `alertmanager.woodhead.tech:9093`.

**Docker service logs** (services that run via Docker Compose):
```bash
ssh -i ~/.ssh/id_ansible root@<lxc-ip>
docker compose logs -f           # from the stack directory
# Stack dirs are usually /opt/<service-name>/
```

**systemd service logs:**
```bash
journalctl -u <service-name> -f
```

**Talos node logs:**
```bash
talosctl --talosconfig talos/_out/talosconfig logs -n 192.168.86.101 kubelet
talosctl --talosconfig talos/_out/talosconfig dmesg -n 192.168.86.101
```

**Kubernetes pod logs:**
```bash
kubectl logs -n <namespace> <pod-name> -f
kubectl logs -n <namespace> -l app=<label> --all-containers
```

---

## Verifying service health after changes

After starting a group or deploying a service, confirm it is actually healthy:

1. Run `make group-status` to verify the group shows RUNNING
2. SSH to the LXC and check the Docker Compose stack: `docker compose ps`
3. If any container is not in state `Up`, check its logs: `docker compose logs -f <service>`
4. For systemd-managed services (Traefik, piboard), check `systemctl status <service>` on the host
5. Hit the public URL to confirm Traefik is routing correctly

For Kubernetes workloads:
1. `kubectl get pods -A | grep -Ev 'Running|Completed'` — surface anything not healthy
2. `kubectl describe pod -n <ns> <name>` — check Events section for root cause
3. `talosctl --talosconfig talos/_out/talosconfig health` — confirm Talos components are healthy

---

## Security model

All public `*.woodhead.tech` services terminate TLS at Traefik (`192.168.86.20`).
Most admin interfaces require Authentik SSO (`auth.woodhead.tech`, `192.168.86.28:9000`).
Services without Authentik either have their own auth (Plex, Jellyfin, Home Assistant, Mailcow)
or are intentionally public (Libby Alert at `alert.woodhead.tech`).

The Pwnagotchi web UI (`192.168.86.38`) is accessible via WireGuard only — avoid
exposing it directly. Access via `192.168.86.39` (WireGuard LXC) tunnel.

The ARR stack LXC (`192.168.86.22`) routes its download traffic through a WireGuard
tunnel configured at deploy time via `WG_PRIVATE_KEY`. If the VPN tunnel drops,
downloads stop but the web UIs remain accessible.

---

## Common workflows

**"Is the arr stack running?"**
1. `make group-status` — check the `media` group; expect RUNNING with all 3 members
2. If PARTIAL or STOPPED, SSH to `192.168.86.22`:
   ```bash
   ssh -i ~/.ssh/id_ansible root@192.168.86.22
   docker compose ps
   docker compose logs -f
   ```
3. To restart the group: `make group-start GROUP=media`

**"Restart the media group"**
1. Check `required_by` in `references/infrastructure.md` — `media` has none, safe to stop
2. `make group-stop GROUP=media`
3. `make group-start GROUP=media`
4. Verify: `make group-status`

**"Stop the monitoring stack"**
1. Check `required_by` — `observability` group has none, safe to stop
2. `make group-stop GROUP=observability`
3. Note: the `monitoring-stack` LXC also hosts the docs, resume, homelab, and landing static sites on ports 8081-8084 — those will also go down

**"What's the status of the K8s cluster?"**
1. Set env vars:
   ```bash
   export KUBECONFIG=~/Workspace/proxmox_kubernetes_cluster/talos/_out/kubeconfig
   export TALOSCONFIG=~/Workspace/proxmox_kubernetes_cluster/talos/_out/talosconfig
   ```
2. `make health` — Talos component health
3. `kubectl get nodes -o wide` — node ready status + IPs
4. `kubectl get pods -A | grep -Ev 'Running|Completed'` — surface problem pods
5. For a specific failing pod: `kubectl describe pod -n <ns> <name>`

**"Deploy a service update"**
1. Run the appropriate `make <service>` target from `references/operations.md`
2. For `traefik`, `monitoring`, `arr-stack` — check required env vars before running
3. Ansible is idempotent; safe to re-run without side effects
4. After deploy, verify the service is up via its URL or `make group-status`

**"Update Docker images for all services"**
```bash
make patch-docker
```
This runs `docker pull` + `docker compose up -d` on all LXCs. Run during a maintenance window; brief downtime per service is expected.

**"Apply OS patches to all LXCs"**
```bash
make patch-lxc
```
Runs `apt-get upgrade` across all LXC containers. Non-destructive; services continue running.

**"Check Grafana / monitoring"**
- Public URL: `https://grafana.woodhead.tech` (requires Authentik login)
- Direct internal: `http://192.168.86.25:3000`
- Prometheus: `http://192.168.86.25:9090`
- Alertmanager: `http://192.168.86.25:9093`
- All three services run via Docker Compose on the `monitoring-stack` LXC

**"Redeploy Traefik config after editing a dynamic route"**
1. Edit the relevant file in `ansible/files/traefik/dynamic/`
2. `make traefik` — redeploys Traefik config via Ansible (requires `scripts/ddns/cloudflare.env`)
3. Alternatively, SSH to `192.168.86.20` and restart: `systemctl restart traefik`

**"Add a new route to Traefik"**
1. Create a new `.yml` file in `ansible/files/traefik/dynamic/`
2. Follow the existing pattern: router + service block, `tls.certResolver: cloudflare`, add `authentik@file` middleware if auth is needed
3. `make traefik` to deploy

**"Check TrueNAS"**
- Public URL: `https://nas.woodhead.tech` (Authentik-protected)
- Direct: `https://192.168.86.40` (self-signed cert; expect browser warning)
- TrueNAS is in the `storage` group (`always_on: true`) — do not stop it while `media` is running
