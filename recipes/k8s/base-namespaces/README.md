# base-namespaces

Base Kubernetes namespace definitions. Apply immediately after cluster bootstrap.

## Namespaces

| Namespace | Purpose |
|---|---|
| `ingress-system` | Ingress controllers, load balancer infra |
| `apps` | Application workloads |
| `monitoring` | Prometheus, Grafana, node-exporter (privileged pod security) |

## Usage

```bash
kubectl apply -f namespaces.yml
kubectl get namespaces
```

## Notes

The `monitoring` namespace sets `pod-security.kubernetes.io/enforce: privileged` because node-exporter requires `hostNetwork`, `hostPID`, and `hostPath` mounts to scrape kernel metrics. Without this label, Talos's default pod security admission will block the daemonset.
