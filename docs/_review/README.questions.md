# README.md Questions & Verification Log

In accordance with [`docs/_review/CONTRACT.md`](CONTRACT.md) § 4.3 ("Never Invent Facts"), any missing parameters, unconfirmed behaviors, or specifications not explicitly established in CONTRACT.md are documented below with their corresponding `TODO(verify)` tags.

| Issue ID | Tag in Documentation | Question / Required Verification | Impacted Section |
| :--- | :--- | :--- | :--- |
| **Q-README-01** | `TODO(verify): default_config_search_order` | What is the exact search hierarchy for configuration files when `cli-leader start` is executed without the `--config` flag? (e.g., `./cli-leader.yaml` $\to$ `~/.config/cli-leader/config.yaml` $\to$ `/etc/cli-leader/config.yaml`). | Quickstart & Configuration |
| **Q-README-02** | `TODO(verify): claude_cli_executable_name` | What is the canonical executable name and discovery mechanism for the Claude CLI brain process? Does `cli-leader` look for `claude`, `claude-code`, or an explicit path defined under `cli-leader.yaml`? | Quickstart & Overview |
| **Q-README-03** | `TODO(verify): env_var_overrides` | What is the precedence order and naming convention for environment variable overrides relative to `cli-leader.yaml` configuration keys (e.g., does `CLI_LEADER_GATEWAY_LISTEN_ADDR` override `gateway.listen_addr`)? | Configuration Snippet |
