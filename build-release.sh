#!/usr/bin/env bash
set -Eeuo pipefail

APP_NAME="RepoRover"
APP_ID="reporover"
VERSION="1.2.2"
ARCH="x86_64"

BINARY_NAME="$APP_ID"
APPDIR="${APP_NAME}.AppDir"
APPIMAGE_NAME="${APP_NAME}-${VERSION}-${ARCH}.AppImage"
RELEASE_DIR_NAME="${APP_NAME}-${VERSION}-linux-${ARCH}"
ZIP_NAME="${RELEASE_DIR_NAME}.zip"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$ROOT_DIR/build"
DIST_DIR="$ROOT_DIR/$RELEASE_DIR_NAME"
APPIMAGE_TOOL="$ROOT_DIR/appimagetool-${ARCH}.AppImage"

DESKTOP_FILE="$ROOT_DIR/${APP_ID}.desktop"
ICON_FILE="$ROOT_DIR/${APP_ID}.png"
README_FILE="$ROOT_DIR/README-LINUX.txt"
LICENSE_FILE="$ROOT_DIR/LICENSE.txt"
INSTALL_FILE="$ROOT_DIR/install.sh"
UNINSTALL_FILE="$ROOT_DIR/uninstall.sh"

log() {
    echo "==> $*"
}

die() {
    echo "Error: $*" >&2
    exit 1
}

require_cmd() {
    command -v "$1" >/dev/null 2>&1 || die "Required command not found: $1"
}

cleanup_old_artifacts() {
    log "Cleaning old build artifacts..."
    rm -rf "$BUILD_DIR" "$ROOT_DIR/$APPDIR" "$DIST_DIR"
    rm -f "$ROOT_DIR/$APPIMAGE_NAME" "$ROOT_DIR/$ZIP_NAME"
}

check_required_files() {
    log "Checking required files..."

    [[ -f "$ROOT_DIR/main.go" ]] || die "main.go not found in repo root"
    [[ -f "$ROOT_DIR/go.mod" ]] || die "Missing go.mod"
    [[ -f "$DESKTOP_FILE" ]] || die "Missing desktop file: ${APP_ID}.desktop"
    [[ -f "$ICON_FILE" ]] || die "Missing icon file: ${APP_ID}.png"
    [[ -f "$README_FILE" ]] || die "Missing README.md"
    [[ -f "$LICENSE_FILE" ]] || die "Missing LICENSE.txt"
    [[ -f "$INSTALL_FILE" ]] || die "Missing install.sh"
    [[ -f "$UNINSTALL_FILE" ]] || die "Missing uninstall.sh"
}

check_tools() {
    log "Checking required tools..."
    require_cmd go
    require_cmd chmod
    require_cmd cp
    require_cmd rm
    require_cmd mkdir
    require_cmd zip

    if [[ ! -f "$APPIMAGE_TOOL" ]]; then
        require_cmd wget
    fi
}

download_appimagetool_if_needed() {
    if [[ -f "$APPIMAGE_TOOL" ]]; then
        log "Using existing $(basename "$APPIMAGE_TOOL")"
        chmod +x "$APPIMAGE_TOOL"
        return
    fi

    log "Downloading appimagetool..."
    wget -O "$APPIMAGE_TOOL" \
        "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-${ARCH}.AppImage"

    chmod +x "$APPIMAGE_TOOL"
}

prepare_directories() {
    log "Preparing build directories..."
    mkdir -p "$BUILD_DIR"
    mkdir -p "$ROOT_DIR/$APPDIR/usr/bin"
    mkdir -p "$ROOT_DIR/$APPDIR/usr/share/applications"
    mkdir -p "$DIST_DIR"
}

build_binary() {
    log "Tidying Go modules..."
    (
        cd "$ROOT_DIR"
        go mod tidy
    )

    log "Building release binary..."
    (
        cd "$ROOT_DIR"
        CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o "$BUILD_DIR/$BINARY_NAME"
    )

    [[ -f "$BUILD_DIR/$BINARY_NAME" ]] || die "Build failed: binary not created"

    chmod +x "$BUILD_DIR/$BINARY_NAME"
}

create_apprun() {
    cat > "$ROOT_DIR/$APPDIR/AppRun" <<EOF
#!/usr/bin/env bash
HERE="\$(dirname "\$(readlink -f "\$0")")"
exec "\$HERE/usr/bin/$BINARY_NAME" "\$@"
EOF
    chmod +x "$ROOT_DIR/$APPDIR/AppRun"
}

create_appdir() {
    log "Creating AppDir..."

    cp "$BUILD_DIR/$BINARY_NAME" "$ROOT_DIR/$APPDIR/usr/bin/$BINARY_NAME"
    cp "$DESKTOP_FILE" "$ROOT_DIR/$APPDIR/${APP_ID}.desktop"
    cp "$DESKTOP_FILE" "$ROOT_DIR/$APPDIR/usr/share/applications/${APP_ID}.desktop"
    cp "$ICON_FILE" "$ROOT_DIR/$APPDIR/${APP_ID}.png"

    create_apprun

    chmod +x "$ROOT_DIR/$APPDIR/usr/bin/$BINARY_NAME"
}

build_appimage() {
    log "Building AppImage..."

    export ARCH
    "$APPIMAGE_TOOL" "$ROOT_DIR/$APPDIR" "$ROOT_DIR/$APPIMAGE_NAME"

    [[ -f "$ROOT_DIR/$APPIMAGE_NAME" ]] || die "AppImage build failed"
    chmod +x "$ROOT_DIR/$APPIMAGE_NAME"
}

prepare_distribution() {
    log "Preparing release folder..."

    cp "$ROOT_DIR/$APPIMAGE_NAME" "$DIST_DIR/"
    cp "$INSTALL_FILE" "$DIST_DIR/"
    cp "$UNINSTALL_FILE" "$DIST_DIR/"
    cp "$README_FILE" "$DIST_DIR/"
    cp "$LICENSE_FILE" "$DIST_DIR/"
    cp "$ICON_FILE" "$DIST_DIR/"
    cp "$DESKTOP_FILE" "$DIST_DIR/"

    chmod +x "$DIST_DIR/install.sh"
    chmod +x "$DIST_DIR/uninstall.sh"
    chmod +x "$DIST_DIR/$APPIMAGE_NAME"
}

create_zip() {
    log "Creating zip archive with top-level folder..."
    (
        cd "$ROOT_DIR"
        zip -r "$ZIP_NAME" "$RELEASE_DIR_NAME"
    )

    [[ -f "$ROOT_DIR/$ZIP_NAME" ]] || die "Zip creation failed"
}

print_summary() {
    echo
    echo "Build complete:"
    echo "  Binary:        $BUILD_DIR/$BINARY_NAME"
    echo "  AppDir:        $APPDIR/"
    echo "  AppImage:      $APPIMAGE_NAME"
    echo "  Release folder:$RELEASE_DIR_NAME/"
    echo "  Zip archive:   $ZIP_NAME"
    echo
    echo "Next test:"
    echo "  ./$APPIMAGE_NAME"
}

main() {
    cleanup_old_artifacts
    check_required_files
    check_tools
    download_appimagetool_if_needed
    prepare_directories
    build_binary
    create_appdir
    build_appimage
    prepare_distribution
    create_zip
    print_summary
}

main "$@"
