# Contributing

Contributions welcome. The most valuable areas:

- **Recipes** — new service templates (Terraform modules, Ansible playbooks, Helm charts)
- **MCP servers** — new tools or improvements to existing ones
- **CLAUDE.md** — better builder prompts, new architecture patterns, corrections
- **Website** — docs, guides, screenshots

## Adding a recipe

Recipes live in `recipes/`. Each is self-contained:

```
recipes/ansible/my-service/
  main.yml          # The playbook
  vars.yml.example  # Variables the user needs to fill in (no defaults for IPs/passwords)
  README.md         # What this deploys, variables reference, any manual steps
```

Requirements:
- No hardcoded IPs, hostnames, or credentials — everything via variables
- A `README.md` that explains what the recipe does and what vars are required
- Tested against at least one real Proxmox setup

## Adding an MCP tool

Each MCP server is in `mcp/<name>/`. Tools should:
- Do one thing
- Return clean, structured output (JSON where possible)
- Return useful error messages — Claude reads these to decide what to do next
- Never execute destructive operations without a dry-run tool variant

## Pull requests

- Keep PRs focused — one recipe or one feature per PR
- Include a short description of what you tested it against
- For recipes: include the `README.md`

## Code style

Go: `gofmt`. No exceptions.
YAML: 2-space indent.
