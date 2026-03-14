#!/usr/bin/env bash
set -euo pipefail

APP_NAME="SysUpdate"
DESKTOP_NAME="com.bytesbreadbbq.sysupdate.desktop"

INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR="$HOME/.local/share/icons/hicolor/256x256/apps"

APPIMAGE_TARGET="$INSTALL_DIR/SysUpdate.AppImage"
DESKTOP_TARGET="$DESKTOP_DIR/$DESKTOP_NAME"
ICON_TARGET="$ICON_DIR/sysupdate.png"

refresh_user_desktop() {
    if command -v update-desktop-database >/dev/null 2>&1; then
        update-desktop-database "$DESKTOP_DIR" >/dev/null 2>&1 || true
    fi

    if command -v gtk-update-icon-cache >/dev/null 2>&1; then
        gtk-update-icon-cache "$HOME/.local/share/icons/hicolor" >/dev/null 2>&1 || true
    fi
}

echo "Uninstalling $APP_NAME..."

rm -f "$APPIMAGE_TARGET"
rm -f "$DESKTOP_TARGET"
rm -f "$ICON_TARGET"

refresh_user_desktop

echo
echo "$APP_NAME has been removed from this user account."