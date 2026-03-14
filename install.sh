#!/usr/bin/env bash
set -euo pipefail

APP_NAME="SysUpdate"
APPIMAGE_NAME="SysUpdate-0.1.0-x86_64.AppImage"
DESKTOP_NAME="com.bytesbreadbbq.sysupdate.desktop"
VERSION="0.1.0"

SOURCE_DIR="$(cd "$(dirname "$0")" && pwd)"
APPIMAGE_SOURCE="$SOURCE_DIR/$APPIMAGE_NAME"

INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR="$HOME/.local/share/icons/hicolor/256x256/apps"

APPIMAGE_TARGET="$INSTALL_DIR/SysUpdate.AppImage"
DESKTOP_TARGET="$DESKTOP_DIR/$DESKTOP_NAME"
ICON_TARGET="$ICON_DIR/sysupdate.png"

LEGACY_BIN="/usr/local/bin/sysupdate-gui"
LEGACY_DESKTOP="/usr/local/share/applications/com.bytesbreadbbq.sysupdate.desktop"
LEGACY_ICON="/usr/local/share/pixmaps/com.bytesbreadbbq.sysupdate.png"

prompt_remove_legacy() {
    if [ -f "$LEGACY_BIN" ] || [ -f "$LEGACY_DESKTOP" ] || [ -f "$LEGACY_ICON" ]; then
        echo
        echo "A legacy system-wide SysUpdate install was found."
        echo "Older versions used /usr/local, which requires administrator permission to remove."
        echo "Without removing it, GNOME or KDE may keep showing or launching the older version."
        echo
        read -r -p "Remove legacy system-wide files now? [y/N] " reply
        case "$reply" in
            [yY]|[yY][eE][sS])
                sudo rm -f "$LEGACY_BIN" "$LEGACY_DESKTOP" "$LEGACY_ICON"
                if command -v update-desktop-database >/dev/null 2>&1; then
                    sudo update-desktop-database /usr/local/share/applications >/dev/null 2>&1 || true
                fi
                ;;
            *)
                echo
                echo "Skipping legacy cleanup."
                echo "You may still see the old menu entry until these are removed:"
                echo "  sudo rm -f \"$LEGACY_BIN\" \"$LEGACY_DESKTOP\" \"$LEGACY_ICON\""
                echo
                ;;
        esac
    fi
}

refresh_user_desktop() {
    if command -v update-desktop-database >/dev/null 2>&1; then
        update-desktop-database "$DESKTOP_DIR" >/dev/null 2>&1 || true
    fi

    if command -v gtk-update-icon-cache >/dev/null 2>&1; then
        gtk-update-icon-cache "$HOME/.local/share/icons/hicolor" >/dev/null 2>&1 || true
    fi
}

extract_icon() {
    local tmp_dir
    local old_pwd

    tmp_dir="$(mktemp -d)"
    old_pwd="$(pwd)"

    cleanup() {
        cd "$old_pwd" >/dev/null 2>&1 || true
        rm -rf "$tmp_dir"
    }

    trap cleanup RETURN

    cd "$tmp_dir"

    if ! "$APPIMAGE_TARGET" --appimage-extract >/dev/null 2>&1; then
        echo "Warning: could not extract AppImage contents for icon setup."
        echo "The app may still work, but the menu/dock icon may be generic."
        return 0
    fi

    if [ -f "$tmp_dir/squashfs-root/sysupdate.png" ]; then
        cp -f "$tmp_dir/squashfs-root/sysupdate.png" "$ICON_TARGET"
    elif [ -f "$tmp_dir/squashfs-root/AppIcon.png" ]; then
        cp -f "$tmp_dir/squashfs-root/AppIcon.png" "$ICON_TARGET"
    elif [ -f "$tmp_dir/squashfs-root/.DirIcon" ]; then
        cp -f "$tmp_dir/squashfs-root/.DirIcon" "$ICON_TARGET"
    elif [ -f "$tmp_dir/squashfs-root/usr/share/icons/hicolor/256x256/apps/sysupdate.png" ]; then
        cp -f "$tmp_dir/squashfs-root/usr/share/icons/hicolor/256x256/apps/sysupdate.png" "$ICON_TARGET"
    else
        echo "Warning: could not find an icon inside the AppImage."
        echo "The app may still work, but the menu/dock icon may be generic."
    fi
}

create_desktop_file() {
    cat > "$DESKTOP_TARGET" <<EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=$APP_NAME
Comment=Linux system update utility
Exec=$APPIMAGE_TARGET
Icon=$ICON_TARGET
Terminal=false
Categories=System;
StartupWMClass=SysUpdate
X-AppImage-Name=$APP_NAME
X-AppImage-Version=$VERSION
EOF

    chmod +x "$DESKTOP_TARGET"
}

echo "Installing $APP_NAME..."

if [ ! -f "$APPIMAGE_SOURCE" ]; then
    echo "Error: $APPIMAGE_NAME not found in:"
    echo "  $SOURCE_DIR"
    exit 1
fi

mkdir -p "$INSTALL_DIR" "$DESKTOP_DIR" "$ICON_DIR"

prompt_remove_legacy

echo "Copying AppImage..."
cp -f "$APPIMAGE_SOURCE" "$APPIMAGE_TARGET"
chmod +x "$APPIMAGE_TARGET"

echo "Extracting icon from AppImage..."
extract_icon

echo "Creating desktop entry..."
create_desktop_file

refresh_user_desktop

echo
echo "$APP_NAME installed successfully."
echo "Launch it from the application menu for proper dock integration."