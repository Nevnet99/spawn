#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

TARGET_RC="$HOME/.zshrc"
STARSHIP_SOURCE="$REPO_ROOT/configs/starship.toml"
STARSHIP_TARGET="$HOME/.config/starship.toml"

echo "⚙️  Synchronizing Shell Environment profiles & system prompts..."

define_aliases_and_runtimes() {
    cat << 'EOF'

# --- Modern CLI Tool Drop-in Replacements ---
alias ls="eza --icons --group-directories-first"
alias ll="eza -lh --icons --grid"
alias cat="bat"
alias g="git"

# --- Zoxide Initialisation ---
if command -v zoxide &> /dev/null; then
    eval "$(zoxide init zsh)"
fi

# --- Starship Prompt Initialisation ---
if command -v starship &> /dev/null; then
    eval "$(starship init zsh)"
fi

# --- Node Version Manager (NVM) Runtime Hooks ---
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
EOF
}

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would append modern aliases and runtime engine environments to $TARGET_RC"
    echo "🪵  [DRY-RUN] Would symlink Starship profile: $STARSHIP_SOURCE -> $STARSHIP_TARGET"
else
    # 1. Handle the profile appends
    touch "$TARGET_RC"
    if grep -q "Starship Prompt Initialisation" "$TARGET_RC"; then
        echo "✅ Shell profiles already mapped out in $TARGET_RC."
    else
        echo "📝 Injecting script hooks into $TARGET_RC..."
        define_aliases_and_runtimes >> "$TARGET_RC"
    fi

    # 2. Link your Starship config file cleanly
    mkdir -p "$HOME/.config"
    ln -sf "$STARSHIP_SOURCE" "$STARSHIP_TARGET"
    echo "✅ Starship environment configurations symlinked."
fi
