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

# Deploy Cask Apps cleanly (with --force to overwrite existing apps)
for cask in "${CASKS[@]}"; do
    if [ "$SPAWN_DRY_RUN" = "true" ]; then
        echo "🪵  [DRY-RUN] brew install --cask --force $cask"
    else
        echo "⏳ Deploying Cask app: $cask..."
        brew install --cask --force "$cask"
    fi
done

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] brew tap webadderallorg/recordly && brew install --cask recordly"
else
    echo "⏳ Deploying Recordly from custom tap..."
    
    # 1. Tap the custom repository so Homebrew trusts it
    brew tap webadderallorg/recordly https://github.com/webadderallorg/Recordly.git || true
    
    # 2. Install the cask by name, not URL
    brew install --cask --force recordly || true
    
    echo "✅ Recordly deployment handled."
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