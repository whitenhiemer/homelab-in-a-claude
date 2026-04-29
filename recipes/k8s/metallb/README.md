# metallb

MetalLB IP address pool for bare-metal Kubernetes LoadBalancer services.

MetalLB gives `LoadBalancer`-type Services real IPs on your LAN, using ARP (L2 mode). This means any service you expose via `type: LoadBalancer` gets a stable IP that clients on your network can reach directly — no cloud load balancer needed.

## Usage

```bash
# Install MetalLB (requires cluster with CRDs disabled in kube-proxy, or Talos)
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.14.8/config/manifests/metallb-native.yaml

# Wait for MetalLB pods to be ready
kubectl -n metallb-system wait --for=condition=ready pod --selector=app=metallb --timeout=90s

# Apply the IP pool (edit ip-pool.yml first to match your subnet)
kubectl apply -f ip-pool.yml
```

## Edit before applying

Change the address range in `ip-pool.yml` to an unused block in your subnet:
```yaml
addresses:
  - 192.168.1.200-192.168.1.250  # adjust to your network
```

Ensure this range:
- Is in the same /24 as your Kubernetes nodes
- Does NOT overlap with your router's DHCP range
- Has no static IPs assigned in that range

## Talos note

On Talos, MetalLB requires `kube-proxy` to be disabled and `strictARP` enabled. Add to your Talos machine config:

```yaml
cluster:
  proxy:
    disabled: true
```

And enable strict ARP via the MetalLB `L2Advertisement` (already in `ip-pool.yml`).
