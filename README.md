# SysUpdate – Linux System Update Utility

![Platform](https://img.shields.io/badge/platform-Linux-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![AppImage](https://img.shields.io/badge/Distribution-AppImage-orange)
![GitHub stars](https://img.shields.io/github/stars/RossContino1/SysUpdate)

⭐ If you find SysUpdate useful, consider starring the project!

---

## Overview

SysUpdate is a simple graphical Linux system update utility.

Instead of remembering package manager commands, SysUpdate provides a clean graphical interface that detects your Linux distribution and runs the appropriate update commands automatically.

The program uses native package managers and requests administrator privileges securely using **pkexec**.

SysUpdate is distributed as a **portable AppImage**, allowing it to run on most Linux systems without requiring a traditional installation.

---

## Homepage

https://bytesbreadbbq.com/sysupdate/

---

## Source Code

https://github.com/RossContino1/SysUpdate

---

## Features

• Simple graphical interface for system updates  
• Automatically detects Linux distribution  
• Uses native package managers  
• Secure privilege escalation using **pkexec**  
• Portable **AppImage** distribution  
• Optional install and uninstall scripts  
• No system-wide installation required  

---

## Supported Distributions

SysUpdate currently supports:

- Ubuntu
- Debian
- Linux Mint
- Fedora
- Arch Linux
- openSUSE Tumbleweed

Support for additional distributions may be added in future releases.

---

## Requirements

SysUpdate requires the following components:

### pkexec (PolicyKit)

Used to securely request administrator privileges.

Install examples:

**Fedora**
sudo dnf install polkit


**Debian / Ubuntu**
sudo apt install policykit-1


---

### FUSE (required for AppImage)

Most Linux distributions already include FUSE.

Example installations:

**Fedora**
sudo dnf install fuse

**Debian / Ubuntu**
sudo apt install fuse


---

## Installation

SysUpdate installs to the user's home directory and does not require administrator privileges.

Installed locations:
~/.local/bin/SysUpdate.AppImage
~/.local/share/applications/com.bytesbreadbbq.sysupdate.desktop
~/.local/share/icons/hicolor/256x256/apps/sysupdate.png


### Install Steps

Open a terminal in the folder containing the files.

Make the installer executable:
chmod +x install.sh

Run the installer:
./install.sh


After installation, SysUpdate will appear in your desktop environment's application menu.

---

## Running SysUpdate

Launch SysUpdate from your desktop application menu.

When an operation requires administrator privileges, the system will prompt for authentication using **pkexec**.

---

## Uninstall

To remove SysUpdate:
chmod +x uninstall.sh
./uninstall.sh


This removes the AppImage, launcher, and icon.

---

## Legacy Installations

Older versions of SysUpdate installed files into:
/usr/local


If a legacy installation is detected, the install or uninstall script may offer to remove it. Administrator permission is required.

---

## Support / Issues

For bug reports or feature requests, please open an issue on GitHub:

https://github.com/RossContino1/SysUpdate/issues

---

## More Projects from Bytes, Bread, and Barbecue

Check out other projects from this developer:

**Leonardo – Linux Media Conversion Application**

https://github.com/RossContino1/Leonardo
https://bytesbreadbbq.com/Leonardo

A fast and simple graphical front-end for FFmpeg designed for modern Linux systems.

---

## ☕ Support SysUpdate

SysUpdate is free to use. If it saves you time (or brisket), consider supporting development:

[![Support via PayPal](https://img.shields.io/badge/Support-PayPal-blue?style=for-the-badge&logo=paypal)](https://www.paypal.com/donate/?hosted_button_id=XS9MXN5AE5P3S)

Your support helps.

## License

This project is licensed under the **MIT License**.

---

## Bytes, Bread, and Barbecue

At **Bytes, Bread, and Barbecue**, we like our code crispy and our software smokin’ hot.
