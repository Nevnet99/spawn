#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}

echo "💥 INITIALIZING SCORCHED EARTH SECURITY PURGE..."
echo "--------------------------------------------------------"

# Target configuration and data paths
TARGET_NVIM="$HOME/.config/nvim"
TARGET_GHOSTTY="$HOME/.config/ghostty"
TARGET_STARSHIP="$HOME/.config/starship.toml"
TARGET_ZSHRC="$HOME/.zshrc"
TARGET_PROJECTS="$HOME/Projects"

# Universal Identity & Security Targets
TARGET_GITCONFIG="$HOME/.gitconfig"
TARGET_SSH_DIR="$HOME/.ssh"

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would vaporize all local code repositories: $TARGET_PROJECTS"
    echo "🪵  [DRY-RUN] Would purge global Git identity parameters and signing rules: $TARGET_GITCONFIG"
    echo "🪵  [DRY-RUN] Would shred and delete all SSH identity keys and access credentials: $TARGET_SSH_DIR"
    echo "🪵  [DRY-RUN] Would flush your active macOS SSH keychain memory routing tables"
    echo "🪵  [DRY-RUN] Would purge active NvChad repository core: $TARGET_NVIM"
    echo "🪵  [DRY-RUN] Would purge active Ghostty symlinks: $TARGET_GHOSTTY"
    echo "🪵  [DRY-RUN] Would purge active Starship profile paths: $TARGET_STARSHIP"
    echo "🪵  [DRY-RUN] Would reset primary shell configuration scripts: $TARGET_ZSHRC"
else
    echo "⚠️  CRITICAL: Executing irrevocable cryptographic device wipe..."

    # 1. Identity & Commit Prevention Layer
    if [ -f "$TARGET_GITCONFIG" ]; then
        echo "🔥 Obliterating global Git configurations (~/.gitconfig)..."
        rm -f "$TARGET_GITCONFIG"
    fi

    # 2. Authentication & SSH Access Layer
    if [ -d "$TARGET_SSH_DIR" ]; then
        echo "🔥 Shredding and deleting secure identity keys (~/.ssh)..."
        # Optional: Overwrite files before unlinking so they cannot be easily recovered from the disk
        find "$TARGET_SSH_DIR" -type f -exec rm -P {} \; 2>/dev/null || true
        rm -rf "$TARGET_SSH_DIR"
    fi

    # Flush the active macOS system ssh-agent memory cache instantly
    echo "🔥 Purging live SSH keys from memory cache..."
    ssh-add -D &>/dev/null || true

    # 3. Local Workspace Layer
    if [ -d "$TARGET_PROJECTS" ]; then
        echo "🔥 Vaporizing local codebase directories inside $TARGET_PROJECTS..."
        rm -rf "$TARGET_PROJECTS"
    fi

    # 4. TUI Editor & Workspace Configuration States
    if [ -d "$TARGET_NVIM" ] || [ -L "$TARGET_NVIM" ]; then
        echo "🔥 Purging Neovim configurations..."
        rm -rf "$TARGET_NVIM"
    fi

    if [ -d "$TARGET_GHOSTTY" ] || [ -L "$TARGET_GHOSTTY" ]; then
        echo "🔥 Purging Ghostty target profiles..."
        rm -rf "$TARGET_GHOSTTY"
    fi

    if [ -f "$TARGET_STARSHIP" ] || [ -L "$TARGET_STARSHIP" ]; then
        echo "🔥 Purging Starship prompt configuration profiles..."
        rm -f "$TARGET_STARSHIP"
    fi

    if [ -f "$TARGET_ZSHRC" ]; then
        echo "🔥 Resetting system .zshrc baseline shell configuration..."
        true > "$TARGET_ZSHRC"
    fi

    echo "--------------------------------------------------------"
    echo "🔒 Security purge completed. No local identity footprint or access tokens remain."
fi
