# wireguard

Install WireGuard VPN on a Debian LXC. Generates server and per-client keypairs, templates configs, starts `wg-quick@wg0`, and fetches client `.conf` files to your local machine.

## What it does

- Installs `wireguard-tools`, `iptables`, `qrencode`
- Enables IP forwarding persistently via sysctl
- Generates server keypair (idempotent — skips if keys exist)
- Generates per-client keypair + preshared key for each entry in `wg_client_list`
- Templates `wg0.conf` (server) and one `<name>.conf` per client
- Generates QR codes for each client (scan with the WireGuard mobile app)
- Fetches client configs to `clients/` on the Ansible control machine

## Usage

```bash
ansible-playbook playbook.yml \
  -i inventory.yml \
  --extra-vars "wg_endpoint=vpn.example.com"
```

Add more clients by editing `wg_client_list` in your inventory or extra-vars:

```yaml
wg_client_list:
  - name: laptop
    ip: "10.10.0.2"
  - name: phone
    ip: "10.10.0.3"
  - name: tablet
    ip: "10.10.0.4"
```

Re-run the playbook after adding clients — it only generates keys for clients that don't have key files yet.

## Variables

| Variable | Default | Description |
|---|---|---|
| `wg_endpoint` | **required** | Public hostname or IP for the VPN endpoint |
| `wg_port` | `51820` | UDP listen port |
| `wg_subnet` | `10.10.0.0/24` | VPN tunnel subnet |
| `wg_server_ip` | `10.10.0.1` | Server's VPN address |
| `wg_dns` | `8.8.8.8` | DNS server pushed to clients |
| `wg_allowed_ips` | LAN + VPN subnet | Routes pushed to clients |
| `wg_client_list` | laptop + phone | List of `{name, ip}` client peers |

## Files

```
wireguard/
├── playbook.yml      # Ansible playbook
├── wg0.conf.j2       # Server config template
├── client.conf.j2    # Client config template
├── clients/          # Fetched client .conf files (gitignored)
└── README.md
```

## After deploy

1. **Port forward** UDP `wg_port` on your router to the WireGuard LXC IP
2. **Import** `clients/<name>.conf` into the WireGuard app on each device
3. **Mobile QR**: `ssh root@<wg-ip> cat /etc/wireguard/clients/phone.qr`

## Notes

- Keys are generated on the remote host and never leave it except via the fetch task. Client configs contain private keys — keep them out of git.
- The `clients/` directory is excluded from git via `.gitignore`.
- LXC containers need `net_admin` capability and `/dev/net/tun` device for WireGuard. In Proxmox, the LXC features block should include `nesting=1` and you may need to set `lxc.cap.keep: net_admin` in the container's config.
