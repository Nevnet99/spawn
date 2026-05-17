#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}
SPAWN_PERSONAL_EMAIL=${SPAWN_PERSONAL_EMAIL:-"luke@personal.dev"}

echo "🔐 INITIALIZING SSH SECURITY GUARD..."
echo "--------------------------------------------------------"

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would generate Ed25519 SSH key for $SPAWN_PERSONAL_EMAIL"
    echo "🪵  [DRY-RUN] Would configure ~/.ssh/config for macOS keychain integration"
    echo "🪵  [DRY-RUN] Would start ssh-agent and add identity securely"
    echo "🪵  [DRY-RUN] Would execute 'pbcopy' to place the public key in your system clipboard"
else

    if [ ! -f "$HOME/.ssh/id_ed25519" ]; then
        echo "⏳ Generating new Ed25519 SSH identity..."
        ssh-keygen -t ed25519 -C "$SPAWN_PERSONAL_EMAIL" -f "$HOME/.ssh/id_ed25519" -q -N ""
    else
        echo "✅ Ed25519 SSH key already exists. Skipping generation."
    fi

    echo "⏳ Configuring macOS Keychain integration..."
    touch "$HOME/.ssh/config"
    if ! grep -q "UseKeychain yes" "$HOME/.ssh/config"; then
        cat << 'EOF' >> "$HOME/.ssh/config"

Host *
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_ed25519
EOF
    fi

    echo "⏳ Adding key to ssh-agent..."
    eval "$(ssh-agent -s)" &>/dev/null
    
    ssh-add --apple-use-keychain "$HOME/.ssh/id_ed25519" &>/dev/null || ssh-add "$HOME/.ssh/id_ed25519" &>/dev/null

    if command -v pbcopy &> /dev/null; then
        cat "$HOME/.ssh/id_ed25519.pub" | pbcopy
        echo "--------------------------------------------------------"
        echo "✨ SUCCESS! Your new SSH public key was just copied to your clipboard! ✨"
        echo ""
        echo "🔗 NEXT STEPS:"
        echo "   1. Open this URL: https://github.com/settings/keys"
        echo "   2. Click the green 'New SSH key' button."
        echo "   3. Give it a title (e.g., 'Spawned MacBook') and hit CMD+V to paste."
        echo "--------------------------------------------------------"
    else
        echo "--------------------------------------------------------"
        echo "✨ SUCCESS! SSH key generated and secured."
        echo "   Please manually copy your key by running:"
        echo "   cat ~/.ssh/id_ed25519.pub"
        echo "--------------------------------------------------------"
    fi
fi
