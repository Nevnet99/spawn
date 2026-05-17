#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}

# Locate the absolute path of your spawn directory root
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

SOURCE_CONFIG="$REPO_ROOT/configs/ghostty"
TARGET_DIR="$HOME/.config/ghostty"
TARGET_CONFIG="$TARGET_DIR/config"

echo "👻 Synchronizing Ghostty symlink connections..."

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would create directory: $TARGET_DIR"
    echo "🪵  [DRY-RUN] Would symlink: $SOURCE_CONFIG -> $TARGET_CONFIG"
else
    # Guarantee the configuration housing folder exists
    mkdir -p "$TARGET_DIR"

    # Safely clear out any existing standalone files so the symlink doesn't collision-fail
    if [ -f "$TARGET_CONFIG" ] && [ ! -L "$TARGET_CONFIG" ]; then
        echo "📝 Backing up existing native Ghostty file to config.bak..."
        mv "$TARGET_CONFIG" "$TARGET_CONFIG.bak"
    fi

    # Create the active symbolic connection link
    ln -sf "$SOURCE_CONFIG" "$TARGET_CONFIG"
    echo "✅ Ghostty symlink established beautifully."
fi
