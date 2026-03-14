#!/usr/bin/env bash
set -e

APP_NAME="SysUpdate"
APPIMAGE_NAME="SysUpdate-x86_64.AppImage"
DESKTOP_NAME="sysupdate.desktop"
ICON_NAME="sysupdate"
VERSION="0.1.0"

SOURCE_DIR="$(cd "$(dirname "$0")" && pwd)"
APPIMAGE_SOURCE="$SOURCE_DIR/$APPIMAGE_NAME"

INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR="$HOME/.local/share/icons/hicolor/256x256/apps"

APPIMAGE_TARGET="$INSTALL_DIR/SysUpdate.AppImage"
DESKTOP_TARGET="$DESKTOP_DIR/$DESKTOP_NAME"
ICON_TARGET="$ICON_DIR/$ICON_NAME.png"

echo "Installing $APP_NAME..."

if [ ! -f "$APPIMAGE_SOURCE" ]; then
    echo "Error: $APPIMAGE_NAME not found in:"
    echo "  $SOURCE_DIR"
    exit 1
fi

mkdir -p "$INSTALL_DIR"
mkdir -p "$DESKTOP_DIR"
mkdir -p "$ICON_DIR"

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

if [ -f "$TMP_DIR/squashfs-root/AppIcon.png" ]; then
    cp -f "$TMP_DIR/squashfs-root/AppIcon.png" "$ICON_TARGET"
elif [ -f "$TMP_DIR/squashfs-root/usr/share/icons/hicolor/256x256/apps/AppIcon.png" ]; then
    cp -f "$TMP_DIR/squashfs-root/usr/share/icons/hicolor/256x256/apps/AppIcon.png" "$ICON_TARGET"
else
    echo "Warning: could not extract icon from AppImage."
    echo "The app may still work, but the dock/menu icon may be generic."
fi

echo "Creating desktop entry..."
cat > "$DESKTOP_TARGET" <<EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=$APP_NAME
Comment=Linux system update utility
Exec=$APPIMAGE_TARGET
Icon=$ICON_NAME
Terminal=false
Categories=Utility;
StartupWMClass=SysUpdate
X-AppImage-Version=$VERSION
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
echo "Launch it from your application menu for proper GNOME/KDE dock integration."