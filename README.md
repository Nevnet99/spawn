**Spawn** is a blazing fast, TUI-driven macOS workspace provisioning and configuration management tool. Written in Go and powered by Bubbletea, it automates the deployment of modern developer environments, terminal workflows, and application payloads.

## 🚀 Overview

Instead of maintaining brittle shell scripts or monolithic dotfiles, Spawn uses a modular, symlink-based architecture managed through a sleek terminal interface. It separates core execution logic from your user configurations, ensuring you can update, push, or completely vaporize your machine state with a single keystroke.

### Key Features
* **Interactive TUI:** Navigate deployment steps, toggle specific modules, and view real-time execution logs.
* **Symlink Architecture:** User configurations (Ghostty, Starship, NvChad) live in the repository and are securely linked to your OS.
* **Non-Destructive Cloud Syncing:** Built-in pipelines to push local config changes or pull remote updates headless.
* **Dry-Run Simulation:** Test your deployment sequences safely without altering your filesystem.
* **Scorched Earth Protocol:** A total cryptographic and local workspace kill-switch for immediate, secure offboarding.

---

## 📂 Architecture & File Tree

The project is structured to keep execution scripts independent from raw configuration data:

```text
└── spawn/
    ├── main.go                 <-- The core Go Bubbletea TUI application
    ├── configs/                <-- Plain-text config files (Your symlink source)
    │   ├── ghostty
    │   ├── starship.toml
    │   └── nvim/
    │       └── lua/custom/     <-- Modular NvChad user configurations
    └── scripts/                <-- Modular bash execution workers
        ├── git_setup.sh
        ├── os_packages.sh
        ├── cli_tools.sh
        ├── terminal_setup.sh
        ├── ssh_setup.sh
        ├── ghostty_setup.sh
        ├── nvim_setup.sh
        ├── macos_defaults.sh
        ├── update_workspace.sh
        ├── push_configs.sh
        └── scorched_earth.sh
```

## 📝 Requirements
Go 1.20+ (Required to compile the Bubbletea TUI)

macOS 

Git (Required for the sync, push, and clone pipelines)
