# Infrastructure Reference

Network: `192.168.86.0/24` | Domain: `woodhead.tech` | DNS: Cloudflare

---

## Proxmox Nodes

| Hostname | Alias        | IP              | User | Notes                                      |
|----------|-------------|-----------------|------|--------------------------------------------|
| pve1     | thinkcentre1 | 192.168.86.29   | root | Primary node                               |
| pve2     | thinkcentre2 | 192.168.86.30   | root | Hosts talos-worker-0 (K8s), SDR scanner   |
| pve3     | thinkcentre3 | 192.168.86.31   | root | RTL8188EUS WiFi dongle; hosts pwnagotchi  |
| tower1   | tower1       | 192.168.86.130  | root | Hosts TrueNAS VM, K8s control plane       |
| zotac    | zotac        | 192.168.86.147  | root | Hosts Zigbee2MQTT, talos-worker-2 (K8s)   |

SSH key: `~/.ssh/id_ansible`

---

## LXC Containers

| Ansible name      | VMID | IP              | User | Service                                           |
|-------------------|------|-----------------|------|---------------------------------------------------|
| traefik-gw        | 200  | 192.168.86.20   | root | Traefik reverse proxy (ingress for all services) |
| recipe-site       | 201  | 192.168.86.21   | root | Recipe site (Go app + nginx)                     |
| arr-stack         | 202  | 192.168.86.22   | root | ARR media stack (Sonarr, Radarr, Prowlarr, Bazarr, SABnzbd, Overseerr) via Docker Compose + WireGuard VPN |
| plex-server       | 203  | 192.168.86.23   | root | Plex Media Server (iGPU passthrough)             |
| jellyfin-server   | 204  | 192.168.86.24   | root | Jellyfin Media Server (iGPU passthrough)         |
| monitoring-stack  | 205  | 192.168.86.25   | root | Prometheus, Grafana, Alertmanager, docs/landing/resume/homelab sites |
| openclaw          | 206  | 192.168.86.26   | root | OpenClaw AI agent gateway                        |
| libby-alert       | 209  | 192.168.86.27   | root | Libby life alert QR website                      |
| authentik         | 207  | 192.168.86.28   | root | Authentik identity provider / SSO                |
| kanboard          | 211  | 192.168.86.33   | root | Kanboard project management                      |
| mailserver        | 212  | 192.168.86.34   | root | Mailcow email server (Mailgun relay)             |
| pxe-server        | 213  | 192.168.86.35   | root | PXE boot server (proxy-DHCP + TFTP + HTTP)      |
| zigbee2mqtt       | 214  | 192.168.86.36   | root | Zigbee2MQTT + Mosquitto (on zotac; Zigbee USB dongle) |
| claude-os         | 215  | 192.168.86.37   | root | Claude OS AI memory/knowledge system; key: id_ed25519 |
| pwnagotchi        | 216  | 192.168.86.38   | root | Pwnagotchi WiFi learning device (on pve3); key: id_ed25519 |
| wireguard         | 208  | 192.168.86.39   | root | WireGuard VPN tunnel                             |
| sdr-scanner       | 210  | 192.168.86.32   | root | SDR scanner: Trunk Recorder + rdio-scanner (on pve2) |

SSH key default: `~/.ssh/id_ansible`. Exceptions noted in last column.

---

## VMs

| Ansible name         | VMID | IP            | Notes                                          |
|----------------------|------|---------------|------------------------------------------------|
| truenas              | 300  | 192.168.86.40 | TrueNAS Scale NAS (on tower1, 16GB RAM); self-signed TLS |
| homeassistant        | 301  | 192.168.86.41 | HAOS VM (on tower1); web UI port 8123         |
| talos-cp-0           | 400  | 192.168.86.101 | Talos K8s control plane (on tower1)           |
| talos-worker-0       | 410  | 192.168.86.111 | Talos K8s worker (on pve2 / thinkcentre2)     |
| talos-worker-1       | 411  | 192.168.86.112 | Talos K8s worker (on pve3 / thinkcentre3)     |
| talos-worker-2       | —    | 192.168.86.113 | Talos K8s worker (on zotac)                   |

---

## Non-Proxmox Devices

| Ansible name       | IP              | User     | Notes                                            |
|--------------------|-----------------|----------|--------------------------------------------------|
| piboard            | 192.168.86.131  | bwoodwar | Pi 3B; Go dashboard + Waveshare display; SSH key: id_ed25519 |
| klipper-ender5pro  | 192.168.86.136  | bwoodwar | Pi 3B; MainsailOS; Klipper for Ender 5 Pro      |
| klipper-ender3     | 192.168.86.138  | bwoodwar | Pi 3B; MainsailOS; Klipper for Ender 3          |
| legion (cachy)     | 192.168.86.173 (eth0) / 192.168.86.152 (wlan0) | bwoodwar | CachyOS laptop; WireGuard client (wg0, 10.10.0.4); NM dispatcher at /etc/NetworkManager/dispatcher.d/99-wireguard-home auto-downs wg0 on home network, ups it when away |

---

## Service Groups

Defined in `ansible/vars/service_groups.yml`. Manage with `make group-start/stop/status GROUP=<name>`.

| Group         | Description                     | VMIDs                    | always_on | hardware_bound | depends_on     | required_by        |
|---------------|---------------------------------|--------------------------|-----------|----------------|----------------|--------------------|
| core          | Ingress + VPN                   | 200 (traefik), 208 (wg) | yes       | —              | —              | —                  |
| storage       | TrueNAS NAS                     | 300 (truenas)            | yes       | —              | —              | [media]            |
| security      | Authentik SSO                   | 207 (authentik)          | —         | —              | —              | [media, apps]      |
| home          | Smart home automation           | 301 (ha), 214 (z2m), 209 (libby) | — | —           | —              | —                  |
| media         | ARR stack, Plex, Jellyfin       | 202, 203, 204            | —         | —              | [core, storage]| —                  |
| observability | Prometheus, Grafana, OpenClaw   | 205 (monitoring), 206 (openclaw) | — | —           | —              | —                  |
| apps          | Lightweight app services        | 201 (recipe), 211 (kanboard), 215 (claude-os) | — | — | — | —           |
| infra         | Email + PXE boot                | 212 (mailserver), 213 (pxe) | —      | —              | —              | —                  |
| sdr           | SDR scanner                     | 210 (sdr)                | —         | —              | —              | —                  |
| special       | USB passthrough services        | 216 (pwnagotchi)         | —         | yes            | —              | —                  |
| k8s           | Talos Kubernetes cluster        | 400 (cp-0), 410 (w-0), 411 (w-1) | — | —          | —              | —                  |

**Safety rules:**
- `always_on=yes`: stop is refused; groups `core` and `storage` cannot be stopped
- `required_by`: stopping `storage` while `media` is running is blocked; stopping `security` while `media` or `apps` is running is blocked
- `hardware_bound=yes`: excluded from bulk group operations; manage manually

---

## Kubernetes Cluster Topology

Cluster name: `talos-proxmox` | Talos v1.12.5 | Kubernetes v1.31.0

| Role          | Hostname             | IP              | Host node |
|---------------|----------------------|-----------------|-----------|
| VIP           | —                    | 192.168.86.100  | —         |
| Control plane | talos-proxmox-cp-0   | 192.168.86.101  | tower1    |
| Worker 0      | talos-proxmox-worker-0 | 192.168.86.111 | pve2     |
| Worker 1      | talos-proxmox-worker-1 | 192.168.86.112 | pve3     |
| Worker 2      | talos-proxmox-worker-2 | 192.168.86.113 | zotac    |

MetalLB pool: `192.168.86.150–199`

Configs: `talos/_out/kubeconfig`, `talos/_out/talosconfig` (gitignored; must exist locally)

API endpoint: `https://192.168.86.100:6443`

---

## Service URL / Port Reference

All public URLs are HTTPS via Traefik at `192.168.86.20`. Internal ports listed for direct access.

### Media

| Service    | URL                                | Internal                    | Auth       |
|------------|------------------------------------|-----------------------------|------------|
| Plex       | https://plex.woodhead.tech         | http://192.168.86.23:32400  | Plex auth  |
| Jellyfin   | https://jellyfin.woodhead.tech     | http://192.168.86.24:8096   | Jellyfin auth |
| Sonarr     | https://sonarr.woodhead.tech       | http://192.168.86.22:8989   | Authentik  |
| Radarr     | https://radarr.woodhead.tech       | http://192.168.86.22:7878   | Authentik  |
| Prowlarr   | https://prowlarr.woodhead.tech     | http://192.168.86.22:9696   | Authentik  |
| Bazarr     | https://bazarr.woodhead.tech       | http://192.168.86.22:6767   | Authentik  |
| SABnzbd    | https://sabnzbd.woodhead.tech      | http://192.168.86.22:8080   | Authentik  |
| Overseerr  | https://requests.woodhead.tech     | http://192.168.86.22:5055   | Authentik  |
| TrueNAS    | https://nas.woodhead.tech          | https://192.168.86.40:443   | Authentik  |

### Monitoring & Observability

| Service      | URL                                    | Internal                    | Auth      |
|--------------|----------------------------------------|-----------------------------|-----------|
| Grafana      | https://grafana.woodhead.tech          | http://192.168.86.25:3000   | Authentik |
| Prometheus   | https://prometheus.woodhead.tech       | http://192.168.86.25:9090   | Authentik |
| Alertmanager | https://alertmanager.woodhead.tech     | http://192.168.86.25:9093   | Authentik |
| OpenClaw     | https://claw.woodhead.tech             | http://192.168.86.26:18789  | Authentik |

### Infrastructure & Auth

| Service      | URL                                | Internal                    | Auth        |
|--------------|------------------------------------|-----------------------------|-------------|
| Authentik    | https://auth.woodhead.tech         | http://192.168.86.28:9000   | —           |
| Traefik dash | https://traefik.woodhead.tech      | —                           | Authentik   |
| Wireguard    | —                                  | 192.168.86.39               | —           |
| PXE server   | —                                  | 192.168.86.35               | —           |
| Mailcow      | https://mail.woodhead.tech         | http://192.168.86.34:8080   | Mailcow auth |

### Smart Home

| Service       | URL                                | Internal                    | Auth        |
|---------------|------------------------------------|-----------------------------|-------------|
| Home Assistant| https://home.woodhead.tech         | http://192.168.86.41:8123   | HA auth     |
| Zigbee2MQTT   | —                                  | 192.168.86.36               | —           |
| Libby Alert   | https://alert.woodhead.tech        | http://192.168.86.27:8080   | None (public)|

### Apps

| Service       | URL                                | Internal                    | Auth        |
|---------------|------------------------------------|-----------------------------|-------------|
| Kanboard      | https://tasks.woodhead.tech        | http://192.168.86.33:8000   | Authentik   |
| Recipe site   | https://recipes.woodhead.tech      | http://192.168.86.21:80     | None        |
| Claude OS     | https://claude-os.woodhead.tech    | http://192.168.86.37:5173   | None        |
| Claude OS API | https://claude-os-api.woodhead.tech| http://192.168.86.37:8051   | None        |

### SDR / Special

| Service       | URL                                | Internal                    | Auth        |
|---------------|------------------------------------|-----------------------------|-------------|
| SDR Scanner   | https://scanner.woodhead.tech      | http://192.168.86.32:3000   | Authentik   |
| Pwnagotchi    | https://pwnagotchi.woodhead.tech   | http://192.168.86.38:8080   | None (use WireGuard) |

### Static Sites (served from monitoring LXC 192.168.86.25)

| Site           | URL                                | Internal port |
|----------------|------------------------------------|---------------|
| Landing page   | https://woodhead.tech              | 8083          |
| Docs           | https://docs.woodhead.tech         | 8081          |
| Resume         | https://resume.woodhead.tech       | 8082          |
| Homelab site   | https://homelab.woodhead.tech      | 8084          |

### 3D Printers (Klipper/Mainsail)

| Printer        | URL                                | Internal                    |
|----------------|------------------------------------|-----------------------------|
| Ender 5 Pro    | https://ender5.woodhead.tech       | http://192.168.86.136:80    |
| Ender 3        | https://ender3.woodhead.tech       | http://192.168.86.138:80    |
