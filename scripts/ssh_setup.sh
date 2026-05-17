#!/usr/bin/env bash

set -e
SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}

echo "🔑 Setting up secure SSH key configurations..."

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would generate Ed25519 key pair at ~/.ssh/id_ed25519"
    echo "🪵  [DRY-RUN] Would configure local ~/.ssh/config for auto-loading keys into the macOS keychain"
else
    mkdir -p "$HOME/.ssh"
    chmod 700 "$HOME/.ssh"

    if [ -f "$HOME/.ssh/id_ed25519" ]; then
        echo "✅ SSH key pair already exists."
    else
        echo "⏳ Generating fresh Ed25519 key pair..."
        # Note: Generates key without a passphrase to keep script automated, 
        # remove the empty quotes if you want it to prompt you interactively.
        ssh-keygen -t ed25519 -C "$SPAWN_PERSONAL_EMAIL" -f "$HOME/.ssh/id_ed25519" -N ""
    fi

    # Append standard macOS keychain configuration helper rules
    cat << 'EOF' > "$HOME/.ssh/config"
Host *
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_ed25519
EOF
    chmod 600 "$HOME/.ssh/config"
    echo "✅ SSH config profiles initialized."
fi
