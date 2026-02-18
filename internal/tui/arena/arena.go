package arena

import (
	"fmt"
	"strings"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
	"github.com/charmbracelet/lipgloss"
)

// Card dimensions
const (
	cardTotalWidth   = 26 // total rendered width including borders
	cardInnerWidth   = 24 // between │ borders
	cardContentWidth = 22 // innerWidth - 2 padding spaces
)

// palette.md colors
var roleColors = map[schema.Role]string{
	schema.RolePlanner:   "#7AA2F7",
	schema.RoleExecutor:  "#9ECE6A",
	schema.RoleReviewer:  "#BB9AF7",
	schema.RoleGuard:     "#F7768E",
	schema.RoleTester:    "#E0AF68",
	schema.RoleWriter:    "#73DACA",
	schema.RoleExplorer:  "#58A6FF",
	schema.RoleArchitect: "#FFA657",
	schema.RoleDebugger:  "#FF7B72",
	schema.RoleVerifier:  "#56D364",
	schema.RoleDesigner:  "#D2A8FF",
	schema.RoleCustom:    "#8A93A5",
}

var stateColors = map[schema.AgentState]string{
	schema.StateRunning:   "#7EE787",
	schema.StateWaiting:   "#7D8590",
	schema.StateBlocked:   "#E3B341",
	schema.StateError:     "#FF7B72",
	schema.StateDone:      "#56D364",
	schema.StateFailed:    "#FF7B72",
	schema.StateCancelled: "#6E7681",
	schema.StateIdle:      "#6E7681",
}

// Model represents the Arena panel state.
type Model struct {
	agents   map[string]*AgentCard
	order    []string // ordered agent IDs for stable rendering
	width    int
	height   int
	selected int  // cursor position for navigation
	focused  bool // whether this panel has focus
	unicode  bool // whether to use unicode sprites
}

// AgentCard holds display state for a single agent.
type AgentCard struct {
	AgentID string
	Role    schema.Role
	State   schema.AgentState
	Summary string // latest event summary
}

// NewModel creates a new Arena model.
func NewModel() Model {
	return Model{
		agents:  make(map[string]*AgentCard),
		order:   make([]string, 0),
		unicode: true,
	}
}

// UpdateAgent adds or updates an agent card.
func (m *Model) UpdateAgent(agentID string, role schema.Role, state schema.AgentState) {
	if _, exists := m.agents[agentID]; !exists {
		m.order = append(m.order, agentID)
	}
	m.agents[agentID] = &AgentCard{
		AgentID: agentID,
		Role:    role,
		State:   state,
	}
}

// UpdateAgentWithSummary adds or updates an agent card with event summary.
func (m *Model) UpdateAgentWithSummary(agentID string, role schema.Role, state schema.AgentState, summary string) {
	if _, exists := m.agents[agentID]; !exists {
		m.order = append(m.order, agentID)
	}
	m.agents[agentID] = &AgentCard{
		AgentID: agentID,
		Role:    role,
		State:   state,
		Summary: summary,
	}
}

// SetSize updates the panel dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetFocused sets whether this panel has keyboard focus.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

// SetUnicode toggles unicode/ASCII rendering mode.
func (m *Model) SetUnicode(unicode bool) {
	m.unicode = unicode
}

// SelectedAgent returns the currently selected agent card, or nil.
func (m Model) SelectedAgent() *AgentCard {
	if len(m.order) == 0 {
		return nil
	}
	if m.selected < 0 || m.selected >= len(m.order) {
		return nil
	}
	id := m.order[m.selected]
	return m.agents[id]
}

// AgentCount returns the number of agents.
func (m Model) AgentCount() int {
	return len(m.order)
}

// HandleKey processes keyboard input for arena navigation.
// Returns true if the key was consumed.
func (m *Model) HandleKey(key string) bool {
	switch key {
	case "j", "down":
		if len(m.order) > 0 {
			m.selected = (m.selected + 1) % len(m.order)
		}
		return true
	case "k", "up":
		if len(m.order) > 0 {
			m.selected = (m.selected - 1 + len(m.order)) % len(m.order)
		}
		return true
	case "h", "left":
		if len(m.order) > 0 {
			m.selected = (m.selected - 1 + len(m.order)) % len(m.order)
		}
		return true
	case "l", "right":
		if len(m.order) > 0 {
			m.selected = (m.selected + 1) % len(m.order)
		}
		return true
	}
	return false
}

// View renders the Arena panel.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	borderColor := "#2A3142"
	if m.focused {
		borderColor = "#58A6FF"
	}

	if len(m.order) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Width(m.width - 2).
			Height(m.height - 2).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("#8A93A5")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(borderColor))
		return emptyStyle.Render("No agents")
	}

	var cards []string
	for i, id := range m.order {
		agent := m.agents[id]
		isSelected := m.focused && i == m.selected
		cards = append(cards, renderCard(agent, isSelected, m.unicode))
	}

	// Join cards horizontally if space allows
	maxHorizontal := m.width / cardTotalWidth
	if maxHorizontal < 1 {
		maxHorizontal = 1
	}

	var rows []string
	for i := 0; i < len(cards); i += maxHorizontal {
		end := i + maxHorizontal
		if end > len(cards) {
			end = len(cards)
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, cards[i:end]...)
		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E6EAF2")).
		Bold(true)
	title := titleStyle.Render(" Agent Arena ")

	style := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		BorderTop(true).
		Padding(0, 1)

	return style.Render(title + "\n" + content)
}

// boxTop returns the top border of a card box.
func boxTop(borderStyle lipgloss.Style) string {
	return borderStyle.Render("╭" + strings.Repeat("─", cardInnerWidth) + "╮")
}

// boxBottom returns the bottom border of a card box.
func boxBottom(borderStyle lipgloss.Style) string {
	return borderStyle.Render("╰" + strings.Repeat("─", cardInnerWidth) + "╯")
}

// boxLine renders a content line wrapped in │ ... │ with right-padding.
// styledContent is already styled (may contain ANSI codes).
// The visual width is measured with lipgloss.Width and padded to cardContentWidth.
func boxLine(styledContent string, borderStyle lipgloss.Style) string {
	visWidth := lipgloss.Width(styledContent)
	pad := cardContentWidth - visWidth
	if pad < 0 {
		pad = 0
	}
	return borderStyle.Render("│") + " " + styledContent + strings.Repeat(" ", pad) + " " + borderStyle.Render("│")
}

// renderCard renders a single agent card with CLCO sprite in fixed-width box.
func renderCard(card *AgentCard, selected bool, useUnicode bool) string {
	roleColor := getRoleColor(card.Role)
	stateColor := getStateColor(card.State)
	indicator := GetStateIndicator(card.State, useUnicode)

	// Border color based on state
	borderColor := "#2A3142"
	if selected {
		borderColor = "#58A6FF"
	}
	if card.State == schema.StateError || card.State == schema.StateFailed {
		borderColor = "#FF7B72"
	}
	if card.State == schema.StateBlocked {
		borderColor = "#E3B341"
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))

	// Sprite tint: role color by default, override for error/done states
	spriteColor := roleColor
	if card.State == schema.StateError || card.State == schema.StateFailed {
		spriteColor = "#FF7B72"
	}
	if card.State == schema.StateDone {
		spriteColor = "#56D364"
	}

	spriteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(spriteColor))
	if card.State == schema.StateWaiting || card.State == schema.StateIdle {
		spriteStyle = spriteStyle.Faint(true)
	}

	roleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(roleColor)).
		Bold(true)

	stateStyle := getStateStyle(card.State)

	idStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8A93A5"))

	indicatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(stateColor)).
		Bold(card.State == schema.StateRunning)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8A93A5")).
		Bold(true)

	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8A93A5")).
		Italic(true)

	// Build card lines
	var lines []string

	// Top border
	lines = append(lines, boxTop(borderStyle))

	// 3-line CLCO sprite, block-centered with uniform x-anchor
	sprite := GetSprite(useUnicode)
	maxW := 0
	for _, line := range sprite {
		if w := len([]rune(line)); w > maxW {
			maxW = w
		}
	}
	leftPad := (cardContentWidth - maxW) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	prefix := strings.Repeat(" ", leftPad)
	for _, spriteLine := range sprite {
		runes := []rune(spriteLine)
		rightPad := maxW - len(runes)
		if rightPad < 0 {
			rightPad = 0
		}
		aligned := prefix + string(runes) + strings.Repeat(" ", rightPad)
		lines = append(lines, boxLine(spriteStyle.Render(aligned), borderStyle))
	}

	// Role + indicator
	roleStr := fmt.Sprintf("%s %s",
		roleStyle.Render(string(card.Role)),
		indicatorStyle.Render(indicator),
	)
	lines = append(lines, boxLine(roleStr, borderStyle))

	// Agent ID
	lines = append(lines, boxLine(idStyle.Render(truncate(card.AgentID, cardContentWidth)), borderStyle))

	// State
	lines = append(lines, boxLine(stateStyle.Render(string(card.State)), borderStyle))

	// Recent activity label
	lines = append(lines, boxLine(labelStyle.Render("Recent activity"), borderStyle))

	// Summary
	summary := card.Summary
	if summary == "" {
		summary = "-"
	}
	lines = append(lines, boxLine(summaryStyle.Render(truncate(summary, cardContentWidth)), borderStyle))

	// Bottom border
	lines = append(lines, boxBottom(borderStyle))

	return strings.Join(lines, "\n")
}

// truncate shortens a string to maxLen with ellipsis.
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// getRoleColor returns the palette.md color for each role.
func getRoleColor(role schema.Role) string {
	if color, ok := roleColors[role]; ok {
		return color
	}
	return "#8A93A5"
}

// getStateColor returns the palette.md color for each state.
func getStateColor(state schema.AgentState) string {
	if color, ok := stateColors[state]; ok {
		return color
	}
	return "#6E7681"
}

// getStateStyle returns the lipgloss style for each state per visual rules.
func getStateStyle(state schema.AgentState) lipgloss.Style {
	color := getStateColor(state)

	switch state {
	case schema.StateRunning:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true)
	case schema.StateWaiting:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Faint(true)
	case schema.StateBlocked:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true)
	case schema.StateError:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true).
			Blink(true)
	case schema.StateDone:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color))
	case schema.StateFailed:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Strikethrough(true)
	case schema.StateCancelled:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Strikethrough(true)
	case schema.StateIdle:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Faint(true)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#8A93A5"))
	}
}

// StripAnsi removes ANSI escape codes from a string for testing.
func StripAnsi(s string) string {
	var result []byte
	inEsc := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') {
				inEsc = false
			}
			continue
		}
		result = append(result, s[i])
	}
	return string(result)
}
