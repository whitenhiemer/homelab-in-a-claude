# truenas

Configure TrueNAS Scale via its REST API after the OS has been installed. Creates a ZFS pool, organizes datasets, enables NFS, and exports shares to the LAN.

Runs entirely from the Ansible control machine over HTTP — no SSH to TrueNAS is needed.

## What it configures

- ZFS pool `tank` (single-disk stripe by default, extend with mirrors/RAIDz in the UI)
- Datasets with LZ4 compression:
  - `tank/media` + subdatasets: `movies`, `tv`, `music`, `downloads`
  - `tank/backups` — Proxmox VM backups and Home Assistant backups
  - `tank/isos` — Proxmox ISO images
- NFS service enabled and started
- NFS exports for `media`, `backups`, `isos` restricted to your LAN CIDR

## Usage

```bash
ansible-playbook playbook.yml \
  --extra-vars "truenas_password=your-admin-password"
```

The playbook is idempotent — re-running it skips resources that already exist.

## Variables

| Variable | Default | Description |
|---|---|---|
| `truenas_host` | `http://192.168.86.40` | TrueNAS URL |
| `truenas_password` | **required** | Admin password |
| `pool_name` | `tank` | ZFS pool name |
| `data_disk` | `sdb` | Disk for the ZFS pool (not the OS disk) |
| `nfs_network` | `192.168.86.0/24` | CIDR allowed to mount NFS shares |

## Prerequisites

1. TrueNAS Scale installed on the OS disk via Proxmox console
2. VM booted and reachable at `truenas_host`
3. At least one data disk attached to the VM (not `/dev/sda` which is the OS disk)

To find the data disk name inside TrueNAS:
```bash
# In TrueNAS Shell (Storage > Shell) or via SSH:
lsblk -d -o NAME,SIZE,MODEL
```

## After this playbook

1. Add NFS storage to Proxmox: **Datacenter → Storage → Add → NFS**
   - ID: `truenas-isos`, Server: `<truenas-ip>`, Share: `/mnt/tank/isos`
   - ID: `truenas-backups`, Server: `<truenas-ip>`, Share: `/mnt/tank/backups`
2. Bind-mount `/mnt/tank/media` into ARR, Plex, and Jellyfin LXCs via Proxmox:
   ```bash
   # On Proxmox host — adds NFS path as LXC mount point
   pct set <vmid> -mp0 /mnt/truenas-media,mp=/media
   ```
3. Configure Proxmox backup job targeting the `truenas-backups` storage

## Disk passthrough note

For a real NAS with dedicated data disks, use disk passthrough rather than a virtual disk. After TrueNAS VM is created:

```bash
# On Proxmox host — list physical disks
lsblk -d -o NAME,SIZE,MODEL,SERIAL

# Pass through a disk by stable ID (survives reboots)
qm set <vmid> -scsi1 /dev/disk/by-id/<disk-id>
```
