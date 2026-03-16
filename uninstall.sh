#!/usr/bin/env bash
set -e

APP_NAME="RepoRover"

INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR_256="$HOME/.local/share/icons/hicolor/256x256/apps"
ICON_DIR_SCALABLE="$HOME/.local/share/icons/hicolor/scalable/apps"

echo "Uninstalling $APP_NAME and removing old SysUpdate files..."

rm -f "$INSTALL_DIR/RepoRover.AppImage"
rm -f "$INSTALL_DIR/SysUpdate.AppImage"

rm -f "$DESKTOP_DIR/reporover.desktop"
rm -f "$DESKTOP_DIR/sysupdate.desktop"

rm -f "$ICON_DIR_256/reporover.png"
rm -f "$ICON_DIR_256/sysupdate.png"

rm -f "$ICON_DIR_SCALABLE/reporover.png"
rm -f "$ICON_DIR_SCALABLE/sysupdate.png"

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$DESKTOP_DIR" >/dev/null 2>&1 || true
fi

if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache "$HOME/.local/share/icons/hicolor" >/dev/null 2>&1 || true
fi

echo "$APP_NAME uninstalled."