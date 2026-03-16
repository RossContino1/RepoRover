#!/usr/bin/env bash
set -e

APP_NAME="RepoRover"
APPIMAGE_NAME="RepoRover-1.0.0-x86_64.AppImage"
DESKTOP_NAME="reporover.desktop"
ICON_NAME="reporover"
VERSION="1.0.0"

SOURCE_DIR="$(cd "$(dirname "$0")" && pwd)"
APPIMAGE_SOURCE="$SOURCE_DIR/$APPIMAGE_NAME"

INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR_256="$HOME/.local/share/icons/hicolor/256x256/apps"
ICON_DIR_SCALABLE="$HOME/.local/share/icons/hicolor/scalable/apps"

APPIMAGE_TARGET="$INSTALL_DIR/RepoRover.AppImage"
DESKTOP_TARGET="$DESKTOP_DIR/$DESKTOP_NAME"
ICON_TARGET_256="$ICON_DIR_256/$ICON_NAME.png"
ICON_TARGET_SCALABLE="$ICON_DIR_SCALABLE/$ICON_NAME.png"

OLD_APPIMAGE_TARGET="$INSTALL_DIR/SysUpdate.AppImage"
OLD_DESKTOP_TARGET="$DESKTOP_DIR/sysupdate.desktop"
OLD_ICON_TARGET_256="$ICON_DIR_256/sysupdate.png"
OLD_ICON_TARGET_SCALABLE="$ICON_DIR_SCALABLE/sysupdate.png"

echo "Installing $APP_NAME..."

if [ ! -f "$APPIMAGE_SOURCE" ]; then
    echo "Error: $APPIMAGE_NAME not found in:"
    echo "  $SOURCE_DIR"
    exit 1
fi

mkdir -p "$INSTALL_DIR"
mkdir -p "$DESKTOP_DIR"
mkdir -p "$ICON_DIR_256"
mkdir -p "$ICON_DIR_SCALABLE"

echo "Removing old SysUpdate / RepoRover files..."
rm -f "$OLD_APPIMAGE_TARGET"
rm -f "$OLD_DESKTOP_TARGET"
rm -f "$OLD_ICON_TARGET_256"
rm -f "$OLD_ICON_TARGET_SCALABLE"

rm -f "$APPIMAGE_TARGET"
rm -f "$DESKTOP_TARGET"
rm -f "$ICON_TARGET_256"
rm -f "$ICON_TARGET_SCALABLE"

chmod +x "$APPIMAGE_SOURCE"

echo "Copying AppImage..."
cp -f "$APPIMAGE_SOURCE" "$APPIMAGE_TARGET"
chmod +x "$APPIMAGE_TARGET"

echo "Extracting icon from AppImage..."
TMP_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

cd "$TMP_DIR"
"$APPIMAGE_TARGET" --appimage-extract >/dev/null 2>&1 || true

ICON_FOUND=""

if [ -f "$TMP_DIR/squashfs-root/AppIcon.png" ]; then
    ICON_FOUND="$TMP_DIR/squashfs-root/AppIcon.png"
elif [ -f "$TMP_DIR/squashfs-root/usr/share/icons/hicolor/256x256/apps/AppIcon.png" ]; then
    ICON_FOUND="$TMP_DIR/squashfs-root/usr/share/icons/hicolor/256x256/apps/AppIcon.png"
elif [ -f "$TMP_DIR/squashfs-root/reporover.png" ]; then
    ICON_FOUND="$TMP_DIR/squashfs-root/reporover.png"
fi

if [ -n "$ICON_FOUND" ]; then
    cp -f "$ICON_FOUND" "$ICON_TARGET_256"
    cp -f "$ICON_FOUND" "$ICON_TARGET_SCALABLE"
else
    echo "Warning: could not extract icon from AppImage."
fi

echo "Creating desktop entry..."
cat > "$DESKTOP_TARGET" <<EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=$APP_NAME
GenericName=System Updater
Comment=Universal Linux package manager updater
Exec=$APPIMAGE_TARGET
Icon=$ICON_NAME
Terminal=false
Categories=System;Utility;
StartupWMClass=RepoRover
Keywords=linux;update;package-manager;apt;dnf;pacman;zypper;
Actions=RunUpdate;
X-AppImage-Version=$VERSION

[Desktop Action RunUpdate]
Name=Run System Update
Exec=$APPIMAGE_TARGET
EOF

chmod +x "$DESKTOP_TARGET"

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$DESKTOP_DIR" >/dev/null 2>&1 || true
fi

if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache "$HOME/.local/share/icons/hicolor" >/dev/null 2>&1 || true
fi

echo
echo "$APP_NAME installed successfully."
echo "Launch it from your application menu for proper desktop integration."