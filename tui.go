package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	minWidth  = 120
	minHeight = 30
)

// Color palette
var (
	colorPrimary      = lipgloss.Color("205")
	colorSecondary    = lipgloss.Color("170")
	colorText         = lipgloss.Color("231")
	colorMuted        = lipgloss.Color("240")
	colorBorder       = lipgloss.Color("240")
	colorError        = lipgloss.Color("196")
	colorBackground   = lipgloss.Color("235")
	colorFocused      = lipgloss.Color("170")
	colorButtonActive = lipgloss.Color("170")
	colorButtonBg     = lipgloss.Color("0")
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	fieldStyle = lipgloss.NewStyle()

	focusedFieldStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorFocused)

	inputStyle = lipgloss.NewStyle().
			Foreground(colorText)

	hintStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)

	buttonBlurredStyle = lipgloss.NewStyle().
				Padding(0, 2).
				MarginRight(1).
				Foreground(colorText).
				Background(colorButtonBg)

	buttonFocusedStyle = lipgloss.NewStyle().
				Padding(0, 2).
				MarginRight(1).
				Bold(true).
				Foreground(colorButtonBg).
				Background(colorButtonActive)

	panelBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(colorMuted)

	errorPopupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorError).
			Background(colorBackground).
			Foreground(colorText).
			Padding(1, 2)

	errorTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorError).
			Align(lipgloss.Center)

	errorMsgStyle = lipgloss.NewStyle().
			Foreground(colorText)

	errorHintStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true).
			Align(lipgloss.Center)

	contentStyle = lipgloss.NewStyle()

	separatorStyle = lipgloss.NewStyle().
			Foreground(colorMuted)
)

type lastOperation struct {
	success      bool
	operation    string
	params       string
	protocol     string
	whenHappened time.Time
}

// Key bindings
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Quit     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Left, k.Right, k.Up, k.Down, k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.ShiftTab, k.Up, k.Down},
		{k.Left, k.Right, k.Enter, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "scroll down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev option"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next option"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev field"),
	),
}

// Message types for async operations
type operationResultMsg struct {
	operation string
	params    string
	protocol  string
	result    string
	err       error
}

type focusField int

const (
	focusProtocol focusField = iota
	focusEnrollment
	focusOperation
	focusParams
	focusSubmit
)

type model struct {
	// Form fields
	protocolIdx  int
	protocols    []string
	enrollment   textinput.Model
	operationIdx int
	operations   []string
	paramsInput  textinput.Model

	// Components
	help     help.Model
	keys     keyMap
	progress progress.Model

	// State
	focusIndex   focusField
	result       string
	errorMsg     string
	showError    bool
	isValidation bool
	loading      bool

	// App settings
	settings *Settings
	width    int
	height   int

	// viewport
	ready           bool
	viewport        viewport.Model
	leftPanelWidth  int
	rightPanelWidth int

	renderer       *glamour.TermRenderer
	lastOperations []lastOperation
}

const defaultEnrollmentID = "538349"

func initialModel(settings *Settings) model {
	enrollment := textinput.New()
	enrollment.Placeholder = "Enter enrollment ID"
	enrollment.CharLimit = 50
	enrollment.PromptStyle = inputStyle
	enrollment.TextStyle = inputStyle
	enrollment.SetValue(defaultEnrollmentID)

	paramsInput := textinput.New()
	paramsInput.Placeholder = "Enter parameters (JSON or comma-separated)"
	paramsInput.CharLimit = 200
	paramsInput.PromptStyle = inputStyle
	paramsInput.TextStyle = inputStyle

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(minWidth/2-10),
	)
	if err != nil {
		panic(err)
	}

	prog := progress.New(progress.WithDefaultGradient())

	return model{
		protocolIdx:    0,
		protocols:      []string{"json", "string", "protobuf"},
		enrollment:     enrollment,
		operationIdx:   0,
		operations:     []string{"echo", "sum", "timestamp", "history", "status"},
		paramsInput:    paramsInput,
		help:           help.New(),
		keys:           keys,
		progress:       prog,
		focusIndex:     focusProtocol,
		settings:       settings,
		width:          minWidth,
		height:         minHeight,
		renderer:       r,
		lastOperations: make([]lastOperation, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

		m.leftPanelWidth = minWidth / 2
		m.rightPanelWidth = minWidth - m.leftPanelWidth

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(m.rightPanelWidth-10, minHeight-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent("")
			m.ready = true
			m.viewport.KeyMap = viewport.DefaultKeyMap()
		} else {
			m.viewport.Width = m.rightPanelWidth
			m.viewport.Height = minHeight - verticalMarginHeight
		}
		return m, nil

	case tea.KeyMsg:
		if m.showError {
			m.showError = false
			m.errorMsg = ""
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			m.focusIndex = (m.focusIndex + 1) % 5
			m.updateFocus()

		case key.Matches(msg, m.keys.ShiftTab):
			m.focusIndex = (m.focusIndex - 1 + 5) % 5
			m.updateFocus()

		case key.Matches(msg, m.keys.Left):
			switch m.focusIndex {
			case focusProtocol:
				if m.protocolIdx > 0 {
					m.protocolIdx--
				}
			case focusOperation:
				if m.operationIdx > 0 {
					m.operationIdx--
				}
			}

		case key.Matches(msg, m.keys.Right):
			switch m.focusIndex {
			case focusProtocol:
				if m.protocolIdx < len(m.protocols)-1 {
					m.protocolIdx++
				}
			case focusOperation:
				if m.operationIdx < len(m.operations)-1 {
					m.operationIdx++
				}
			}
		case key.Matches(msg, m.keys.Enter):
			if m.focusIndex == focusSubmit {
				if err := m.validate(); err != nil {
					m.showError = true
					m.isValidation = true
					m.errorMsg = err.Error()
					return m, nil
				}
				m.loading = true
				return m, m.executeOperation()
			}
		case key.Matches(msg, m.keys.Down):
			m.viewport.ScrollDown(1)

		case key.Matches(msg, m.keys.Up):
			m.viewport.ScrollUp(1)
		}

	case operationResultMsg:
		newOp := lastOperation{
			whenHappened: time.Now(),
			operation:    msg.operation,
			protocol:     msg.protocol,
			params:       msg.params,
		}
		m.loading = false
		if msg.err != nil {
			newOp.success = false
			m.showError = true
			m.isValidation = false
			m.errorMsg = msg.err.Error()
		} else {
			newOp.success = true
			m.result = msg.result
			content := ""

			markdown := "```json\n" + m.result + "\n```"
			rendered, err := m.renderer.Render(markdown)
			if err == nil {
				content = rendered
			} else {
				// Fallback to plain text if glamour fails
				content = m.result
			}
			m.viewport.SetContent(content)
		}

		if len(m.lastOperations) >= 5 {
			m.lastOperations = m.lastOperations[1:]
		}

		m.lastOperations = append(m.lastOperations, newOp)

		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(msg)

	// Handle text input updates
	switch m.focusIndex {
	case focusEnrollment:
		m.enrollment, cmd = m.enrollment.Update(msg)
	case focusParams:
		m.paramsInput, cmd = m.paramsInput.Update(msg)
	}

	return m, cmd
}

func (m *model) updateFocus() {
	m.enrollment.Blur()
	m.paramsInput.Blur()

	switch m.focusIndex {
	case focusEnrollment:
		m.enrollment.Focus()
	case focusParams:
		m.paramsInput.Focus()
	}
}

func (m model) headerView() string {
	localTitleStyle := titleStyle.Width(m.viewport.Width)
	title := localTitleStyle.Render("╔═ Response ═╗")
	return title
}

func (m model) footerView() string {
	info := panelBorderStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m model) validate() error {
	if m.enrollment.Value() == "" {
		return fmt.Errorf("enrollment ID is required")
	}

	op := m.operations[m.operationIdx]
	params := m.paramsInput.Value()

	switch op {
	case "echo":
		if params == "" {
			return fmt.Errorf("echo operation requires a message parameter")
		}
	case "sum":
		if params == "" {
			return fmt.Errorf("sum operation requires numbers (comma-separated)")
		}
		// Validate numbers
		parts := strings.Split(params, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if _, err := strconv.Atoi(p); err != nil {
				return fmt.Errorf("invalid number: %s", p)
			}
		}
	case "history":
		if params == "" {
			return fmt.Errorf("history operation requires a limit parameter")
		}
		limit, err := strconv.Atoi(params)
		if err != nil {
			return fmt.Errorf("limit must be a number")
		}
		if limit < 1 || limit > 100 {
			return fmt.Errorf("limit must be between 1 and 100")
		}
	case "status":
		if params != "" && params != "true" && params != "false" {
			return fmt.Errorf("status parameter must be 'true' or 'false' for detailed mode")
		}
	}

	return nil
}

func (m model) executeOperation() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// Get protocol serde and server address
		var serde Serde
		var serverAddress string

		switch m.protocols[m.protocolIdx] {
		case "json":
			serde = &JSONSerde{}
			serverAddress = m.settings.App.JSONProtocolServerAddress
		case "string":
			serde = &StringSerde{}
			serverAddress = m.settings.App.StringProtocolServerAddress
		case "protobuf":
			serde = &ProtobufSerde{}
			serverAddress = m.settings.App.ProtobufProtocolServerAddress
		}

		// 1. Authenticate
		authClient := NewAppLayerClient[*AuthRequest, *AuthResponse](
			serde,
			DefaultTCPRoundTripper,
			&m.settings.App,
		)

		authReq := &AuthRequest{
			StudentID: m.enrollment.Value(),
			Timestamp: time.Now(),
		}

		authResp, err := authClient.Auth(ctx, serverAddress, authReq)
		if err != nil {
			return operationResultMsg{err: fmt.Errorf("authentication failed: %w", err)}
		}

		token := authResp.Token
		var result string

		// 2. Perform operation
		op := m.operations[m.operationIdx]
		switch op {
		case "echo":
			client := NewAppLayerClient[*EchoRequest, *EchoResponse](
				serde,
				DefaultTCPRoundTripper,
				&m.settings.App,
			)
			req := &EchoRequest{Message: m.paramsInput.Value()}
			var resp EchoResponse
			err = client.Do(ctx, serverAddress, req, &resp, token)
			if err == nil {
				result = formatResponse("Echo Response", resp)
			}

		case "sum":
			client := NewAppLayerClient[*SumRequest, *SumResponse](
				serde,
				DefaultTCPRoundTripper,
				&m.settings.App,
			)
			parts := strings.Split(m.paramsInput.Value(), ",")
			numbers := make([]int, 0, len(parts))
			for _, p := range parts {
				n, _ := strconv.Atoi(strings.TrimSpace(p))
				numbers = append(numbers, n)
			}
			req := &SumRequest{Numbers: numbers}
			var resp SumResponse
			err = client.Do(ctx, serverAddress, req, &resp, token)
			if err == nil {
				result = formatResponse("Sum Response", resp)
			}

		case "timestamp":
			client := NewAppLayerClient[*TimestampRequest, *TimestampResponse](
				serde,
				DefaultTCPRoundTripper,
				&m.settings.App,
			)
			req := &TimestampRequest{}
			var resp TimestampResponse
			err = client.Do(ctx, serverAddress, req, &resp, token)
			if err == nil {
				result = formatResponse("Timestamp Response", resp)
			}

		case "history":
			client := NewAppLayerClient[*HistoryRequest, *HistoryResponse](
				serde,
				DefaultTCPRoundTripper,
				&m.settings.App,
			)
			limit, _ := strconv.Atoi(m.paramsInput.Value())
			req := &HistoryRequest{Limit: limit}
			var resp HistoryResponse
			err = client.Do(ctx, serverAddress, req, &resp, token)
			if err == nil {
				result = formatResponse("History Response", resp)
			}

		case "status":
			client := NewAppLayerClient[*StatusRequest, *StatusResponse](
				serde,
				DefaultTCPRoundTripper,
				&m.settings.App,
			)
			detailed := m.paramsInput.Value() == "true"
			req := &StatusRequest{Detailed: detailed}
			var resp StatusResponse
			err = client.Do(ctx, serverAddress, req, &resp, token)
			if err == nil {
				result = formatResponse("Status Response", resp)
			}
		}

		if err != nil {
			// Still try to logout
			logoutClient := NewAppLayerClient[*LogoutRequest, *LogoutResponse](
				serde,
				DefaultTCPRoundTripper,
				&m.settings.App,
			)
			logoutReq := &LogoutRequest{}
			logoutClient.Logout(ctx, serverAddress, logoutReq, token)
			return operationResultMsg{
				err:       fmt.Errorf("operation failed: %w", err),
				protocol:  m.protocols[m.protocolIdx],
				params:    m.paramsInput.Value(),
				operation: op,
			}
		}

		// 3. Logout
		logoutClient := NewAppLayerClient[*LogoutRequest, *LogoutResponse](
			serde,
			DefaultTCPRoundTripper,
			&m.settings.App,
		)
		logoutReq := &LogoutRequest{}
		_, err = logoutClient.Logout(ctx, serverAddress, logoutReq, token)
		if err != nil {
			result += "\n\nWarning: Logout failed: " + err.Error()
		}

		return operationResultMsg{
			result:    result,
			protocol:  m.protocols[m.protocolIdx],
			params:    m.paramsInput.Value(),
			operation: op,
		}
	}
}

func formatResponse(title string, data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%s\nError formatting response: %v", title, err)
	}
	return fmt.Sprintf("%s\n%s", title, string(jsonData))
}

func (m model) View() string {
	if m.width < minWidth || m.height < minHeight {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("Terminal too small. Please resize. Minimum size: %dx%d, found %dx%d", minWidth, minHeight, m.width, m.height),
		)
	}

	leftPanel := m.renderLeftPanel(m.leftPanelWidth)
	rightPanel := m.renderRightPanel(m.rightPanelWidth)

	// Combine panels side by side
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	// Add help at the bottom
	helpView := m.help.View(m.keys)
	separator := separatorStyle.Width(minWidth).Render(strings.Repeat("─", minWidth))

	title := titleStyle.Render("╔═ Triprotocol Client ═╗")

	fullView := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		mainContent,
		"",
		separator,
		helpView,
	)

	// Center everything both horizontally and vertically
	centeredView := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		fullView,
	)

	// Show error popup if needed
	if m.showError {
		return m.renderErrorPopup(centeredView)
	}

	return centeredView
}

func (m model) renderLeftPanel(width int) string {
	// Apply width to styles
	localTitleStyle := titleStyle.Width(width)
	localFieldStyle := fieldStyle.Width(width)
	localFocusedStyle := focusedFieldStyle.Width(width)
	localHintStyle := hintStyle.Width(width)

	// Title
	title := localTitleStyle.Render("╔═ Options ═╗")

	// Protocol section (styled as tabs)
	protocolLabel := "Protocol:"
	protocolLabelRendered := ""
	if m.focusIndex == focusProtocol {
		protocolLabelRendered = localFocusedStyle.Render("▸ " + protocolLabel)
	} else {
		protocolLabelRendered = localFieldStyle.Render("  " + protocolLabel)
	}

	// Render protocols as tabs
	var protocolTabs []string
	for i, protocol := range m.protocols {
		var style lipgloss.Style
		isActive := i == m.protocolIdx
		if isActive {
			style = buttonFocusedStyle
		} else {
			style = buttonBlurredStyle
		}
		protocolTabs = append(protocolTabs, style.Render(protocol))
	}
	protocolValue := localFieldStyle.Render("  " + lipgloss.JoinHorizontal(lipgloss.Top, protocolTabs...))

	// Enrollment section
	enrollmentLabel := "Enrollment ID:"
	enrollmentLabelRendered := ""
	if m.focusIndex == focusEnrollment {
		enrollmentLabelRendered = localFocusedStyle.Render("▸ " + enrollmentLabel)
	} else {
		enrollmentLabelRendered = localFieldStyle.Render("  " + enrollmentLabel)
	}
	enrollmentValue := localFieldStyle.Render("  " + m.enrollment.View())

	// Operation section
	operationLabel := "Operation:"
	operationLabelRendered := ""
	if m.focusIndex == focusOperation {
		operationLabelRendered = localFocusedStyle.Render("▸ " + operationLabel)
	} else {
		operationLabelRendered = localFieldStyle.Render("  " + operationLabel)
	}

	// Render operations as tabs
	var operationTabs []string
	for i, op := range m.operations {
		var style lipgloss.Style
		isActive := i == m.operationIdx
		if isActive {
			style = buttonFocusedStyle
		} else {
			style = buttonBlurredStyle
		}
		operationTabs = append(operationTabs, style.Render(op))
	}
	operationsList := localFieldStyle.Render("  " + lipgloss.JoinHorizontal(lipgloss.Top, operationTabs...))

	// Params section
	paramsLabel := "Parameters:"
	paramsLabelRendered := ""
	if m.focusIndex == focusParams {
		paramsLabelRendered = localFocusedStyle.Render("▸ " + paramsLabel)
	} else {
		paramsLabelRendered = localFieldStyle.Render("  " + paramsLabel)
	}
	paramsValue := localFieldStyle.Render("  " + m.paramsInput.View())
	paramsHint := localHintStyle.Render("  " + m.getParamsHint())

	// Submit button
	buttonRendered := ""
	if m.focusIndex == focusSubmit {
		buttonRendered = localFieldStyle.Render("  " + buttonFocusedStyle.Render("SUBMIT"))
	} else {
		buttonRendered = localFieldStyle.Render("  " + buttonBlurredStyle.Render("SUBMIT"))
	}

	lastOperations := []string{}

	for _, op := range m.lastOperations {
		operation := inputStyle.Render(op.operation)

		status := inputStyle.Render("✅")

		if !op.success {
			status = errorTitleStyle.Render("❌")
		}

		protocol := titleStyle.Render(op.protocol)
		params := inputStyle.Render(op.params)

		happenedAt := inputStyle.Render(time.Since(op.whenHappened).Truncate(time.Second).String())
		renderedOperation := lipgloss.JoinHorizontal(lipgloss.Top, status, "|", operation, " at ", happenedAt, "|", protocol, "|", params)
		lastOperations = append([]string{""}, lastOperations...)
		lastOperations = append([]string{renderedOperation}, lastOperations...)
	}

	operationsPane := ""
	if len(lastOperations) != 0 {
		operationsPane = lipgloss.JoinVertical(lipgloss.Left, lastOperations...)
		operationsPane = panelBorderStyle.Width(width).Padding(1).Render(operationsPane)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		protocolLabelRendered,
		protocolValue,
		"",
		enrollmentLabelRendered,
		enrollmentValue,
		"",
		operationLabelRendered,
		operationsList,
		"",
		paramsLabelRendered,
		paramsValue,
		paramsHint,
		"",
		buttonRendered,
		"",
		operationsPane,
	)
}

func (m model) renderRightPanel(width int) string {
	var content string
	if m.loading {
		progressBar := m.progress.ViewAs(0.5) // Simple animated progress
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			"⏳ Processing request...",
			"",
			progressBar,
		)
	} else if m.result == "" {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			"No results yet.",
			"Fill in the form and",
			"submit a request.",
		)
	} else {
		content = m.viewport.View()
	}

	panel := lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		"",
		content,
		"",
		m.footerView(),
	)

	return panelBorderStyle.Padding(1).Render(panel)
}

func (m model) renderErrorPopup(background string) string {
	errorTitle := "⚠️  Error"
	if m.isValidation {
		errorTitle = "⚠️  Validation Error"
	} else {
		errorTitle = "⚠️  Internal Error"
	}

	popupWidth := 60

	// Word wrap error message
	var msgLines []string
	words := strings.Fields(m.errorMsg)
	var line string
	maxLineWidth := popupWidth - 8
	for _, word := range words {
		if len(line)+len(word)+1 > maxLineWidth {
			if line != "" {
				msgLines = append(msgLines, line)
			}
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}
	if line != "" {
		msgLines = append(msgLines, line)
	}

	// Build popup content
	popupContent := lipgloss.JoinVertical(
		lipgloss.Center,
		errorTitleStyle.Render(errorTitle),
		"",
		errorMsgStyle.Render(strings.Join(msgLines, "\n")),
		"",
		errorHintStyle.Render("Press any key to close"),
	)

	popupWidget := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorError).
		Background(colorBackground).
		Foreground(colorText).
		Padding(1, 2).
		Width(popupWidth).
		Render(popupContent)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		popupWidget,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}

func (m model) getParamsHint() string {
	op := m.operations[m.operationIdx]
	switch op {
	case "echo":
		return "Enter a message"
	case "sum":
		return "Comma-separated numbers"
	case "timestamp":
		return "No parameters needed"
	case "history":
		return "Limit (1-100)"
	case "status":
		return "true/false for detailed"
	}
	return ""
}

func RunTUI() error {
	// Load configuration
	settings, err := LoadConfig[Settings]("TUI", BaseSettings)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	defer f.Close()
	logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	p := tea.NewProgram(
		initialModel(settings),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err = p.Run()
	return err
}
