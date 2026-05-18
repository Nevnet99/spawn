# --- Modern CLI Tool Drop-in Replacements ---
alias ls="eza --icons --group-directories-first"
alias ll="eza -lh --icons --grid"
alias cat="bat"
alias g="git"

# --- Git Shorthand Aliases (Oh My Zsh Port) ---
alias gst='git status'
alias ga='git add'
alias gaa='git add --all'
alias gco='git checkout'
alias gcb='git checkout -b'
alias gcmsg='git commit -m'
alias gl='git pull'
alias gp='git push'
alias ggp='git push origin "$(git branch --show-current)"'
alias ggpull='git pull origin "$(git branch --show-current)"'
alias glog='git log --oneline --decorate --graph'

# --- High-Velocity Autocompletion Engine ---
# Initialize the native Zsh completion system (compsys) with speed caching
autoload -Uz compinit
if [ $(date +'%j') -ne $(stat -f '%Sm' -t '%j' ~/.zcompdump 2>/dev/null || echo 0) ]; then
    compinit
else
    compinit -C
fi

# Set case-insensitive, substring, and partial matching logic for standard tab menus
zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}' 'r:|=*' 'l:|=* r:|=*'
zstyle ':completion:*' list-colors "${(s.:.)LS_COLORS}"
zstyle ':completion:*' menu select

# --- FZF Advanced History Integration ---
if command -v fzf &> /dev/null; then
    source <(fzf --zsh)
fi

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

# --- Advanced Interactive Plugin Suite ---
# Dynamically locate the Homebrew cell prefix layout for portable architecture executions
if command -v brew &>/dev/null; then
    BREW_PREFIX="$(brew --prefix)"
    
    # 1. FZF-Tab Hook (Must load BEFORE autocomplete scripts initialize)
    if [ -f "$BREW_PREFIX/share/fzf-tab/fzf-tab.plugin.zsh" ]; then
        source "$BREW_PREFIX/share/fzf-tab/fzf-tab.plugin.zsh"
    fi
    
    # 2. Fish-like Ghost Text Autosuggestions Hook
    if [ -f "$BREW_PREFIX/share/zsh-autosuggestions/zsh-autosuggestions.zsh" ]; then
        source "$BREW_PREFIX/share/zsh-autosuggestions/zsh-autosuggestions.zsh"
        
        # Bind Tab key to instantly accept the prediction layout line
        bindkey '^I' autosuggest-accept
    fi
    
    # 3. You-Should-Use Alias Reminder Engine
    if [ -f "$BREW_PREFIX/share/zsh-you-should-use/you-should-use.plugin.zsh" ]; then
        source "$BREW_PREFIX/share/zsh-you-should-use/you-should-use.plugin.zsh"
    fi
    
    # 4. Real-time Command Validation Syntax Highlighting
    # Note: Syntax Highlighting MUST be loaded dead-last in the execution stack to prevent script interception bugs
    if [ -f "$BREW_PREFIX/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh" ]; then
        source "$BREW_PREFIX/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
    fi
fi