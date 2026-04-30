---
title: "Why I Built a Homelab with an AI"
date: 2026-03-10
---

I've run a homelab for years. At various points it has included a Proxmox cluster, a NAS, a Kubernetes cluster, a full media server stack, WireGuard, Home Assistant, and a monitoring setup that pages me when my blood glucose goes out of range. I enjoy building it. I also forget how I built it.

The documentation problem is real. You set something up, it works, you move on. Six months later you want to add a second node or migrate a service and you're staring at config files with no memory of why they're structured the way they are. You check your notes and there are no notes. You Google around and piece something together that's subtly different from the working version.

I wanted to solve this without actually writing documentation — because I will never consistently write documentation. The idea was: what if the system itself knew how it was built?

---

## What I tried first

My first instinct was IaC discipline. Terraform for everything, Ansible playbooks with comments, commit messages that explained the why. This works well enough while you're in the flow. It doesn't hold up when you've been away from a project for months, or when you're handing it to someone else (or your future self after a hard drive failure).

The deeper problem is that infrastructure knowledge isn't just in the files — it's in the sequence of decisions, the things you tried that didn't work, the constraints that shaped the design. A Terraform module tells you *what* you created but not *why the IP was chosen that way* or *why you used a VM instead of a container for that service*.

---

## The Claude Code angle

Claude Code shipped with MCP (Model Context Protocol) support — a way to give Claude persistent tools it can use across a session. You can point it at your infrastructure directly: Proxmox API, Cloudflare DNS, SSH sessions, kubectl, Terraform. Instead of copying and pasting output into a chat window, Claude can run commands and read results itself.

The implication is that Claude can hold the full context of your homelab while you're working on it. Not just the files, but the live state — what's running, what the logs say, whether the cert is expiring, whether the container started.

That's a different kind of documentation. It's not a README you write once and forget to update. It's a system that can answer questions about itself, propose changes with awareness of what already exists, and walk you through operations it's done before.

---

## What this project is

The repo is three things bundled together:

**MCP servers** that give Claude access to your infrastructure. There's one for Proxmox, one for Cloudflare, one for SSH, one for Terraform, one for Ansible, one for kubectl. Each one translates Claude's intent into API calls or shell commands.

**Recipes** — parameterized Terraform modules and Ansible playbooks that Claude can select and apply. The traefik recipe sets up the reverse proxy. The arr-stack recipe deploys the media management stack behind a VPN. The monitoring recipe wires up Prometheus and Grafana. All of them are standalone — you can use them without Claude if you want.

**A CLAUDE.md** that tells Claude how the project is structured, what conventions to follow, and how to pick the right recipe for a given task. This is where the institutional knowledge lives in a form Claude can actually read.

The result is a homelab that can onboard you to itself. Point Claude at it and ask "what services are running?" and it checks. Ask "how do I add a new LXC for Jellyfin?" and it finds the right recipe, fills in your network conventions, and walks you through it.

---

## What it doesn't do

It doesn't make decisions for you. Claude proposes things and waits for confirmation before touching anything. If you're adding a VM and you haven't thought through the IP, Claude will ask — it won't just pick one.

It also doesn't replace knowing your infrastructure. The value is in the acceleration and the context retention, not in being able to operate something you don't understand. If anything, walking through an operation with Claude tends to surface more of the reasoning than doing it alone, which is good for actually learning how things work.

---

The next post walks through a concrete example: setting up a new homelab from scratch with this project, from bare metal to a working Traefik + Authentik + media server stack.
