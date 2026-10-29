#!/bin/sh
set -e

DIR="$(cd "$(dirname "$0")" && pwd)"

if [ "$1" = "--user" ] || [ "$(id -u)" -ne 0 ]; then
    BIN_DIR="$HOME/.local/bin"
    APP_DIR="$HOME/.local/share/applications"
    ICON_DIR="$HOME/.local/share/icons/hicolor"
    echo "Installing Antlion locally for $(whoami)..."
else
    BIN_DIR="/usr/local/bin"
    APP_DIR="/usr/local/share/applications"
    ICON_DIR="/usr/local/share/icons/hicolor"
    echo "Installing Antlion system-wide..."
fi

mkdir -p "$BIN_DIR" "$APP_DIR" "$ICON_DIR/scalable/apps" "$ICON_DIR/512x512/apps"

install -m 755 "$DIR/antlion-app" "$BIN_DIR/antlion"
install -m 644 "$DIR/antlion-app.desktop" "$APP_DIR/antlion.desktop"
install -m 644 "$DIR/antlion.svg" "$ICON_DIR/scalable/apps/antlion.svg"
install -m 644 "$DIR/antlion-512.png" "$ICON_DIR/512x512/apps/antlion.png"

if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$APP_DIR" 2>/dev/null || true
fi

echo "Antlion installed successfully to $BIN_DIR/antlion"
