# talos-cluster

Provision Talos Linux Kubernetes cluster VMs on Proxmox. Creates control plane and worker node VMs booting from the Talos ISO. Talos cluster initialization is done separately with `talosctl` after VMs exist.

## What it creates

- N control plane VMs booting Talos ISO, spread across Proxmox nodes
- M worker VMs booting Talos ISO, spread across Proxmox nodes
- All VMs on virtio networking, raw SCSI disk, serial console

## Usage

```bash
# 1. Download Talos ISO to Proxmox
# On Proxmox host:
wget -O /var/lib/vz/template/iso/talos-amd64.iso \
  https://github.com/siderolabs/talos/releases/download/v1.9.0/metal-amd64.iso

# 2. Apply Terraform
cp terraform.tfvars.example terraform.tfvars
# Edit with your Proxmox details and desired IPs
terraform init
terraform plan
terraform apply

# 3. Bootstrap the cluster (run from your workstation)
export CP_IP="192.168.86.101"
export VIP="192.168.86.100"
export WORKERS="192.168.86.111 192.168.86.112"

talosctl gen config talos https://${VIP}:6443 \
  --config-patch '[{"op":"add","path":"/machine/network/interfaces/0/vip","value":{"ip":"'${VIP}'"}}]'

talosctl apply-config --insecure --nodes ${CP_IP} --file controlplane.yaml

for w in ${WORKERS}; do
  talosctl apply-config --insecure --nodes ${w} --file worker.yaml
done

talosctl bootstrap --nodes ${CP_IP} --endpoints ${CP_IP} --talosconfig talosconfig
talosctl kubeconfig --nodes ${CP_IP} --endpoints ${CP_IP} --talosconfig talosconfig
kubectl get nodes
```

## Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `proxmox_endpoint` | yes | — | Proxmox API URL |
| `proxmox_api_token` | yes | — | `USER@REALM!TOKENID=SECRET` |
| `proxmox_node` | yes | — | Default node for VMs |
| `cluster_name` | no | `talos` | Name prefix for VMs |
| `talos_iso` | yes | — | ISO file ID on Proxmox |
| `cluster_vip` | yes | — | Kubernetes API VIP |
| `controlplane_count` | no | `1` | Control plane replicas |
| `controlplane_ips` | yes | — | Control plane IPs (informational) |
| `worker_count` | no | `2` | Worker replicas |
| `worker_ips` | yes | — | Worker IPs (informational) |

## Notes

- IPs are informational only — Terraform doesn't configure Talos networking. Talos reads IPs from the DHCP lease initially, then locks them via the machine config you apply.
- For HA control plane (3 nodes), set `controlplane_count = 3` and spread across 3 `controlplane_nodes`.
- The VIP requires all control plane nodes to be on the same L2 subnet.
- Use `talhelper` or `talosctl gen config` to generate machine configs. The `talconfig.yaml` approach scales better for multi-node clusters.
