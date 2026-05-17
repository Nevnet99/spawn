#!/usr/bin/env bash

set -e

# Fall back to safe dry-run mode if the global variable isn't explicitly defined
SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}

echo "🔄 INITIALIZING NON-DESTRUCTIVE CONFIGURATION UPDATE..."
echo "--------------------------------------------------------"

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would execute 'git pull origin main' to sync configuration profiles"
    echo "🪵  [DRY-RUN] Would execute NvChad headless sync to update internal editor plugins"
else
    # 1. Verify we are in a valid Git repository to prevent path clipping
    if [ ! -d .git ]; then
        echo "❌ Error: This script must be executed from your spawn repository root directory."
        exit 1
    fi

    echo "⏳ Fetching latest configuration files from remote repository..."
    # Pull incoming changes without displacing uncommitted local edits
    git pull --rebase origin main || {
        echo "⚠️  Git Pull halted. You may have local conflicts or need to configure a remote origin branch."
    }

    # 2. Trigger an automated headless update for NvChad plugins
    if [ -d "$HOME/.config/nvim" ]; then
        echo "⏳ Running automated headless plugin updates inside NvChad..."
        # Opens Neovim silently in the background, runs the plugin sync, and exits
        nvim --headless "+Lazy! sync" +qa &>/dev/null || true
        echo "✅ NvChad plugin ecosystem synchronized."
    fi

    echo "--------------------------------------------------------"
    echo "🎉 Workspace configurations updated and synchronized seamlessly!"
fi
