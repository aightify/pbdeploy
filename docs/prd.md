Here is the updated **Product Requirements Document (PRD)** reflecting your latest directive: deployment is allowed with `sudo`, and the installation must be **super minimal**, using only **systemd** and no third-party tools.

---

# 🧾 Product Requirements Document (PRD)

**Project Name:** PocketBase Deploy CLI & Agent
**Prepared For:** Engineering Team
**Prepared By:** Senthuraan Ponnampalam
**Date:** 2025-06-23
**Status:** Final – Approved for Development

---

## 🧭 1. Purpose

To create a **Go CLI tool** (`pbdeploy`) and a minimal **remote deployment agent** (`pbdeploy-agent`) that allows developers to deploy PocketBase apps to Ubuntu-based servers with one command. The installation is **sudo-based**, **systemd-backed**, and intentionally minimal with **zero third-party dependencies**.

---

## 🎯 2. Goals & Objectives

* One-command deployment for PocketBase apps via CLI.
* Remote agent installed and managed with `systemd`.
* Minimal install footprint: just one binary and one service file.
* Secure, persistent, and restart-on-failure design.
* Optional GitHub webhook integration for auto-updates.

---

## 📦 3. Features

### ✅ 3.1 CLI Tool (`pbdeploy`)

| Feature     | Description                                             |
| ----------- | ------------------------------------------------------- |
| `init`      | Creates `pbdeploy.yml` with deployment descriptor       |
| `deploy`    | Deploys the app by sending a command to the agent       |
| `install`   | Installs `pbdeploy-agent` remotely via SSH using `sudo` |
| `status`    | Checks if agent is running                              |
| `webhook`   | Installs GitHub webhook on repo                         |
| `uninstall` | Removes the agent and disables the service              |

---

### ✅ 3.2 Deployment Agent (`pbdeploy-agent`)

| Feature                | Description                          |
| ---------------------- | ------------------------------------ |
| Minimal Go binary      | \~3MB binary, no dependencies        |
| systemd service        | Starts on boot, restarts on crash    |
| GitHub webhook handler | Accepts push events to auto-deploy   |
| `/deploy` endpoint     | Pulls repo, builds, and restarts app |
| `/healthz` endpoint    | Shows agent status, uptime, version  |

---

### ✅ 3.3 `pbdeploy.yml` Configuration File

```yaml
server: ubuntu@1.2.3.4
ssh_key: ~/.ssh/id_rsa
repo: git@github.com/user/myapp.git
branch: main
app_name: myapp
env:
  PORT: 8080
  DATABASE_URL: sqlite://...
post_deploy: |
  systemctl restart myapp
webhook:
  enable: true
  secret: abcdef123456
```

---

## ⚙️ 4. System Architecture

```
Developer Laptop        Ubuntu Server
+---------------+       +------------------------+
| pbdeploy CLI  |  ==>  | pbdeploy-agent         |
|               |       | - /usr/local/bin/      |
|               |       | - systemd service      |
+---------------+       +------------------------+
     |                           |
     |-- webhook setup --> GitHub
```

---

## 🛠️ 5. Installation Details (with sudo)

### 🔩 Agent Installation Flow (`pbdeploy install`)

1. Copy binary to `/usr/local/bin/pbdeploy-agent`
2. Create systemd service at `/etc/systemd/system/pbdeploy-agent.service`
3. Enable and start the service:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now pbdeploy-agent
   ```

### 🧾 Sample systemd service file:

```ini
[Unit]
Description=PocketBase Deploy Agent
After=network.target

[Service]
ExecStart=/usr/local/bin/pbdeploy-agent
Restart=on-failure
User=ubuntu
WorkingDirectory=/home/ubuntu

[Install]
WantedBy=multi-user.target
```

---

## 🧰 6. Tech Stack

| Component          | Tech                   |
| ------------------ | ---------------------- |
| CLI Tool           | Go + Cobra             |
| Agent              | Go                     |
| Remote Deployment  | SSH + `scp`            |
| Service Management | `systemd`              |
| GitHub Webhooks    | HMAC-secured, JSON API |

---

## 🔐 7. Security

* All communication is via SSH (with key authentication)
* GitHub webhook payloads are validated with HMAC secret
* Agent only exposes required endpoints
* No external ports need to be opened manually (optional reverse proxy allowed)

---

## 🧪 8. Testing & Validation

| Test              | Description                                     |
| ----------------- | ----------------------------------------------- |
| CLI command tests | Verify all CLI subcommands behave as expected   |
| Agent integration | Ensure agent builds and installs cleanly        |
| Remote deployment | Ensure full flow from CLI to server works       |
| Webhook trigger   | Push to GitHub should auto-deploy               |
| Systemd recovery  | Kill process; verify `Restart=on-failure` works |
| Reboot recovery   | Ensure service comes up on reboot               |

---

## 📆 9. Development Timeline

| Phase                | Time   | Responsibility |
| -------------------- | ------ | -------------- |
| CLI Base (`cobra`)   | 2 days | CLI Dev        |
| Agent Binary         | 2 days | Backend Dev    |
| SSH Install Logic    | 2 days | CLI Dev        |
| GitHub Webhook       | 2 days | Fullstack Dev  |
| systemd integration  | 1 day  | DevOps         |
| Tests & Finalization | 2 days | QA             |
| Buffer               | 2 days | –              |

**Total:** 11 working days (\~2 weeks)

---

## 📘 10. Deliverables

* `pbdeploy` CLI binary
* `pbdeploy-agent` binary
* `pbdeploy.yml` schema and validator
* systemd unit file template
* Documentation: README + Quickstart + Webhook Setup Guide

---

## ✨ 11. Out of Scope / Future Ideas

| Feature             | Reason                    |
| ------------------- | ------------------------- |
| Docker-based deploy | Adds complexity           |
| SFTP backup         | Optional, future          |
| Rollbacks           | To be considered in v2    |
| Multi-environment   | Staging/prod coming later |

---

## ✅ Summary

The updated PRD now supports `sudo`-based, **minimalist installation via systemd**. It is production-ready, simple, and developer-friendly. No additional packages or dependencies required on the server side except standard Ubuntu + systemd.

---

Let me know if you’d like:

* A markdown export
* A `.pdf` version
* A code scaffold for CLI + agent setup
