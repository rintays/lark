package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"lark/internal/updatecheck"
	"lark/internal/version"
)

const updateCheckInterval = 30 * time.Minute

func autoUpdateEnabled() bool {
	if v := strings.TrimSpace(os.Getenv("LARK_NO_UPDATE_CHECK")); v != "" && v != "0" {
		return false
	}
	if v := strings.TrimSpace(os.Getenv("LARK_UPDATE_CHECK")); v != "" {
		// allow 0/false/off
		vv := strings.ToLower(v)
		if vv == "0" || vv == "false" || vv == "off" {
			return false
		}
	}
	return true
}

func isSystemCommand(command string) bool {
	command = strings.TrimSpace(command)
	if command == "" {
		return true
	}
	root := strings.Fields(command)
	if len(root) == 0 {
		return true
	}
	switch root[0] {
	case "help", "completion", "version", "upgrade", "auth", "config":
		return true
	default:
		return false
	}
}

func handleAutoUpdate(state *appState) {
	if state == nil {
		return
	}
	if !autoUpdateEnabled() {
		return
	}
	if isSystemCommand(state.Command) {
		return
	}

	st, err := updatecheck.Load(state.ConfigPath)
	if err != nil {
		return
	}
	now := time.Now()
	if st.Suppressed(now) {
		// Still refresh in the background.
		goto ASYNC
	}

	// If we have a pending update discovered in the background, handle it
	// before executing the command (minor/major intercept, patch info).
	if !state.JSON {
		switch st.PendingLevel {
		case "patch":
			if st.PendingTo != "" && st.PendingTo != st.LastPromptedVersion {
				fmt.Fprintf(errWriter(state), "INFO: lark update available (%s -> %s). Run: lark upgrade\n", st.PendingFrom, st.PendingTo)
				st.LastPromptedVersion = st.PendingTo
				st.PendingLevel = ""
				_ = updatecheck.Save(state.ConfigPath, st)
			}
		case "minor", "major":
			if st.PendingTo != "" && st.PendingTo != st.LastPromptedVersion {
				m := newAutoUpgradePromptModel(st.PendingFrom, st.PendingTo, st.PendingLevel)
				p := tea.NewProgram(m)
				res, err := p.Run()
				if err == nil {
					if mm, ok := res.(autoUpgradePromptModel); ok {
						if mm.choice == "upgrade" {
							_ = upgradeToLatest(context.Background(), state, upgradeOptions{Yes: true, Interactive: false, Source: "github"})
						} else {
							fmt.Fprintln(errWriter(state), "Skipped upgrade.")
						}
					}
				}
				st.LastPromptedVersion = st.PendingTo
				st.PendingLevel = ""
				// Suppress for a while to avoid repeated prompts.
				st.SuppressUntilUnix = time.Now().Add(24 * time.Hour).Unix()
				_ = updatecheck.Save(state.ConfigPath, st)
			}
		}
	}

ASYNC:
	// Async version check (throttled).
	if !st.ShouldCheck(now, updateCheckInterval) {
		return
	}
	st.LastCheckedAtUnix = now.Unix()
	_ = updatecheck.Save(state.ConfigPath, st)

	go func() {
		latest, found, err := detectLatest(context.Background(), false)
		if err != nil || !found || latest == "" {
			return
		}
		cur := version.Version
		if cur == "" {
			cur = "v0.0.0"
		}
		level, ok := updatecheck.DiffLevel(cur, latest)
		if !ok {
			return
		}
		st2, err := updatecheck.Load(state.ConfigPath)
		if err != nil {
			return
		}
		st2.LatestVersion = latest
		if level == "" {
			st2.PendingLevel = ""
			st2.PendingFrom = ""
			st2.PendingTo = ""
		} else {
			st2.PendingLevel = level
			st2.PendingFrom = cur
			st2.PendingTo = latest
		}
		_ = updatecheck.Save(state.ConfigPath, st2)
	}()
}

// ---- bubbletea prompt ----

type autoUpgradePromptModel struct {
	from   string
	to     string
	level  string
	choice string // upgrade|skip
}

func newAutoUpgradePromptModel(from, to, level string) autoUpgradePromptModel {
	return autoUpgradePromptModel{from: from, to: to, level: level}
}

func (m autoUpgradePromptModel) Init() tea.Cmd { return nil }

func (m autoUpgradePromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "enter":
			m.choice = "upgrade"
			return m, tea.Quit
		case "n", "q", "esc", "ctrl+c":
			m.choice = "skip"
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m autoUpgradePromptModel) View() string {
	head := strings.ToUpper(m.level)
	return fmt.Sprintf("New %s update available: %s -> %s\nUpgrade now? [y/N]\n", head, m.from, m.to)
}
