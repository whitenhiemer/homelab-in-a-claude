# proxmox-lxc

Provision a single Debian 12 LXC container on Proxmox using the `bpg/proxmox` provider.

## What it creates

- Unprivileged Debian 12 LXC container
- Static IP, DNS, and SSH key injection via cloud-init
- Optional nesting for Docker workloads

## Usage

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your Proxmox endpoint, API token, IP, etc.

terraform init
terraform plan
terraform apply
```

## Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `proxmox_endpoint` | yes | — | Proxmox API URL |
| `proxmox_api_token` | yes | — | `USER@REALM!TOKENID=SECRET` |
| `proxmox_node` | yes | — | Node to place the container on |
| `vm_id` | yes | — | Unique container ID |
| `hostname` | yes | — | Container hostname |
| `ip_address` | yes | — | Static IPv4 (no prefix) |
| `gateway` | yes | — | Default gateway |
| `cpu_cores` | no | `1` | CPU cores |
| `memory_mb` | no | `512` | RAM in MB |
| `disk_size_gb` | no | `8` | Root disk size |
| `storage` | no | `local-lvm` | Proxmox storage pool |
| `enable_nesting` | no | `false` | Enable Docker inside LXC |

## Prerequisites

1. Proxmox API token with `PVEVMAdmin` role (or `Administrator` for simplicity)
2. Debian 12 LXC template downloaded on the node:
   ```bash
   pveam update
   pveam download local debian-12-standard_12.7-1_amd64.tar.zst
   ```
3. `terraform >= 1.5` installed locally
