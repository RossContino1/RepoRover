SysUpdate
A simple graphical Linux system update utility.

Author: Ross Contino

Homepage:
https://bytesbreadbbq.com/sysupdate/

Source Code (GitHub):
https://github.com/RossContino1/SysUpdate

---

## Overview

SysUpdate provides a simple graphical interface for running
system update commands on Linux systems.

The application is designed to simplify system updates while
still using the native package management tools of your
Linux distribution.

SysUpdate runs update commands using **pkexec**, which allows
the program to securely request administrator privileges
only when required.

The program is distributed as an **AppImage**, allowing it
to run on most Linux systems without a traditional install
process.

---

## Features

• Simple graphical interface for system updates
• Uses native package managers
• Secure privilege escalation using pkexec
• Distributed as a portable AppImage
• Easy install and uninstall scripts
• No system-wide installation required

---

## Requirements

SysUpdate requires the following components:

1. pkexec (PolicyKit)
2. FUSE support for AppImage

Most Linux distributions already include pkexec.

Example installations:

Fedora:
sudo dnf install polkit

Debian / Ubuntu:
sudo apt install policykit-1

FUSE support may also be required depending on your system.

Fedora:
sudo dnf install fuse

Debian / Ubuntu:
sudo apt install fuse

---

## Installation

SysUpdate installs to the user's home directory and does
not require administrator privileges.

Installed locations:

```
~/.local/bin/SysUpdate.AppImage
~/.local/share/applications/com.bytesbreadbbq.sysupdate.desktop
~/.local/share/icons/hicolor/256x256/apps/sysupdate.png
```

To install:

1. Open a terminal in the folder containing the files.

2. Make the installer executable:

   chmod +x install.sh

3. Run the installer:

   ./install.sh

After installation, SysUpdate will appear in your desktop
environment's application menu.

---

## Running SysUpdate

Launch SysUpdate from your desktop application menu.

When an operation requires administrator privileges,
the system will prompt for authentication using pkexec.

---

## Uninstalling

To remove SysUpdate:

1. Make the uninstall script executable (if needed):

   chmod +x uninstall.sh

2. Run the uninstall script:

   ./uninstall.sh

This removes the AppImage, launcher, and icon.

---

## Legacy Installations

Older versions of SysUpdate installed files into /usr/local.

If a legacy installation is detected, the install or uninstall
script may offer to remove it. Administrator permission is
required for this step.

---

## Support

For help or bug reports, please open an issue on GitHub:

https://github.com/RossContino1/SysUpdate
