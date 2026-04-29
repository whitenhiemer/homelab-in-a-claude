# arr-stack

Deploy the ARR media management stack via Docker Compose on a Debian LXC.

## Services

| Service | Port | Purpose |
|---|---|---|
| Prowlarr | 9696 | Indexer manager (Usenet + torrents) |
| Sonarr | 8989 | TV show library manager |
| Radarr | 7878 | Movie library manager |
| Bazarr | 6767 | Subtitle downloader |
| Overseerr | 5055 | User request portal |
| SABnzbd | 8080 | Usenet downloader (behind gluetun VPN) |
| gluetun | — | WireGuard VPN killswitch for SABnzbd |

SABnzbd routes all download traffic through gluetun's WireGuard tunnel. If the VPN drops, SABnzbd loses network access entirely (killswitch behavior).

## Usage

```bash
ansible-playbook playbook.yml \
  -i inventory.yml \
  --extra-vars "
    wg_private_key=<your_wireguard_private_key>
    wg_public_key=<server_public_key>
    wg_address=100.64.x.x/32
    wg_endpoint_ip=<vpn_server_ip>
    nfs_server=192.168.86.40
    nfs_share=/mnt/pool/media
  "
```

## WireGuard keys

Get these from your VPN provider's WireGuard config file. The private key must stay secret — never commit it.

## NFS media mount

If you have TrueNAS or another NFS server, set `nfs_server` and `nfs_share`. Media will mount at `/media` inside the container. The LXC must have the `nfs-common` package and the Proxmox host must bind-mount the NFS path into the container.

## Files

```
arr-stack/
├── playbook.yml              # Ansible playbook
├── docker-compose.yml.j2     # Stack template (customize for your VPN provider)
└── README.md
```

## Post-deploy order

1. Configure Prowlarr indexers first
2. Connect Sonarr and Radarr to Prowlarr
3. Connect Sonarr, Radarr to SABnzbd
4. Set up Overseerr pointing to Sonarr + Radarr
5. Set hardlink paths — SABnzbd downloads to `/media/downloads`, Sonarr/Radarr import from same path
