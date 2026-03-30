# RepoRover

🚀 **Update ALL your Linux packages in ONE command — including AUR**

No more juggling:
- apt
- snap
- flatpak
- zypper
- pacman
- yay / paru

👉 RepoRover detects everything and updates it for you.

---

⭐ **If this saves you time, please star the repo — it helps more than you think**

⬇️ **Download RepoRover**  
https://bytesbreadbbq.com/reporover

📰 **Featured on LinuxLinks (independent Linux review site)**  
https://www.linuxlinks.com/reporover-universal-linux-package-updater/

🎥 Seen on TikTok / YouTube

---

<p align="center">
  <img src="docs/output.gif" height="400" style="border:3px solid black" alt="RepoRover demo"/><br>
  <sub>One command. Every package manager. Zero hassle.</sub>
</p>

---

## ⚡ Quick Start (30 seconds)

1. Download RepoRover  
2. Run the AppImage  
3. Click “Update”  

Done.

No setup. No configuration.

---

## ⭐ Why people star RepoRover

- Saves time every single update
- Works across multiple package managers automatically
- Makes Linux easier without removing power

👉 If that sounds useful, give it a ⭐ — it helps others discover the project.

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
- 🚀 Improved AUR handling in latest versions

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

👉 https://bytesbreadbbq.com/reporover

---

## 🔐 Permissions (Important)

RepoRover uses **pkexec** (part of PolicyKit) to securely request administrator privileges when running system updates.

### Distribution Notes

- **Ubuntu / Debian / Linux Mint**  
  Works out of the box.

- **Arch / CachyOS**  
  Works out of the box (requires PolicyKit / polkit).

- **openSUSE**  
  Works out of the box.

- **Fedora**  
  If needed:

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

### ✅ Smart AUR Handling

RepoRover:
- Detects which helper actually works
- Uses the same helper for detection and updates
- Prevents failures when both are installed

👉 No more guessing or broken AUR updates.

---

## 📈 Growing Fast

RepoRover is gaining traction across:
- TikTok
- YouTube
- Linux communities

👉 Join early and be part of it.

---

## 🧰 More Tools

Check out my other projects:

- 🎬 **Domenico** — Turn videos into looping GIFs instantly  
- 🎞️ **Leonardo** — Convert videos for DaVinci Resolve on Linux  

---

## ⭐ Support the Project

If RepoRover saves you time:

- ⭐ Star the repo  
- 🔁 Share it with other Linux users  
- ☕ Support development  

[![Support via PayPal](https://img.shields.io/badge/Support-PayPal-blue?style=for-the-badge&logo=paypal)](https://www.paypal.com/donate/?hosted_button_id=XS9MXN5AE5P3S)

---

## 📄 License

MIT License
