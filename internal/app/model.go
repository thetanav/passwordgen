package app

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewState represents the current screen
type ViewState string

const (
	ViewWelcome     ViewState = "welcome"
	ViewMain        ViewState = "main"
	ViewSave        ViewState = "save"
	ViewList        ViewState = "list"
	ViewConfirmQuit ViewState = "confirm_quit"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Align(lipgloss.Left).
			Width(70)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("250")).
			Padding(1, 2)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("35"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("160"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true)
)

// Model holds the application state
type Model struct {
	Password       string
	SiteInput      textinput.Model
	UsernameInput  textinput.Model
	FilterInput    textinput.Model
	CurrentView    ViewState
	Length         int
	StatusMessage  string
	StatusExpiry   time.Time
	SavedPasswords [][]string
	FilterText     string
	Cursor         int
	MenuCursor     int
}

// filterPasswords filters the saved passwords based on the filter text
func (m Model) filterPasswords() [][]string {
	if m.FilterText == "" {
		return m.SavedPasswords
	}

	var filtered [][]string
	for _, record := range m.SavedPasswords {
		if len(record) >= 2 {
			site := strings.ToLower(record[0])
			username := strings.ToLower(record[1])
			filter := strings.ToLower(m.FilterText)

			if strings.Contains(site, filter) || strings.Contains(username, filter) {
				filtered = append(filtered, record)
			}
		}
	}
	return filtered
}

// clearStatusMsg is a command that clears the status message after a delay
type clearStatusMsg struct{}

func clearStatusAfterDelay() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

// setStatus sets a temporary status message
func (m *Model) setStatus(msg string) tea.Cmd {
	m.StatusMessage = msg
	m.StatusExpiry = time.Now().Add(3 * time.Second)
	return clearStatusAfterDelay()
}

// initialModel creates the initial application state
func InitialModel() Model {
	// Site name input for save view
	siteInput := textinput.New()
	siteInput.Placeholder = "example.com"
	siteInput.Prompt = "Site/Service Name: "
	siteInput.CharLimit = 100
	siteInput.Width = 40

	// Username input for save view
	usernameInput := textinput.New()
	usernameInput.Placeholder = "username"
	usernameInput.Prompt = "Username: "
	usernameInput.CharLimit = 100
	usernameInput.Width = 40

	// Filter input for list view
	filterInput := textinput.New()
	filterInput.Placeholder = "filter..."
	filterInput.Prompt = "Filter: "
	filterInput.CharLimit = 100
	filterInput.Width = 40

	return Model{
		SiteInput:     siteInput,
		UsernameInput: usernameInput,
		FilterInput:   filterInput,
		CurrentView:   ViewWelcome,
		Length:        16, // defaultPasswordLength
		MenuCursor:    0,
	}
}
