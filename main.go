package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	selectScripts step = iota
	gatherIdentity
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
	pushButtonFocus
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
	isPushOnly   bool
	err          error
	outputLog    string
	termWidth    int

	// Identity Form Variables
	nameInput      textinput.Model
	emailInput     textinput.Model
	workEmailInput textinput.Model
	inputFocus     int // 0 = Name, 1 = Email, 2 = Work Email
	userName       string
	userEmail      string
	userWorkEmail  string
}

func initialModel(dryRun bool) model {
	nameIn := textinput.New()
	nameIn.Placeholder = "e.g. Luke Brannagan"
	nameIn.Focus()
	nameIn.CharLimit = 50

	emailIn := textinput.New()
	emailIn.Placeholder = "e.g. luke@personal.dev"
	emailIn.CharLimit = 100

	workEmailIn := textinput.New()
	workEmailIn.Placeholder = "e.g. luke.brannagan@company.com (Leave blank for personal only)"
	workEmailIn.CharLimit = 100

	envName := os.Getenv("USER_FULL_NAME")
	envEmail := os.Getenv("USER_GIT_EMAIL")
	envWorkEmail := os.Getenv("USER_WORK_EMAIL")

	return model{
		activeStep:     selectScripts,
		activeTab:      creationTab,
		activeFocus:    listFocus,
		nameInput:      nameIn,
		emailInput:     emailIn,
		workEmailInput: workEmailIn,
		userName:       envName,
		userEmail:      envEmail,
		userWorkEmail:  envWorkEmail,
		inputFocus:     0,
		scriptList: []scriptItem{
			{
				name:        "Git Environment Setup",
				description: "Configures workspace directories and profiles",
				selected:    true,
			},
			{
				name: "App Payload Provisioning",
				description: "Installs desktop GUI applications " +
					"(Chrome, Spotify, Ghostty, Obsidian, Recordly)",
				selected: true,
			},
			{
				name: "Terminal Utilities Setup",
				description: "Installs modern command line binary upgrades " +
					"(eza, bat, zoxide, ripgrep, fzf, golang, nvm, starship)",
				selected: true,
			},
			{
				name: "Shell Profile Configuration",
				description: "Binds operational shorthand aliases and runtime environment " +
					"hooks to your shell rc files",
				selected: true,
			},
			{
				name: "SSH Configuration Guard",
				description: "Generates secure Ed25519 identity keys and sets up " +
					"automatic keychain routing",
				selected: true,
			},
			{
				name: "Ghostty Preferences Setup",
				description: "Configures terminal margins, font rendering dimensions, " +
					"and window presentation styles",
				selected: true,
			},
			{
				name: "NvChad Configuration Link",
				description: "Deploys a blazing fast, hyper-optimized NvChad layout " +
					"structure via user profile symlinks",
				selected: true,
			},
			{
				name: "macOS Core Velocity Tuning",
				description: "Overrides operating system delays to enable lightning-fast " +
					"key repeat rates and directional scrolling adjustments",
				selected: true,
			},
		},
		cursor:    0,
		dryRun:    dryRun,
		termWidth: 80,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func runSelectedScriptsCmd(m model) tea.Cmd {
	var bashCommand string

	if m.activeTab == destructionTab {
		bashCommand = "/bin/bash ./scripts/scorched_earth.sh"
	} else if m.isUpdateOnly {
		bashCommand = "/bin/bash ./scripts/update_workspace.sh"
	} else if m.isPushOnly {
		bashCommand = "/bin/bash ./scripts/push_configs.sh"
	} else {
		var cmds []string
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

			// Add an echo so you can see which script is starting in the terminal output
			cmds = append(cmds, fmt.Sprintf("echo '\n▶️ Executing: %s...'", script.name))
			cmds = append(cmds, fmt.Sprintf("/bin/bash %s", scriptPath))
		}
		if len(cmds) == 0 {
			return func() tea.Msg { return successMsg{log: "No scripts selected."} }
		}
		// String the scripts together so they run in exact sequence natively
		bashCommand = strings.Join(cmds, " && ")
	}

	c := exec.Command("/bin/bash", "-c", bashCommand)

	workspaceMode := "personal"
	if m.userWorkEmail != "" {
		workspaceMode = "work"
	}

	c.Env = append(os.Environ(),
		fmt.Sprintf("SPAWN_DRY_RUN=%t", m.dryRun),
		fmt.Sprintf("SPAWN_WORKSPACE=%s", workspaceMode),
		fmt.Sprintf("SPAWN_NAME=%s", m.userName),
		fmt.Sprintf("SPAWN_PERSONAL_EMAIL=%s", m.userEmail),
		fmt.Sprintf("SPAWN_WORK_EMAIL=%s", m.userWorkEmail),
		"HOMEBREW_NO_AUTO_UPDATE=1",
		"HOMEBREW_NO_INSTALL_CLEANUP=1",
		"NONINTERACTIVE=1",
	)

	// tea.ExecProcess pauses the UI, streams the command natively, and resumes the UI on finish
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return errMsg{err: fmt.Errorf("Pipeline execution failed. Scroll up to see the exact bash error.\n\n%w", err)}
		}
		return successMsg{log: "✅ Pipeline completed successfully!\n\nAll real-time execution logs were printed directly to your terminal. Scroll up to review them."}
	})
}

type errMsg struct{ err error }
type successMsg struct{ log string }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "q" && m.activeStep != gatherIdentity {
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
	case gatherIdentity:
		var cmd tea.Cmd
		var cmds []tea.Cmd

		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				m.activeStep = selectScripts
				return m, nil
			case "tab", "down":
				m.inputFocus = (m.inputFocus + 1) % 3

				m.nameInput.Blur()
				m.emailInput.Blur()
				m.workEmailInput.Blur()

				if m.inputFocus == 0 {
					cmd = m.nameInput.Focus()
				} else if m.inputFocus == 1 {
					cmd = m.emailInput.Focus()
				} else {
					cmd = m.workEmailInput.Focus()
				}
				cmds = append(cmds, cmd)
			case "shift+tab", "up":
				m.inputFocus--
				if m.inputFocus < 0 {
					m.inputFocus = 2
				}

				m.nameInput.Blur()
				m.emailInput.Blur()
				m.workEmailInput.Blur()

				if m.inputFocus == 0 {
					cmd = m.nameInput.Focus()
				} else if m.inputFocus == 1 {
					cmd = m.emailInput.Focus()
				} else {
					cmd = m.workEmailInput.Focus()
				}
				cmds = append(cmds, cmd)
			case "enter":
				// If not on the last field, drop down a line
				if m.inputFocus < 2 {
					m.inputFocus++
					m.nameInput.Blur()
					m.emailInput.Blur()
					if m.inputFocus == 1 {
						cmd = m.emailInput.Focus()
					} else {
						cmd = m.workEmailInput.Focus()
					}
					cmds = append(cmds, cmd)
				} else {
					// Require Name and Personal Email. Work email can remain blank.
					if m.nameInput.Value() != "" && m.emailInput.Value() != "" {
						m.userName = m.nameInput.Value()
						m.userEmail = m.emailInput.Value()
						m.userWorkEmail = m.workEmailInput.Value()
						m.activeStep = reviewPlan
						return m, nil
					}
				}
			}
		}

		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)
		m.emailInput, cmd = m.emailInput.Update(msg)
		cmds = append(cmds, cmd)
		m.workEmailInput, cmd = m.workEmailInput.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

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

					// Smart Trigger: Only ask for identity if Git or SSH setup is actively checked
					needsIdentity := false
					for _, item := range m.scriptList {
						if item.selected && (item.name == "Git Environment Setup" || item.name == "SSH Configuration Guard") {
							needsIdentity = true
							break
						}
					}

					if needsIdentity && (m.userName == "" || m.userEmail == "") {
						m.activeStep = gatherIdentity
					} else {
						m.activeStep = reviewPlan
					}
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
			case "n", "esc":
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
		errorBoxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(dangerRed).Padding(1, 2).Width(innerContentWidth - 2)
		var errS string
		errS += lipgloss.NewStyle().Foreground(dangerRed).Bold(true).
			Render("❌ SEQUENCE EXECUTION FAILURE") + "\n\n"
		errS += fmt.Sprintf("%v\n\n", m.err)
		errS += lipgloss.NewStyle().Foreground(dimGrey).
			Render("Review the output log history context above:") + "\n"
		errS += fmt.Sprintf("%s\n\n", m.outputLog)
		errS += lipgloss.NewStyle().Foreground(dimGrey).
			Render("Press 'q' or 'ctrl+c' to close.")

		var finalView strings.Builder
		finalView.WriteString(renderHeader() + "\n")
		appWrapperWindowStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(appBorderColor).Padding(1, 2).Width(outerChromeWidth)
		finalView.WriteString(appWrapperWindowStyle.Render(errorBoxStyle.Render(errS)) + "\n")
		return finalView.String()
	}

	if m.activeStep == showLogs {
		logsBoxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(mintGreen).Padding(1, 2).Width(innerContentWidth - 2)
		var logsS string
		logsS += lipgloss.NewStyle().Foreground(mintGreen).Bold(true).
			Render("📋 SEQUENCE LOG EXECUTION OUTPUT HISTORY") + "\n"

		dividerLine := lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C")).
			Render(strings.Repeat("─", innerContentWidth-4))
		logsS += dividerLine + "\n"
		logsS += m.outputLog + "\n"
		logsS += lipgloss.NewStyle().Foreground(dimGrey).
			Render("Press 'q' or 'ctrl+c' to close this log shell and exit spawn.")

		var finalView strings.Builder
		finalView.WriteString(renderHeader() + "\n")
		appWrapperWindowStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(appBorderColor).Padding(1, 2).Width(outerChromeWidth)
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
	case gatherIdentity:
		boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(mintGreen).Padding(1, 2).Width(innerContentWidth - 2)
		var form string
		form += lipgloss.NewStyle().Foreground(mintGreen).Bold(true).
			Render("👤 SECURE IDENTITY CONFIGURATION") + "\n"
		form += lipgloss.NewStyle().Foreground(dimGrey).
			Render("Spawn requires your identity to provision Git and SSH keys natively.") + "\n\n"

		form += lipgloss.NewStyle().Bold(true).Render("Full Name:") + "\n"
		form += m.nameInput.View() + "\n\n"

		form += lipgloss.NewStyle().Bold(true).Render("Personal Email:") + "\n"
		form += m.emailInput.View() + "\n\n"

		form += lipgloss.NewStyle().Bold(true).Render("Work Email (Optional):") + "\n"
		form += m.workEmailInput.View() + "\n\n"

		form += lipgloss.NewStyle().Foreground(darkMuted).
			Render("(Tab/Arrows to navigate, Enter to submit, Esc to go back)")
		bodyContent += boxStyle.Render(form)

	case selectScripts:
		var tabs []string
		activeTabStyle := lipgloss.NewStyle().Background(steelBlue).
			Foreground(lipgloss.Color("#000000")).Bold(true).Padding(0, 2)
		inactiveTabStyle := lipgloss.NewStyle().Background(darkMuted).
			Foreground(dimGrey).Padding(0, 2)

		if m.activeTab == creationTab {
			tabs = append(tabs, activeTabStyle.Render("🛠  CREATION"))
			tabs = append(tabs, inactiveTabStyle.Render("🔥 DESTRUCTION"))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render("🛠  CREATION"))
			tabs = append(tabs, activeTabStyle.Render("🔥 DESTRUCTION"))
		}
		bodyContent += lipgloss.JoinHorizontal(lipgloss.Top, tabs...) + "\n\n"

		panelStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(appBorderColor).Padding(1, 2).Width(innerContentWidth - 2)
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
					lineTitle = fmt.Sprintf("%s %s %s", pointer,
						activeStyle.Render(checkbox),
						activeStyle.Render(item.name))
				} else if item.selected {
					lineTitle = fmt.Sprintf("%s %s %s", pointer,
						activeStyle.Render(checkbox),
						lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC")).Render(item.name))
				} else {
					lineTitle = fmt.Sprintf("%s %s %s", pointer,
						inactiveStyle.Render(checkbox),
						inactiveStyle.Render(item.name))
				}
				panelContent.WriteString(lineTitle + "\n" +
					fmt.Sprintf("   %s\n", descStyle.Render(item.description)))
			}
			panelContent.WriteString("\n")

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

			bodyContent += panelStyle.BorderForeground(steelBlue).
				Render(strings.TrimSuffix(panelContent.String(), "\n")) + "\n\n"

			navHint := "(Arrows/hjkl to navigate layout cells, Space to toggle checks, " +
				"Tab to switch paths, Enter to execute)"
			bodyContent += lipgloss.NewStyle().Foreground(dimGrey).Render(navHint)

		} else {
			destructBtnStyle := lipgloss.NewStyle().Padding(0, 3).Bold(true)
			if m.activeFocus == runButtonFocus {
				destructBtnStyle = destructBtnStyle.Background(dangerRed).Foreground(lipgloss.Color("#000000"))
			} else {
				destructBtnStyle = destructBtnStyle.Background(darkMuted).Foreground(dangerRed)
			}

			warningTitle := lipgloss.NewStyle().Foreground(dangerRed).Bold(true).Render("⚠️  CRITICAL WARNING")
			warningDesc := lipgloss.NewStyle().Foreground(dimGrey).
				Render("Executing this procedure runs a full clear-slate reset on your configuration data roots.")

			panelContent.WriteString(warningTitle + "\n" + warningDesc + "\n\n")
			panelContent.WriteString(fmt.Sprintf("👉 %s\n", destructBtnStyle.Render("💥 TRIGGER SCORCHED EARTH")))

			bodyContent += panelStyle.BorderForeground(dangerRed).
				Render(strings.TrimSuffix(panelContent.String(), "\n")) + "\n\n"

			navHint := "(Tab to switch back to creation workspace, Enter to confirm button action)"
			bodyContent += lipgloss.NewStyle().Foreground(dimGrey).Render(navHint)
		}

	case reviewPlan:
		boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(mintGreen).Padding(1, 2).Width(innerContentWidth - 2).MarginBottom(1)
		planText := lipgloss.NewStyle().Foreground(mintGreen).Bold(true).
			Render("📋 Scheduled Pipeline Actions") + "\n\n"

		if m.activeTab == destructionTab {
			planText += fmt.Sprintf(" 💥 %s\n", lipgloss.NewStyle().Foreground(dangerRed).Bold(true).
				Render("Scorched Earth Protocol (Full-Purge Slate Cleanup)"))
		} else if m.isUpdateOnly {
			planText += fmt.Sprintf(" 🔄 %s\n", lipgloss.NewStyle().Foreground(steelBlue).Bold(true).
				Render("Non-Destructive Configuration Sync (Pull changes from cloud)"))
		} else if m.isPushOnly {
			planText += fmt.Sprintf(" 📤 %s\n", lipgloss.NewStyle().Foreground(accentOrange).Bold(true).
				Render("Cloud Configuration Backup (Commit and push local changes)"))
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
		prompt := "👉 Execute pipeline? (y/n/esc): "
		if m.dryRun {
			prompt = "👉 Evaluate simulated execution pipeline printouts? (y/n/esc): "
		}
		bodyContent += lipgloss.NewStyle().Bold(true).Render(prompt)

	case executingScripts:
		bodyContent += "⚙️  Executing scheduled architecture scripts... Please wait.\n"
	}

	appWrapperWindowStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
		BorderForeground(appBorderColor).Padding(1, 2).Width(outerChromeWidth)

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
