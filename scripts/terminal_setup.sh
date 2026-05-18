#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Define paths cleanly
ZSHRC_SOURCE="$REPO_ROOT/configs/.zshrc"
ZSHRC_TARGET="$HOME/.zshrc"

STARSHIP_SOURCE="$REPO_ROOT/configs/starship.toml"
STARSHIP_TARGET="$HOME/.config/starship.toml"

echo "⚙️  Synchronizing Shell Environment profiles & system prompts..."
echo "--------------------------------------------------------"

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would symlink Zsh Profile: $ZSHRC_SOURCE -> $ZSHRC_TARGET"
    echo "🪵  [DRY-RUN] Would symlink Starship profile: $STARSHIP_SOURCE -> $STARSHIP_TARGET"
else
	# 1. Clean up potential dead/ghost configurations safely
	if [ -L "$ZSHRC_TARGET" ] || [ -f "$ZSHRC_TARGET" ]; then
		echo "⏳ Clearing previous shell profile targets..."
		rm -rf "$ZSHRC_TARGET"
	fi

    # 2. Wire up the .zshrc symlink
    echo "🔗 Symlinking fresh terminal runtime profiles..."
    ln -sf "$ZSHRC_SOURCE" "$ZSHRC_TARGET"
    echo "✅ Zsh operational runtime profiles successfully symlinked."

    # 3. Link your Starship config file cleanly
    mkdir -p "$HOME/.config"
    ln -sf "$STARSHIP_SOURCE" "$STARSHIP_TARGET"
    echo "✅ Starship environment configurations symlinked."
fi