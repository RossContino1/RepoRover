# RepoRover
**Update apt, snap, flatpak, pacman, zypper, and AUR — all in one click.**

<p align="center">
  <img src="docs/output.gif" height="400" style="border:3px solid black" alt="RepoRover demo"/><br>
  <sub>One tool. Every package manager. Zero command memorization.</sub>
</p>

---

⭐ **If RepoRover saves you time, consider starring the project — it really helps visibility.**

---

## 🚀 What is RepoRover?

**RepoRover** is a graphical Linux update utility that eliminates the need to remember package manager commands.

It automatically:
- Detects your Linux distribution
- Detects installed package managers
- Runs the correct update workflow for your system

👉 No more:
- remembering `apt`, `dnf`, `pacman`, etc.
- running multiple update commands
- guessing which package manager to use

---

## ⚡ Why use RepoRover?

Because Linux updates shouldn’t look like this:

```bash
sudo apt update && sudo apt upgrade
flatpak update
snap refresh
```

Or this:

```bash
sudo pacman -Syu
yay -Syu
```

👉 RepoRover does it all for you — in one click.

---

## ✨ Features

- 🧠 Auto-detects Linux distribution
- 🔍 Auto-detects installed package managers
- ⚙️ Supports:
  - `apt`
  - `snap`
  - `flatpak`
  - `zypper`
  - `pacman`
  - AUR helpers: `yay`, `paru`
- 🔐 Secure privilege escalation using **pkexec**
- 📦 Portable **AppImage** (no install required)
- 🧰 Optional install/uninstall scripts for menu integration
- 🚀 Improved AUR handling in **v1.2.0**

---

## 🐧 Supported Distributions

- Ubuntu
- Debian
- Linux Mint
- Fedora
- Arch Linux
- openSUSE Tumbleweed
- CachyOS

---

## 📸 Screenshots

<p align="center">
  <b>🧭 Main Dashboard</b><br>
  <img src="docs/1.png" width="700" style="border:3px solid black" alt="RepoRover main dashboard"/>
</p>

<p align="center">
  <b>🔍 Distro Detection</b><br>
  <img src="docs/2.png" width="700" style="border:3px solid black" alt="RepoRover distro detection"/>
</p>

<p align="center">
  <b>🔐 Privilege Prompt</b><br>
  <img src="docs/3.png" width="700" style="border:3px solid black" alt="RepoRover privilege prompt"/>
</p>

<p align="center">
  <b>📊 Update Results</b><br>
  <img src="docs/4.png" width="700" style="border:3px solid black" alt="RepoRover results"/>
</p>

---

## 📦 Download

[bytesbreadbbq.com/reporover](https://bytesbreadbbq.com/reporover/)

---

## 🔐 Permissions (Important)

RepoRover uses **pkexec** (part of PolicyKit) to securely request administrator privileges when running system updates.

### Distribution Notes

- **Ubuntu / Debian / Linux Mint**  
  Works out of the box.

- **Arch / CachyOS**  
  Works out of the box. GUI privilege prompts rely on PolicyKit / polkit being present, which is typical on these systems.

- **openSUSE**  
  Works out of the box.

- **Fedora**  
  Usually includes PolicyKit already, but if it is missing:

```bash
sudo dnf install polkit
```

### If no password prompt appears

Your system may be missing or misconfigured PolicyKit.

RepoRover does **not** bypass system security. It uses your system’s native authentication flow through `pkexec`.

---

## ⚠️ AUR (Arch-based systems)

RepoRover supports both:
- `yay`
- `paru`

### How AUR helper detection works

- RepoRover checks whether `yay` and/or `paru` are installed.
- If only one is present, RepoRover uses that helper.
- If both are present, RepoRover currently selects an available helper automatically.

### Important note about `yay` vs `paru`

On some systems, both helpers may be installed, but one may work better than the other for a particular setup.

For example, a package may have originally been installed using `yay`, while `paru` is also present on the machine. In rare cases, `paru` may fail even though `yay` succeeds when run manually.

If that happens:
- Run updates manually using your preferred helper:
  - `yay -Syu`
  - `paru -Syu`
- Then re-run RepoRover if needed.

Future versions of RepoRover will continue improving AUR helper selection logic.

---

## 🌐 Homepage

[bytesbreadbbq.com/reporover](https://bytesbreadbbq.com/reporover)

---

## 🌐 Source Code

[github.com/RossContino1/RepoRover](https://github.com/RossContino1/RepoRover)

---

## ⭐ Support the Project

If RepoRover saves you time, consider:

- ⭐ Starring the repo
- 🔁 Sharing it with other Linux users
- ☕ Supporting development

[![Support via PayPal](https://img.shields.io/badge/Support-PayPal-blue?style=for-the-badge&logo=paypal)](https://www.paypal.com/donate/?hosted_button_id=XS9MXN5AE5P3S)

---

## 📄 License

This project is licensed under the MIT License.
