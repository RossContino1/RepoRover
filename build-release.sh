#!/usr/bin/env bash
set -euo pipefail

APP_NAME="RepoRover"
APP_ID="reporover"
VERSION="1.0.0"
ARCH="x86_64"

BINARY_NAME="$APP_ID"
APPDIR="${APP_NAME}.AppDir"
APPIMAGE_NAME="${APP_NAME}-${VERSION}-${ARCH}.AppImage"
DIST_DIR="distribution"
ZIP_NAME="${APP_NAME}-${VERSION}-linux-${ARCH}.zip"

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
BUILD_DIR="$ROOT_DIR/build"

echo "==> Cleaning old build artifacts..."
rm -rf "$BUILD_DIR" "$APPDIR" "$DIST_DIR"
rm -f "$APPIMAGE_NAME" "$ZIP_NAME"

mkdir -p "$BUILD_DIR"
mkdir -p "$APPDIR/usr/bin"
mkdir -p "$APPDIR/usr/share/applications"
mkdir -p "$DIST_DIR"

echo "==> Building Go binary..."
go build -o "$BUILD_DIR/$BINARY_NAME"

chmod +x "$BUILD_DIR/$BINARY_NAME"

echo "==> Creating AppDir..."
cp "$BUILD_DIR/$BINARY_NAME" "$APPDIR/usr/bin/$BINARY_NAME"
cp "$APP_ID.desktop" "$APPDIR/$APP_ID.desktop"
cp "$APP_ID.desktop" "$APPDIR/usr/share/applications/$APP_ID.desktop"
cp "$APP_ID.png" "$APPDIR/$APP_ID.png"

cat > "$APPDIR/AppRun" <<EOF
#!/usr/bin/env bash
HERE="\$(dirname "\$(readlink -f "\$0")")"
exec "\$HERE/usr/bin/$BINARY_NAME" "\$@"
EOF

chmod +x "$APPDIR/AppRun"
chmod +x "$APPDIR/usr/bin/$BINARY_NAME"

echo "==> Building AppImage..."
if [ ! -f "$ROOT_DIR/appimagetool-x86_64.AppImage" ]; then
    echo "Error: appimagetool-x86_64.AppImage not found in repo root."
    exit 1
fi

chmod +x "$ROOT_DIR/appimagetool-x86_64.AppImage"
"$ROOT_DIR/appimagetool-x86_64.AppImage" "$APPDIR" "$APPIMAGE_NAME"

echo "==> Preparing distribution folder..."
cp "$APPIMAGE_NAME" "$DIST_DIR/"
cp install.sh "$DIST_DIR/"
cp uninstall.sh "$DIST_DIR/"
cp README.md "$DIST_DIR/"
cp LICENSE.txt "$DIST_DIR/"
cp "$APP_ID.png" "$DIST_DIR/"
cp "$APP_ID.desktop" "$DIST_DIR/"

chmod +x "$DIST_DIR/install.sh"
chmod +x "$DIST_DIR/uninstall.sh"
chmod +x "$DIST_DIR/$APPIMAGE_NAME"

echo "==> Creating zip archive..."
(
    cd "$DIST_DIR"
    zip -r "../$ZIP_NAME" .
)

echo
echo "Build complete:"
echo "  AppImage: $APPIMAGE_NAME"
echo "  Folder:   $DIST_DIR/"
echo "  Zip:      $ZIP_NAME"