package main

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the TUI
type Theme struct {
	Primary      lipgloss.Color
	Secondary    lipgloss.Color
	Text         lipgloss.Color
	Muted        lipgloss.Color
	Border       lipgloss.Color
	Error        lipgloss.Color
	Background   lipgloss.Color
	Focused      lipgloss.Color
	ButtonActive lipgloss.Color
	ButtonBg     lipgloss.Color
}

// DefaultTheme returns the default color scheme
func DefaultTheme() Theme {
	return Theme{
		Primary:      lipgloss.Color("205"), // Pink/Magenta
		Secondary:    lipgloss.Color("170"), // Purple
		Text:         lipgloss.Color("231"), // White
		Muted:        lipgloss.Color("240"), // Gray
		Border:       lipgloss.Color("240"), // Gray
		Error:        lipgloss.Color("196"), // Red
		Background:   lipgloss.Color("235"), // Dark Gray
		Focused:      lipgloss.Color("170"), // Purple
		ButtonActive: lipgloss.Color("170"), // Purple
		ButtonBg:     lipgloss.Color("0"),   // Black
	}
}

// CyberpunkTheme returns a cyberpunk-inspired color scheme
func CyberpunkTheme() Theme {
	return Theme{
		Primary:      lipgloss.Color("51"),  // Cyan
		Secondary:    lipgloss.Color("201"), // Magenta
		Text:         lipgloss.Color("231"), // White
		Muted:        lipgloss.Color("240"), // Gray
		Border:       lipgloss.Color("51"),  // Cyan
		Error:        lipgloss.Color("196"), // Red
		Background:   lipgloss.Color("235"), // Dark Gray
		Focused:      lipgloss.Color("201"), // Magenta
		ButtonActive: lipgloss.Color("51"),  // Cyan
		ButtonBg:     lipgloss.Color("0"),   // Black
	}
}

// DraculaTheme returns a Dracula-inspired color scheme
func DraculaTheme() Theme {
	return Theme{
		Primary:      lipgloss.Color("141"), // Purple
		Secondary:    lipgloss.Color("212"), // Pink
		Text:         lipgloss.Color("231"), // White
		Muted:        lipgloss.Color("61"),  // Comment Gray
		Border:       lipgloss.Color("61"),  // Comment Gray
		Error:        lipgloss.Color("203"), // Red
		Background:   lipgloss.Color("235"), // Dark Background
		Focused:      lipgloss.Color("141"), // Purple
		ButtonActive: lipgloss.Color("212"), // Pink
		ButtonBg:     lipgloss.Color("0"),   // Black
	}
}

// GruvboxTheme returns a Gruvbox-inspired color scheme
func GruvboxTheme() Theme {
	return Theme{
		Primary:      lipgloss.Color("214"), // Orange
		Secondary:    lipgloss.Color("142"), // Green
		Text:         lipgloss.Color("223"), // Light Beige
		Muted:        lipgloss.Color("245"), // Gray
		Border:       lipgloss.Color("245"), // Gray
		Error:        lipgloss.Color("167"), // Red
		Background:   lipgloss.Color("235"), // Dark Brown
		Focused:      lipgloss.Color("142"), // Green
		ButtonActive: lipgloss.Color("214"), // Orange
		ButtonBg:     lipgloss.Color("0"),   // Black
	}
}

// ApplyTheme applies a theme to the global color variables
func ApplyTheme(theme Theme) {
	colorPrimary = theme.Primary
	colorSecondary = theme.Secondary
	colorText = theme.Text
	colorMuted = theme.Muted
	colorBorder = theme.Border
	colorError = theme.Error
	colorBackground = theme.Background
	colorFocused = theme.Focused
	colorButtonActive = theme.ButtonActive
	colorButtonBg = theme.ButtonBg

	// Recreate styles with new colors
	initStyles()
}

func init() {
	// Apply default theme on package initialization
	ApplyTheme(DefaultTheme())
}

// initStyles initializes all styles with current colors
func initStyles() {
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
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorBackground)

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
}
