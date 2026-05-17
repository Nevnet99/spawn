#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}

# Desktop Apps Only
CASKS=(
    "google-chrome"
    "spotify"
    "ghostty"
    "obsidian"
)

echo "📦 Commencing desktop application deployment..."

if ! command -v brew &> /dev/null; then
    if [ "$SPAWN_DRY_RUN" = "true" ]; then
        echo "🪵  [DRY-RUN] Homebrew missing. Would install before application provisioning."
    else
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        eval "$(/opt/homebrew/bin/brew shellenv)"
    fi
fi

# Deploy Cask Apps cleanly
for cask in "${CASKS[@]}"; do
    if [ "$SPAWN_DRY_RUN" = "true" ]; then
        echo "🪵  [DRY-RUN] brew install --cask $cask"
    else
        echo "⏳ Deploying Cask app: $cask..."
        brew install --cask "$cask"
    fi
done

# Deploy Recordly via Standalone Formula
RECORDLY_FORMULA_URL="https://raw.githubusercontent.com/webadderallorg/Recordly/main/recordly.rb"
if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] brew install --cask $RECORDLY_FORMULA_URL"
else
    echo "⏳ Deploying Recordly from remote formula source..."
    brew install --cask "$RECORDLY_FORMULA_URL"
fi

# Handle Default Browser configuration
if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] brew install defaultbrowser && defaultbrowser chrome"
else
    echo "🌐 Synchronizing default browser definitions..."
    if ! command -v defaultbrowser &> /dev/null; then
        brew install defaultbrowser
    fi
    defaultbrowser chrome
    echo "✅ Google Chrome configured as primary browser."
fi

echo "🎉 Desktop app provisioning complete."
