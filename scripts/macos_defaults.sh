#!/usr/bin/env bash

set -e
SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}

echo "🍏 Calibrating native macOS engine latency parameters..."

if [ "$SPAWN_DRY_RUN" = "true" ]; then
    echo "🪵  [DRY-RUN] Would set ApplePressAndHoldEnabled to false"
    echo "🪵  [DRY-RUN] Would set InitialKeyRepeat to 15"
    echo "🪵  [DRY-RUN] Would set KeyRepeat to 2"
    echo "🪵  [DRY-RUN] Would disable 'Natural Scrolling' (Sets trackpad direction to classic sliding)"
else
    defaults write -g ApplePressAndHoldEnabled -bool false
    defaults write -g InitialKeyRepeat -int 15
    defaults write -g KeyRepeat -int 2

    defaults write -g com.apple.swipescrolldirection -bool false
    
    echo "✅ OS acceleration values and trackpad direction profiles locked down."
    echo "   (Note: Trackpad updates will fully apply to active window panes upon system restart)."
fi
