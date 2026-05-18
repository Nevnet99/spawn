#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}

# Core utilities including Golang
FORMULAS=(
    "zoxide"
    "eza"
    "bat"
    "ripgrep"
    "fzf"
    "btop"
    "gh"
    "golang"
    "starship"
    "neovim"
    "zsh-autosuggestions"
    "zsh-syntax-highlighting"
    "zsh-you-should-use"
    "fzf-tab"
)

echo "💻 Deploying core terminal utility binaries and runtime languages..."

# Verify Homebrew is present
if ! command -v brew &> /dev/null; then
    if [ "$SPAWN_DRY_RUN" = "true" ]; then
        echo "🪵  [DRY-RUN] Homebrew missing. Would install before binary provisioning."
    else
        echo "⏳ Homebrew missing. Triggering deployment layer..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        eval "$(/opt/homebrew/bin/brew shellenv)"
    fi
else
    echo "✅ Homebrew framework verified."
fi

# 1. Install standard binary formulas (including Go)
for formula in "${FORMULAS[@]}"; do
    if [ "$SPAWN_DRY_RUN" = "true" ]; then
        echo "🪵  [DRY-RUN] brew install $formula"
    else
        echo "⏳ Installing utility: $formula..."
        brew install "$formula"
    fi
done

# 2. Deploy NVM via its official script runner safely
if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash"
else
    if [ -d "$HOME/.nvm" ]; then
        echo "✅ NVM environment root directory already present (skipping install)."
    else
        echo "⏳ Fetching and installing Node Version Manager (NVM)..."
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash
    fi
fi

echo "🎉 Terminal utilities, Go runtime, and NVM layout complete."
