#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}

echo "📤 INITIALIZING UPSTREAM CONFIGURATION BACKUP..."
echo "--------------------------------------------------------"

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would detect modified configuration files across the repository"
    echo "🪵  [DRY-RUN] Would stage and commit changes with an automated timestamp message"
    echo "🪵  [DRY-RUN] Would push changes securely to origin main"
else
    # 1. Verify path integrity
    if [ ! -d .git ]; then
        echo "❌ Error: This script must be executed from your spawn repository root directory."
        exit 1
    fi

    # 2. Check if there are actual changes to commit
    if [ -z "$(git status --porcelain)" ]; then
        echo "✅ No configuration changes detected. Your cloud repository is already up to date."
    else
        echo "⏳ Staging modified configurations and scripts..."
        git add .
        
        # Create an automated commit message tagging the exact time of the backup
        TIMESTAMP=$(date +'%Y-%m-%d %H:%M')
        git commit -m "chore(sync): automated workspace configuration backup $TIMESTAMP"
        
        echo "⏳ Pushing updates to the cloud..."
        git push origin main
        echo "✅ Configurations successfully synchronized to the cloud!"
    fi
    echo "--------------------------------------------------------"
fi
