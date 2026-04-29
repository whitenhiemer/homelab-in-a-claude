# proxmox-vm

Provision a generic VM on Proxmox using the `bpg/proxmox` provider.

Boots from an ISO for initial OS installation. After the OS is installed, Terraform ignores further ISO changes (so re-running `apply` won't unmount your CD-ROM).

## What it creates

- KVM VM with virtio networking and SCSI disk
- Serial console device for headless access via Proxmox UI
- QEMU guest agent support (install `qemu-guest-agent` in the OS)
- Optional RAM ballooning (set `memory_balloon_mb > 0`)

## Usage

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit with your Proxmox details, ISO path, and sizing

terraform init
terraform plan
terraform apply
```

After Terraform creates the VM, open the Proxmox console (or use the serial device) to complete the OS installation.

## Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `proxmox_endpoint` | yes | — | Proxmox API URL |
| `proxmox_api_token` | yes | — | `USER@REALM!TOKENID=SECRET` |
| `proxmox_node` | yes | — | Node to create the VM on |
| `vm_id` | yes | — | Unique VM ID |
| `name` | yes | — | VM display name |
| `iso` | yes | — | ISO file ID for installation |
| `cpu_cores` | no | `2` | CPU cores |
| `memory_mb` | no | `2048` | RAM ceiling in MB |
| `memory_balloon_mb` | no | `0` | RAM floor (0 disables ballooning) |
| `disk_size_gb` | no | `32` | OS disk size |
| `storage` | no | `local-lvm` | Storage pool |

## Notes

- Static IPs must be configured inside the guest OS — this recipe has no cloud-init support. For cloud-init VMs, clone from a cloud image instead of installing from ISO.
- Disk passthrough (e.g., for TrueNAS ZFS pools) must be added manually after VM creation: `qm set <vmid> -scsi1 /dev/disk/by-id/<id>`
