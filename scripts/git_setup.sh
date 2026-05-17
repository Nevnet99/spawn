#!/usr/bin/env bash

set -e

SPAWN_DRY_RUN=${SPAWN_DRY_RUN:-false}
SPAWN_WORKSPACE=${SPAWN_WORKSPACE:-"personal"}
SPAWN_NAME=${SPAWN_NAME:-"Luke Brannagan"}
SPAWN_PERSONAL_EMAIL=${SPAWN_PERSONAL_EMAIL:-"luke@personal.dev"}
SPAWN_WORK_EMAIL=${SPAWN_WORK_EMAIL:-""}

PROJECTS_DIR="$HOME/Projects"
GLOBAL_GITCONFIG="$HOME/.gitconfig"

write_config_file() {
    local target_path=$1
    local content=$2

    if [ "$SPAWN_DRY_RUN" = "true" ]; then
        echo -e "🪵  [DRY-RUN] Would write to: $target_path\n---"
        echo -e "$content"
        echo -e "---\n"
    else
        echo -e "$content" > "$target_path"
        echo "✅ Created/Updated: $target_path"
    fi
}

make_directory() {
    local target_dir=$1
    if [ "$SPAWN_DRY_RUN" = "true" ]; then
        echo "🪵  [DRY-RUN] Would create directory: $target_dir"
    else
        mkdir -p "$target_dir"
        echo "✅ Verified directory: $target_dir"
    fi
}

echo "🛠️  Running Git Environment Spawning..."

if [ "$SPAWN_WORKSPACE" = "work" ]; then
    echo "💼 Designing a split Work/Personal Git environment layout..."

    make_directory "$PROJECTS_DIR/Work"
    make_directory "$PROJECTS_DIR/Personal"

    PERSONAL_CONFIG=$(cat << EOF
[user]
    name = $SPAWN_NAME
    email = $SPAWN_PERSONAL_EMAIL
EOF
)
    write_config_file "$HOME/.gitconfig-personal" "$PERSONAL_CONFIG"

    WORK_CONFIG=$(cat << EOF
[user]
    name = $SPAWN_NAME
    email = $SPAWN_WORK_EMAIL
EOF
)
    write_config_file "$HOME/.gitconfig-work" "$WORK_CONFIG"

    GLOBAL_CONFIG=$(cat << EOF
[core]
    editor = nvim
    excludesfile = ~/.gitignore_global

[includeIf "gitdir:~/Projects/Personal/"]
    path = ~/.gitconfig-personal

[includeIf "gitdir:~/Projects/Work/"]
    path = ~/.gitconfig-work
EOF
)
    write_config_file "$GLOBAL_GITCONFIG" "$GLOBAL_CONFIG"

else
    echo "🏡 Designing a flat, unified Personal Git environment layout..."
    
    make_directory "$PROJECTS_DIR"

    GLOBAL_CONFIG=$(cat << EOF
[user]
    name = $SPAWN_NAME
    email = $SPAWN_PERSONAL_EMAIL
[core]
    editor = nvim
    excludesfile = ~/.gitignore_global
EOF
)
    write_config_file "$GLOBAL_GITCONFIG" "$GLOBAL_CONFIG"
fi

# --- Global Ignore Sync ---
GLOBAL_IGNORE=$(cat << EOF
# macOS System Garbage
.DS_Store
.AppleDouble
.LSOverride

# Local IDE/Editor configs
.idea/
.vscode/
*.suo
*.swp

# Logs and Runtime outputs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
EOF
)

write_config_file "$HOME/.gitignore_global" "$GLOBAL_IGNORE"
echo "🎉 Git identity script tasks complete."
