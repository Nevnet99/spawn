package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	selectScripts step = iota
	reviewPlan
	executingScripts
	showLogs
)

type tabIndex int

const (
	creationTab tabIndex = iota
	destructionTab
)

type focusElement int

const (
	listFocus focusElement = iota
	runButtonFocus
	updateButtonFocus
	pushButtonFocus // New focus target for pushing configs up
)

type scriptItem struct {
	name        string
	description string
	selected    bool
}

type model struct {
	activeStep   step
	activeTab    tabIndex
	activeFocus  focusElement
	scriptList   []scriptItem
	cursor       int
	dryRun       bool
	isUpdateOnly bool
	isPushOnly   bool // New flag to route the push script
	err          error
	done         bool
	outputLog    string
	termWidth    int
}

func initialModel(dryRun bool) model {
	return model{
		activeStep:  selectScripts,
		activeTab:   creationTab,
		activeFocus: listFocus,
		scriptList: []scriptItem{
			{name: "Git Environment Setup", description: "Configures workspace directories and profiles", selected: true},
			{name: "App Payload Provisioning", description: "Installs desktop GUI applications (Chrome, Spotify, Ghostty, Obsidian, Recordly)", selected: true},
			{name: "Terminal Utilities Setup", description: "Installs modern command line binary upgrades (eza, bat, zoxide, ripgrep, fzf, golang, nvm, starship)", selected: true},
			{name: "Shell Profile Configuration", description: "Binds operational shorthand aliases and runtime environment hooks to your shell rc files", selected: true},
			{name: "SSH Configuration Guard", description: "Generates secure Ed25519 identity keys and sets up automatic keychain routing", selected: true},
			{name: "Ghostty Preferences Setup", description: "Configures terminal margins, font rendering dimensions, and window presentation styles", selected: true},
			{name: "NvChad Configuration Link", description: "Deploys a blazing fast, hyper-optimized NvChad layout structure via user profile symlinks", selected: true},
			{name: "macOS Core Velocity Tuning", description: "Overrides operating system delays to enable lightning-fast key repeat rates and directional scrolling adjustments", selected: true},
		},
		cursor:    0,
		dryRun:    dryRun,
		termWidth: 80,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func runSelectedScriptsCmd(m model) tea.Cmd {
	return func() tea.Msg {
		var output strings.Builder

		// Route path 1: Critical Security Purge
		if m.activeTab == destructionTab {
			output.WriteString("▶️ Executing: Scorched Earth Protocol...\n")
			cmd := exec.Command("/bin/bash", "./scripts/scorched_earth.sh")
			cmd.Env = append(os.Environ(), fmt.Sprintf("SPAWN_DRY_RUN=%t", m.dryRun))
			out, err := cmd.CombinedOutput()
			output.Write(out)
			if err != nil {
				return errMsg{err: fmt.Errorf("Scorched Earth failed: %w", err)}
			}
			return successMsg{log: output.String()}
		}

		// Route path 2: Pull Updates Down
		if m.isUpdateOnly {
			output.WriteString("▶️ Executing: Cloud Configuration Sync (Pull)...\n")
			cmd := exec.Command("/bin/bash", "./scripts/update_workspace.sh")
			cmd.Env = append(os.Environ(), fmt.Sprintf("SPAWN_DRY_RUN=%t", m.dryRun))
			out, err := cmd.CombinedOutput()
			output.Write(out)
			if err != nil {
				return successMsg{log: output.String()}
			}
			return successMsg{log: output.String()}
		}

		// Route path 3: Push Updates Up
		if m.isPushOnly {
			output.WriteString("▶️ Executing: Cloud Configuration Backup (Push)...\n")
			cmd := exec.Command("/bin/bash", "./scripts/push_configs.sh")
			cmd.Env = append(os.Environ(), fmt.Sprintf("SPAWN_DRY_RUN=%t", m.dryRun))
			out, err := cmd.CombinedOutput()
			output.Write(out)
			if err != nil {
				return successMsg{log: output.String()}
			}
			return successMsg{log: output.String()}
		}

		// Route path 4: Full Workspace Deployment Checklist Sequence
		for _, script := range m.scriptList {
			if !script.selected {
				continue
			}

			var scriptPath string
			switch script.name {
			case "Git Environment Setup":
				scriptPath = "./scripts/git_setup.sh"
			case "App Payload Provisioning":
				scriptPath = "./scripts/os_packages.sh"
			case "Terminal Utilities Setup":
				scriptPath = "./scripts/cli_tools.sh"
			case "Shell Profile Configuration":
				scriptPath = "./scripts/terminal_setup.sh"
			case "SSH Configuration Guard":
				scriptPath = "./scripts/ssh_setup.sh"
			case "Ghostty Preferences Setup":
				scriptPath = "./scripts/ghostty_setup.sh"
			case "NvChad Configuration Link":
				scriptPath = "./scripts/nvim_setup.sh"
			case "macOS Core Velocity Tuning":
				scriptPath = "./scripts/macos_defaults.sh"
			}

			output.WriteString(fmt.Sprintf("▶️ Executing: %s...\n", script.name))

			cmd := exec.Command("/bin/bash", scriptPath)
			cmd.Env = append(os.Environ(),
				fmt.Sprintf("SPAWN_DRY_RUN=%t", m.dryRun),
				"SPAWN_WORKSPACE=personal",
				"SPAWN_NAME=Luke Brannagan",
				"SPAWN_PERSONAL_EMAIL=luke@personal.dev",
			)

			out, err := cmd.CombinedOutput()
			output.Write(out)
			output.WriteString("\n")

			if err != nil {
				return successMsg{log: output.String()}
			}
		}

		return successMsg{log: output.String()}
	}
}

type errMsg struct{ err error }
type successMsg struct{ log string }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case successMsg:
		m.outputLog = msg.log
		m.activeStep = showLogs
		return m, nil
	case errMsg:
		m.err = msg.err
		m.activeStep = showLogs
		return m, nil
	}

	switch m.activeStep {
	case selectScripts:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "tab":
				if m.activeTab == creationTab {
					m.activeTab = destructionTab
					m.activeFocus = runButtonFocus
				} else {
					m.activeTab = creationTab
					m.activeFocus = listFocus
				}
				m.cursor = 0
				m.isUpdateOnly = false
				m.isPushOnly = false
			case "up", "k":
				if m.activeTab == creationTab {
					if m.activeFocus == runButtonFocus || m.activeFocus == updateButtonFocus || m.activeFocus == pushButtonFocus {
						m.activeFocus = listFocus
						m.cursor = len(m.scriptList) - 1
					} else if m.cursor > 0 {
						m.cursor--
					}
				}
			case "down", "j":
				if m.activeTab == creationTab {
					if m.activeFocus == listFocus {
						if m.cursor < len(m.scriptList)-1 {
							m.cursor++
						} else {
							m.activeFocus = runButtonFocus
						}
					}
				}
			case "left", "h":
				if m.activeTab == creationTab {
					if m.activeFocus == pushButtonFocus {
						m.activeFocus = updateButtonFocus
					} else if m.activeFocus == updateButtonFocus {
						m.activeFocus = runButtonFocus
					}
				}
			case "right", "l":
				if m.activeTab == creationTab {
					if m.activeFocus == runButtonFocus {
						m.activeFocus = updateButtonFocus
					} else if m.activeFocus == updateButtonFocus {
						m.activeFocus = pushButtonFocus
					}
				}
			case "a":
				if m.activeTab == creationTab && m.activeFocus == listFocus {
					anyUnselected := false
					for _, item := range m.scriptList {
						if !item.selected {
							anyUnselected = true
							break
						}
					}
					for i := range m.scriptList {
						m.scriptList[i].selected = anyUnselected
					}
				}
			case " ":
				if m.activeTab == creationTab && m.activeFocus == listFocus {
					m.scriptList[m.cursor].selected = !m.scriptList[m.cursor].selected
				}
			case "enter":
				if m.activeTab == destructionTab {
					m.activeStep = reviewPlan
				} else if m.activeFocus == runButtonFocus {
					m.isUpdateOnly = false
					m.isPushOnly = false
					m.activeStep = reviewPlan
				} else if m.activeFocus == updateButtonFocus {
					m.isUpdateOnly = true
					m.isPushOnly = false
					m.activeStep = reviewPlan
				} else if m.activeFocus == pushButtonFocus {
					m.isUpdateOnly = false
					m.isPushOnly = true
					m.activeStep = reviewPlan
				}
			}
		}

	case reviewPlan:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch strings.ToLower(keyMsg.String()) {
			case "y", "enter":
				m.activeStep = executingScripts
				return m, runSelectedScriptsCmd(m)
			case "n":
				m.activeStep = selectScripts
			}
		}
	}

	return m, nil
}

func renderHeader() string {
	frogGreen := lipgloss.Color("#44B78B")
	banner := "▗▄▄▖▗▄▄▖  ▗▄▖ ▗▖ ▗▖▗▖  ▗▖\n" +
		"▐▌   ▐▌ ▐▌▐▌ ▐▌▐▌ ▐▌▐▛▚▖▐▌\n" +
		" ▝▀▚▖▐▛▀▘ ▐▛▀▜▌▐▌ ▐▌▐▌ ▝▜▌\n" +
		"▗▄▄▞▘▐▌   ▐▌ ▐▌▐▙█▟▌▐▌  ▐▌"
	return lipgloss.NewStyle().Foreground(frogGreen).Bold(true).Render(banner)
}

func (m model) View() string {
	outerChromeWidth := m.termWidth - 2
	if outerChromeWidth < 20 {
		outerChromeWidth = 20
	}
	innerContentWidth := outerChromeWidth - 4

	mintGreen := lipgloss.Color("#44B78B")
	steelBlue := lipgloss.Color("#4A90E2")
	accentOrange := lipgloss.Color("#F2A65A")
	dangerRed := lipgloss.Color("#E05C5C")
	dimGrey := lipgloss.Color("#626262")
	darkMuted := lipgloss.Color("#444444")
	appBorderColor := lipgloss.Color("#2B2B2B")

	if m.err != nil {
		errorBoxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(dangerRed).Padding(1, 2).Width(innerContentWidth - 2)
		var errS string
		errS += lipgloss.NewStyle().Foreground(dangerRed).Bold(true).Render("❌ SEQUENCE EXECUTION FAILURE") + "\n\n"
		errS += fmt.Sprintf("%v\n\n", m.err)
		errS += lipgloss.NewStyle().Foreground(dimGrey).Render("Review the output log history context above:") + "\n"
		errS += fmt.Sprintf("%s\n\n", m.outputLog)
		errS += lipgloss.NewStyle().Foreground(dimGrey).Render("Press 'q' or 'ctrl+c' to close.")

		var finalView strings.Builder
		finalView.WriteString(renderHeader() + "\n")
		appWrapperWindowStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(appBorderColor).Padding(1, 2).Width(outerChromeWidth)
		finalView.WriteString(appWrapperWindowStyle.Render(errorBoxStyle.Render(errS)) + "\n")
		return finalView.String()
	}

	if m.activeStep == showLogs {
		logsBoxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(mintGreen).Padding(1, 2).Width(innerContentWidth - 2)
		var logsS string
		logsS += lipgloss.NewStyle().Foreground(mintGreen).Bold(true).Render("📋 SEQUENCE LOG EXECUTION OUTPUT HISTORY") + "\n"

		dividerLine := lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C")).Render(strings.Repeat("─", innerContentWidth-4))
		logsS += dividerLine + "\n"
		logsS += m.outputLog + "\n"
		logsS += lipgloss.NewStyle().Foreground(dimGrey).Render("Press 'q' or 'ctrl+c' to close this log shell and exit spawn.")

		var finalView strings.Builder
		finalView.WriteString(renderHeader() + "\n")
		appWrapperWindowStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(appBorderColor).Padding(1, 2).Width(outerChromeWidth)
		finalView.WriteString(appWrapperWindowStyle.Render(logsBoxStyle.Render(logsS)) + "\n")
		return finalView.String()
	}

	activeStyle := lipgloss.NewStyle().Foreground(mintGreen).Bold(true)
	inactiveStyle := lipgloss.NewStyle().Foreground(dimGrey)
	descStyle := lipgloss.NewStyle().Foreground(darkMuted)

	var bodyContent string

	if m.dryRun {
		bannerText := "⚠️  DRY-RUN MODE ACTIVE"
		paddingLength := innerContentWidth - len(bannerText)
		if paddingLength < 0 {
			paddingLength = 0
		}
		bodyContent += lipgloss.NewStyle().
			Background(lipgloss.Color("#F2A65A")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1).
			Render(bannerText+strings.Repeat(" ", paddingLength)) + "\n\n"
	}

	switch m.activeStep {
	case selectScripts:
		var tabs []string
		activeTabStyle := lipgloss.NewStyle().Background(steelBlue).Foreground(lipgloss.Color("#000000")).Bold(true).Padding(0, 2)
		inactiveTabStyle := lipgloss.NewStyle().Background(darkMuted).Foreground(dimGrey).Padding(0, 2)

		if m.activeTab == creationTab {
			tabs = append(tabs, activeTabStyle.Render("🛠  CREATION"))
			tabs = append(tabs, inactiveTabStyle.Render("🔥 DESTRUCTION"))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render("🛠  CREATION"))
			tabs = append(tabs, activeTabStyle.Render("🔥 DESTRUCTION"))
		}
		bodyContent += lipgloss.JoinHorizontal(lipgloss.Top, tabs...) + "\n\n"

		panelStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(appBorderColor).Padding(1, 2).Width(innerContentWidth - 2)
		var panelContent strings.Builder

		if m.activeTab == creationTab {
			for i, item := range m.scriptList {
				isCursor := m.cursor == i && m.activeFocus == listFocus
				checked := " "
				if item.selected {
					checked = "✓"
				}
				checkbox := fmt.Sprintf("[%s]", checked)
				pointer := "  "
				if isCursor {
					pointer = "👉"
				}

				var lineTitle string
				if isCursor {
					lineTitle = fmt.Sprintf("%s %s %s", pointer, activeStyle.Render(checkbox), activeStyle.Render(item.name))
				} else if item.selected {
					lineTitle = fmt.Sprintf("%s %s %s", pointer, activeStyle.Render(checkbox), lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC")).Render(item.name))
				} else {
					lineTitle = fmt.Sprintf("%s %s %s", pointer, inactiveStyle.Render(checkbox), inactiveStyle.Render(item.name))
				}
				panelContent.WriteString(lineTitle + "\n" + fmt.Sprintf("   %s\n", descStyle.Render(item.description)))
			}
			panelContent.WriteString("\n")

			// Dynamic Three-Button Layout Rendering
			runBtnStyle := lipgloss.NewStyle().Padding(0, 2).Bold(true)
			updateBtnStyle := lipgloss.NewStyle().Padding(0, 2).Bold(true)
			pushBtnStyle := lipgloss.NewStyle().Padding(0, 2).Bold(true)

			var rPointer, uPointer, pPointer string

			if m.activeFocus == runButtonFocus {
				runBtnStyle = runBtnStyle.Background(mintGreen).Foreground(lipgloss.Color("#000000"))
				rPointer = "👉"
				uPointer = "  "
				pPointer = "  "
			} else if m.activeFocus == updateButtonFocus {
				updateBtnStyle = updateBtnStyle.Background(steelBlue).Foreground(lipgloss.Color("#000000"))
				rPointer = "  "
				uPointer = "👉"
				pPointer = "  "
			} else if m.activeFocus == pushButtonFocus {
				pushBtnStyle = pushBtnStyle.Background(accentOrange).Foreground(lipgloss.Color("#000000"))
				rPointer = "  "
				uPointer = "  "
				pPointer = "👉"
			} else {
				rPointer = "  "
				uPointer = "  "
				pPointer = "  "
			}

			if m.activeFocus != runButtonFocus {
				runBtnStyle = runBtnStyle.Background(darkMuted).Foreground(lipgloss.Color("#CCCCCC"))
			}
			if m.activeFocus != updateButtonFocus {
				updateBtnStyle = updateBtnStyle.Background(darkMuted).Foreground(lipgloss.Color("#CCCCCC"))
			}
			if m.activeFocus != pushButtonFocus {
				pushBtnStyle = pushBtnStyle.Background(darkMuted).Foreground(lipgloss.Color("#CCCCCC"))
			}

			panelContent.WriteString(fmt.Sprintf("%s %s   %s %s   %s %s\n",
				rPointer, runBtnStyle.Render("🚀 RUN SEQUENCE"),
				uPointer, updateBtnStyle.Render("🔄 PULL UPDATES"),
				pPointer, pushBtnStyle.Render("📤 PUSH CONFIGS")))

			bodyContent += panelStyle.BorderForeground(steelBlue).Render(strings.TrimSuffix(panelContent.String(), "\n")) + "\n\n"
			bodyContent += lipgloss.NewStyle().Foreground(dimGrey).Render("(Arrows/hjkl to navigate layout cells, Space to toggle checks, Tab to switch paths, Enter to execute)")
		} else {
			destructBtnStyle := lipgloss.NewStyle().Padding(0, 3).Bold(true)
			if m.activeFocus == runButtonFocus {
				destructBtnStyle = destructBtnStyle.Background(dangerRed).Foreground(lipgloss.Color("#000000"))
			} else {
				destructBtnStyle = destructBtnStyle.Background(darkMuted).Foreground(dangerRed)
			}

			warningText := lipgloss.NewStyle().Foreground(dangerRed).Bold(true).Render("⚠️  CRITICAL WARNING") + "\n" +
				lipgloss.NewStyle().Foreground(dimGrey).Render("Executing this procedure runs a full clear-slate reset on your configuration data roots.") + "\n\n"

			panelContent.WriteString(warningText)
			panelContent.WriteString(fmt.Sprintf("👉 %s\n", destructBtnStyle.Render("💥 TRIGGER SCORCHED EARTH")))

			bodyContent += panelStyle.BorderForeground(dangerRed).Render(strings.TrimSuffix(panelContent.String(), "\n")) + "\n\n"
			bodyContent += lipgloss.NewStyle().Foreground(dimGrey).Render("(Tab to switch back to creation workspace, Enter to confirm button action)")
		}

	case reviewPlan:
		boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(mintGreen).Padding(1, 2).Width(innerContentWidth - 2).MarginBottom(1)
		planText := lipgloss.NewStyle().Foreground(mintGreen).Bold(true).Render("📋 Scheduled Pipeline Actions") + "\n\n"

		if m.activeTab == destructionTab {
			planText += fmt.Sprintf(" 💥 %s\n", lipgloss.NewStyle().Foreground(dangerRed).Bold(true).Render("Scorched Earth Protocol (Full-Purge Slate Cleanup)"))
		} else if m.isUpdateOnly {
			planText += fmt.Sprintf(" 🔄 %s\n", lipgloss.NewStyle().Foreground(steelBlue).Bold(true).Render("Non-Destructive Configuration Sync (Pull changes from cloud)"))
		} else if m.isPushOnly {
			planText += fmt.Sprintf(" 📤 %s\n", lipgloss.NewStyle().Foreground(accentOrange).Bold(true).Render("Cloud Configuration Backup (Commit and push local changes)"))
		} else {
			hasSelection := false
			for _, item := range m.scriptList {
				if item.selected {
					planText += fmt.Sprintf(" ⚡️ %s\n", item.name)
					hasSelection = true
				}
			}
			if !hasSelection {
				planText += " 🚫 No components selected. Application will exit.\n"
			}
		}

		bodyContent += boxStyle.Render(planText) + "\n"
		prompt := "👉 Execute pipeline? (y/n): "
		if m.dryRun {
			prompt = "👉 Evaluate simulated execution pipeline printouts? (y/n): "
		}
		bodyContent += lipgloss.NewStyle().Bold(true).Render(prompt)

	case executingScripts:
		bodyContent += "⚙️  Executing scheduled architecture scripts... Please wait.\n"
	}

	appWrapperWindowStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(appBorderColor).Padding(1, 2).Width(outerChromeWidth)

	var finalView strings.Builder
	finalView.WriteString(renderHeader() + "\n")
	finalView.WriteString(appWrapperWindowStyle.Render(strings.TrimSuffix(bodyContent, "\n")) + "\n")

	return finalView.String()
}

func main() {
	dryRunPtr := flag.Bool("dry-run", false, "Simulate execution context safely")
	flag.BoolVar(dryRunPtr, "d", false, "Simulate execution context safely")
	flag.Parse()

	p := tea.NewProgram(initialModel(*dryRunPtr))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running spawn: %v", err)
		os.Exit(1)
	}
}
