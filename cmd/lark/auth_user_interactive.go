package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"

	"lark/internal/authregistry"
	"lark/internal/output"
)

type userOAuthSelectionMode int

const (
	userOAuthSelectServices userOAuthSelectionMode = iota
	userOAuthSelectScopes
)

type userOAuthInteractiveSelection struct {
	Mode     userOAuthSelectionMode
	Services []string
	Scopes   []string
}

func promptUserOAuthSelection(state *appState, account string) (userOAuthInteractiveSelection, error) {
	if state == nil {
		return userOAuthInteractiveSelection{}, errors.New("missing app state")
	}
	if state.Printer.JSON || !interactiveAuthAvailable() {
		return userOAuthInteractiveSelection{}, errors.New("interactive login requires a TTY; use --scopes or --services")
	}

	prevServices, prevScopes := previousUserOAuthSelections(state, account)
	defaultMode := userOAuthSelectServices
	if len(prevServices) == 0 && len(prevScopes) > 0 {
		defaultMode = userOAuthSelectScopes
	}

	modeIndex, canceled, err := runModeSelect(defaultMode)
	if err != nil {
		return userOAuthInteractiveSelection{}, err
	}
	if canceled {
		return userOAuthInteractiveSelection{}, errors.New("login canceled")
	}

	if modeIndex == 0 {
		services, err := promptUserOAuthServices(prevServices)
		if err != nil {
			return userOAuthInteractiveSelection{}, err
		}
		return userOAuthInteractiveSelection{Mode: userOAuthSelectServices, Services: services}, nil
	}

	scopes, err := promptUserOAuthScopes(state, prevScopes)
	if err != nil {
		return userOAuthInteractiveSelection{}, err
	}
	return userOAuthInteractiveSelection{Mode: userOAuthSelectScopes, Scopes: scopes}, nil
}

func interactiveAuthAvailable() bool {
	if !output.AutoStyle(os.Stdout) {
		return false
	}
	fd := os.Stdin.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func previousUserOAuthSelections(state *appState, account string) ([]string, []string) {
	if state == nil || state.Config == nil {
		return nil, nil
	}
	acct, ok := loadUserAccount(state.Config, account)
	if !ok {
		return nil, nil
	}
	var services []string
	var scopes []string
	if acct.UserRefreshTokenPayload != nil {
		if len(acct.UserRefreshTokenPayload.Services) > 0 {
			services = normalizeServices(acct.UserRefreshTokenPayload.Services)
		}
		if strings.TrimSpace(acct.UserRefreshTokenPayload.Scopes) != "" {
			scopes = normalizeScopes(parseScopeList(acct.UserRefreshTokenPayload.Scopes))
		}
	}
	if len(scopes) == 0 && len(acct.UserScopes) > 0 {
		scopes = normalizeScopes(acct.UserScopes)
	}
	if len(scopes) == 0 && strings.TrimSpace(acct.UserAccessTokenScope) != "" {
		scopes = normalizeScopes(parseScopeList(acct.UserAccessTokenScope))
	}
	return services, scopes
}

func promptUserOAuthServices(previous []string) ([]string, error) {
	services := authregistry.ListUserOAuthServices()
	if len(services) == 0 {
		return nil, errors.New("no user OAuth services available")
	}

	selectedSet := make(map[string]struct{})
	defaults := previous
	if len(defaults) == 0 {
		defaults = authregistry.DefaultUserOAuthServices
	}
	for _, svc := range normalizeServices(defaults) {
		selectedSet[svc] = struct{}{}
	}

	items := make([]optionItem, 0, len(services))
	for _, svc := range services {
		_, selected := selectedSet[svc]
		items = append(items, optionItem{
			Label:    svc,
			Value:    svc,
			Selected: selected,
		})
	}

	selected, canceled, err := runMultiSelect("Select OAuth services", items, false)
	if err != nil {
		return nil, err
	}
	if canceled {
		return nil, errors.New("login canceled")
	}
	if len(selected) == 0 {
		return nil, errors.New("no services selected")
	}
	return normalizeServices(selected), nil
}

func promptUserOAuthScopes(state *appState, previous []string) ([]string, error) {
	available := userOAuthAvailableScopes()
	defaults := previous
	if len(defaults) == 0 {
		if scopes, _, err := resolveUserOAuthScopes(state, userOAuthScopeOptions{}); err == nil && len(scopes) > 0 {
			defaults = scopes
		} else {
			defaults = []string{defaultUserOAuthScope}
		}
	}

	available = appendMissingScopes(available, defaults)
	selectedSet := make(map[string]struct{})
	for _, scope := range canonicalizeUserOAuthScopes(defaults) {
		selectedSet[scope] = struct{}{}
	}

	items := make([]optionItem, 0, len(available))
	for _, scope := range available {
		_, selected := selectedSet[scope]
		locked := scope == defaultUserOAuthScope
		items = append(items, optionItem{
			Label:    scope,
			Value:    scope,
			Selected: selected || locked,
			Locked:   locked,
			Tag:      "required",
		})
	}

	selected, canceled, err := runMultiSelect("Select OAuth scopes", items, true)
	if err != nil {
		return nil, err
	}
	if canceled {
		return nil, errors.New("login canceled")
	}
	selected = ensureOfflineAccess(selected)
	return canonicalizeUserOAuthScopes(selected), nil
}

func appendMissingScopes(available []string, defaults []string) []string {
	seen := make(map[string]struct{}, len(available))
	for _, scope := range available {
		seen[scope] = struct{}{}
	}
	for _, scope := range normalizeScopes(defaults) {
		if _, ok := seen[scope]; ok {
			continue
		}
		available = append(available, scope)
		seen[scope] = struct{}{}
	}
	return available
}

func runModeSelect(defaultMode userOAuthSelectionMode) (int, bool, error) {
	delegate := brandDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	items := []list.Item{
		modeItem{label: "Select by service (recommended)"},
		modeItem{label: "Select by scope"},
	}

	model := list.New(items, delegate, 0, 0)
	model.Title = "Choose OAuth selection mode"
	model.Styles = brandListStyles()
	model.SetShowHelp(false)
	model.SetFilteringEnabled(false)
	model.SetShowFilter(false)
	model.SetShowStatusBar(false)
	model.DisableQuitKeybindings()
	model.Select(modeIndexFor(defaultMode))

	keys := modeKeyMap{listKeys: list.DefaultKeyMap()}
	m := &singleSelectModel{
		list: model,
		help: brandHelpModel(),
		keys: keys,
	}

	program := tea.NewProgram(m, tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
	result, err := program.Run()
	if err != nil {
		return 0, false, err
	}
	final, ok := result.(*singleSelectModel)
	if !ok {
		return 0, false, errors.New("unexpected selection result")
	}
	return final.list.Index(), final.canceled, nil
}

func runMultiSelect(title string, items []optionItem, allowEmpty bool) ([]string, bool, error) {
	delegate := brandDelegate()
	delegate.ShowDescription = false

	listItems := make([]list.Item, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, item)
	}

	model := list.New(listItems, delegate, 0, 0)
	model.Title = title
	model.Styles = brandListStyles()
	model.SetShowHelp(false)
	model.SetFilteringEnabled(false)
	model.SetShowFilter(false)
	model.DisableQuitKeybindings()

	keys := multiKeyMap{listKeys: list.DefaultKeyMap()}
	m := &multiSelectModel{
		list:       model,
		help:       brandHelpModel(),
		keys:       keys,
		allowEmpty: allowEmpty,
	}

	program := tea.NewProgram(m, tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
	result, err := program.Run()
	if err != nil {
		return nil, false, err
	}
	final, ok := result.(*multiSelectModel)
	if !ok {
		return nil, false, errors.New("unexpected selection result")
	}
	return final.selectedValues(), final.canceled, nil
}

type modeItem struct {
	label string
}

func (i modeItem) Title() string       { return i.label }
func (i modeItem) Description() string { return "" }
func (i modeItem) FilterValue() string { return i.label }

type optionItem struct {
	Label    string
	Value    string
	Selected bool
	Locked   bool
	Tag      string
}

func (i optionItem) Title() string {
	box := "[ ]"
	if i.Selected {
		box = "[x]"
	}
	label := i.Label
	if i.Locked && i.Tag != "" {
		label = fmt.Sprintf("%s (%s)", label, i.Tag)
	}
	return fmt.Sprintf("%s %s", box, label)
}

func (i optionItem) Description() string { return "" }
func (i optionItem) FilterValue() string { return i.Label + " " + i.Value }

type modeKeyMap struct {
	listKeys list.KeyMap
	Confirm  key.Binding
	Quit     key.Binding
}

func (k modeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.listKeys.CursorUp, k.listKeys.CursorDown, k.Confirm, k.Quit}
}

func (k modeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.listKeys.CursorUp, k.listKeys.CursorDown, k.Confirm, k.Quit}}
}

type multiKeyMap struct {
	listKeys list.KeyMap
	Toggle   key.Binding
	Confirm  key.Binding
	Quit     key.Binding
}

func (k multiKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.listKeys.CursorUp, k.listKeys.CursorDown, k.Toggle, k.Confirm, k.Quit}
}

func (k multiKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.listKeys.CursorUp, k.listKeys.CursorDown, k.Toggle, k.Confirm, k.Quit}}
}

type singleSelectModel struct {
	list     list.Model
	help     help.Model
	keys     modeKeyMap
	canceled bool
}

func (m *singleSelectModel) Init() tea.Cmd {
	m.keys.Confirm = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm"))
	m.keys.Quit = key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "cancel"))
	return nil
}

func (m *singleSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, listHeight(msg.Height, 1, 0))
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.canceled = true
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *singleSelectModel) View() string {
	helpView := m.help.View(m.keys)
	if helpView == "" {
		return m.list.View()
	}
	return lipgloss.JoinVertical(lipgloss.Left, m.list.View(), helpView)
}

type multiSelectModel struct {
	list       list.Model
	help       help.Model
	keys       multiKeyMap
	allowEmpty bool
	canceled   bool
	errMessage string
}

func (m *multiSelectModel) Init() tea.Cmd {
	m.keys.Toggle = key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle"))
	m.keys.Confirm = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm"))
	m.keys.Quit = key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "cancel"))
	return nil
}

func (m *multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		extra := 1
		if m.errMessage != "" {
			extra = 2
		}
		m.list.SetSize(msg.Width, listHeight(msg.Height, extra, 0))
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.canceled = true
			return m, tea.Quit
		case " ":
			m.toggleSelected()
			return m, nil
		case "enter":
			if !m.allowEmpty && len(m.selectedValues()) == 0 {
				m.errMessage = "Select at least one item."
				return m, nil
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *multiSelectModel) View() string {
	helpView := m.help.View(m.keys)
	blocks := make([]string, 0, 3)
	if m.errMessage != "" {
		blocks = append(blocks, errorStyle().Render(m.errMessage))
	}
	blocks = append(blocks, m.list.View())
	if helpView != "" {
		blocks = append(blocks, helpView)
	}
	return lipgloss.JoinVertical(lipgloss.Left, blocks...)
}

func (m *multiSelectModel) toggleSelected() {
	index := m.list.Index()
	items := m.list.Items()
	if index < 0 || index >= len(items) {
		return
	}
	current, ok := items[index].(optionItem)
	if !ok {
		return
	}
	if current.Locked {
		return
	}
	current.Selected = !current.Selected
	items[index] = current
	_ = m.list.SetItems(items)
}

func (m *multiSelectModel) selectedValues() []string {
	items := m.list.Items()
	selected := make([]string, 0, len(items))
	for _, item := range items {
		opt, ok := item.(optionItem)
		if !ok {
			continue
		}
		if opt.Selected {
			selected = append(selected, opt.Value)
		}
	}
	return selected
}

func listHeight(total, footerLines, headerLines int) int {
	height := total - footerLines - headerLines
	if height < 4 {
		return 4
	}
	return height
}

func brandDelegate() list.DefaultDelegate {
	brand := output.BrandColor()
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Copy().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(brand).
		Foreground(brand)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Copy().Foreground(brand)
	return delegate
}

func brandListStyles() list.Styles {
	brand := output.BrandColor()
	styles := list.DefaultStyles()
	styles.Title = styles.Title.Copy().
		Background(brand).
		Foreground(lipgloss.Color("255")).
		Bold(true)
	styles.TitleBar = styles.TitleBar.Copy().Padding(0, 0, 1, 1)
	styles.PaginationStyle = styles.PaginationStyle.Copy().Foreground(brand)
	styles.ActivePaginationDot = styles.ActivePaginationDot.Copy().Foreground(brand)
	styles.ArabicPagination = styles.ArabicPagination.Copy().Foreground(brand)
	styles.HelpStyle = styles.HelpStyle.Copy().Foreground(brand)
	styles.StatusBar = styles.StatusBar.Copy().Foreground(brand)
	styles.StatusBarActiveFilter = styles.StatusBarActiveFilter.Copy().Foreground(brand)
	styles.StatusBarFilterCount = styles.StatusBarFilterCount.Copy().Foreground(brand)
	return styles
}

func brandHelpModel() help.Model {
	brand := output.BrandColor()
	m := help.New()
	m.Styles.ShortKey = m.Styles.ShortKey.Copy().Foreground(brand).Bold(true)
	m.Styles.FullKey = m.Styles.FullKey.Copy().Foreground(brand).Bold(true)
	m.Styles.ShortSeparator = m.Styles.ShortSeparator.Copy().Foreground(brand)
	m.Styles.FullSeparator = m.Styles.FullSeparator.Copy().Foreground(brand)
	return m
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#E02E2E")).Bold(true)
}

func modeIndexFor(mode userOAuthSelectionMode) int {
	if mode == userOAuthSelectScopes {
		return 1
	}
	return 0
}
