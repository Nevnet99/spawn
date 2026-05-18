#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

SOURCE_CUSTOM="$REPO_ROOT/configs/nvim/lua/custom"
TARGET_NVIM="$HOME/.config/nvim"
TARGET_CUSTOM="$TARGET_NVIM/lua/custom"
TARGET_BAK="$TARGET_NVIM.bak"

echo "🌌 Initializing high-velocity NvChad environment configurations..."

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would fetch upstream NvChad baseline repository to: $TARGET_NVIM"
    echo "🪵  [DRY-RUN] Would create custom configuration symlink pointer: $SOURCE_CUSTOM -> $TARGET_CUSTOM"
else
    if [ -d "$TARGET_NVIM" ] && [ ! -L "$TARGET_NVIM" ]; then
        # Clean up any stale backup directory to guarantee the 'mv' operation can't collide
        if [ -d "$TARGET_BAK" ]; then
            echo "🧹 Removing obsolete Neovim backup traces..."
            rm -rf "$TARGET_BAK"
        fi

        echo "📝 Moving existing Neovim directory folder out of the way to nvim.bak..."
        mv "$TARGET_NVIM" "$TARGET_BAK"
    fi

    if [ ! -d "$TARGET_NVIM" ]; then
        echo "⏳ Fetching pristine NvChad framework core wrapper..."
        git clone https://github.com/NvChad/starter "$TARGET_NVIM" --depth 1
    fi

    echo "🔗 Anchoring user configuration symlinks..."
    rm -rf "$TARGET_CUSTOM" 
    ln -sf "$SOURCE_CUSTOM" "$TARGET_CUSTOM"
    
    echo "✅ NvChad environment setup successfully synchronized."
fi