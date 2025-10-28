package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// viewState represents the current screen
type viewState string

const (
	viewWelcome  viewState = "welcome"
	viewSettings viewState = "settings"
	viewMain     viewState = "main"
	viewSave     viewState = "save"
	viewList     viewState = "list"
)

// model holds the application state
type model struct {
	password       string
	lengthInput    textinput.Model
	siteInput      textinput.Model
	usernameInput  textinput.Model
	filterInput    textinput.Model
	view           viewState
	length         int
	statusMessage  string
	statusExpiry   time.Time
	savedPasswords [][]string
	filterText     string
	cursor         int
	menuCursor     int
}

// Character sets for password generation
const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
	symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

const (
	defaultPasswordLength = 16
	minPasswordLength     = 4
	maxPasswordLength     = 128
	csvFilename           = "passwords.csv"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Align(lipgloss.Center).
			Width(70)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("7")).
			Padding(1, 2)

	passwordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true).
			Align(lipgloss.Center)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("34"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true)
)

// generatePassword creates a secure random password of the specified length
func generatePassword(length int) (string, error) {
	if length < minPasswordLength {
		return "", fmt.Errorf("password length must be at least %d", minPasswordLength)
	}
	if length > maxPasswordLength {
		return "", fmt.Errorf("password length must not exceed %d", maxPasswordLength)
	}

	// Build character set
	var charset strings.Builder
	charset.WriteString(lowercase)
	charset.WriteString(uppercase)
	charset.WriteString(numbers)
	charset.WriteString(symbols)
	charsetStr := charset.String()

	// Generate random password
	password := make([]rune, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetStr))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = rune(charsetStr[n.Int64()])
	}

	return string(password), nil
}

// savePasswordToCSV appends the password to a CSV file
func savePasswordToCSV(siteName, username, password string) error {
	file, err := os.OpenFile(csvFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{siteName, username, password}
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
}

// loadPasswordsFromCSV reads all passwords from the CSV file
func loadPasswordsFromCSV() ([][]string, error) {
	file, err := os.Open(csvFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return [][]string{}, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	// Ensure all records have at least 3 columns (backward compatibility)
	for i, record := range records {
		for len(record) < 3 {
			record = append(record, "")
		}
		records[i] = record
	}

	return records, nil
}

// filterPasswords filters the saved passwords based on the filter text
func (m model) filterPasswords() [][]string {
	if m.filterText == "" {
		return m.savedPasswords
	}

	var filtered [][]string
	for _, record := range m.savedPasswords {
		if len(record) >= 2 {
			site := strings.ToLower(record[0])
			username := strings.ToLower(record[1])
			filter := strings.ToLower(m.filterText)

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
func (m *model) setStatus(msg string) tea.Cmd {
	m.statusMessage = msg
	m.statusExpiry = time.Now().Add(3 * time.Second)
	return clearStatusAfterDelay()
}

// initialModel creates the initial application state
func initialModel() model {
	// Length input for settings view
	lengthInput := textinput.New()
	lengthInput.Placeholder = "16"
	lengthInput.Focus()
	lengthInput.Prompt = "Password Length: "
	lengthInput.CharLimit = 3
	lengthInput.Width = 20

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

	return model{
		lengthInput:   lengthInput,
		siteInput:     siteInput,
		usernameInput: usernameInput,
		filterInput:   filterInput,
		view:          viewWelcome,
		length:        defaultPasswordLength,
		menuCursor:    0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "esc":
			return m.handleEscape()

		case "r":
			if m.view == viewMain {
				return m.refreshPassword()
			}

		case "g":
			if m.view == viewWelcome {
				return m.startGenerate()
			}

		case "l":
			if m.view == viewMain {
				return m.startList()
			}
			if m.view == viewWelcome {
				return m.startList()
			}

		case "s":
			if m.view == viewMain && m.password != "" {
				return m.startSave()
			}
			if m.view == viewWelcome {
				return m.startSettings()
			}

		case "c":
			if m.view == viewMain && m.password != "" {
				return m.copyToClipboard()
			}

		case "tab":
			if m.view == viewSave {
				if m.siteInput.Focused() {
					m.siteInput.Blur()
					m.usernameInput.Focus()
				} else {
					m.usernameInput.Blur()
					m.siteInput.Focus()
				}
				return m, nil
			}

		case "up":
			if m.view == viewWelcome && m.menuCursor > 0 {
				m.menuCursor--
				return m, nil
			} else if m.view == viewList && m.cursor > 0 {
				m.cursor--
				return m, nil
			}

		case "down":
			if m.view == viewWelcome && m.menuCursor < 3 {
				m.menuCursor++
				return m, nil
			} else if m.view == viewList {
				filtered := m.filterPasswords()
				if m.cursor < len(filtered)-1 {
					m.cursor++
				}
				return m, nil
			}
		}

	case clearStatusMsg:
		if time.Now().After(m.statusExpiry) {
			m.statusMessage = ""
		}
		return m, nil
	}

	// Update the appropriate text input based on current view
	switch m.view {
	case viewSettings:
		m.lengthInput, cmd = m.lengthInput.Update(msg)
	case viewSave:
		if m.siteInput.Focused() {
			m.siteInput, cmd = m.siteInput.Update(msg)
		} else {
			m.usernameInput, cmd = m.usernameInput.Update(msg)
		}
	case viewList:
		oldFilter := m.filterText
		m.filterInput, cmd = m.filterInput.Update(msg)
		m.filterText = m.filterInput.Value()
		if m.filterText != oldFilter {
			m.cursor = 0 // Reset cursor when filter changes
		}
	}

	return m, cmd
}

// handleEnter processes the Enter key for different views
func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.view {
	case viewWelcome:
		switch m.menuCursor {
		case 0:
			return m.startGenerate()
		case 1:
			return m.startList()
		case 2:
			return m.startSettings()
		case 3:
			return m, tea.Quit
		}

	case viewSettings:
		inputLen, err := strconv.Atoi(m.lengthInput.Value())
		if err != nil || inputLen < minPasswordLength || inputLen > maxPasswordLength {
			inputLen = defaultPasswordLength
			m.statusMessage = fmt.Sprintf("Invalid length. Using default: %d", defaultPasswordLength)
		}
		m.length = inputLen

		newPass, err := generatePassword(m.length)
		if err != nil {
			m.password = ""
			return m, m.setStatus(fmt.Sprintf("Error: %v", err))
		}

		m.password = newPass
		m.view = viewMain
		m.lengthInput.Blur()
		return m, m.setStatus(fmt.Sprintf("Generated %d-character password", m.length))

	case viewSave:
		siteName := strings.TrimSpace(m.siteInput.Value())
		username := strings.TrimSpace(m.usernameInput.Value())
		if siteName == "" {
			return m, m.setStatus("Site name cannot be empty")
		}

		if err := savePasswordToCSV(siteName, username, m.password); err != nil {
			return m, m.setStatus(fmt.Sprintf("Save failed: %v", err))
		}

		// Copy to clipboard as well
		if err := clipboard.WriteAll(m.password); err != nil {
			return m, m.setStatus(fmt.Sprintf("✓ Saved to %s (clipboard copy failed)", csvFilename))
		}

		m.view = viewMain
		m.siteInput.Blur()
		m.siteInput.SetValue("")
		m.usernameInput.Blur()
		m.usernameInput.SetValue("")
		return m, m.setStatus(fmt.Sprintf("Saved to %s & copied to clipboard", csvFilename))

	case viewList:
		filtered := m.filterPasswords()
		if len(filtered) > 0 && m.cursor >= 0 && m.cursor < len(filtered) {
			record := filtered[m.cursor]
			if len(record) >= 3 {
				password := record[2]
				if err := clipboard.WriteAll(password); err != nil {
					return m, m.setStatus(fmt.Sprintf("Failed to copy: %v", err))
				}
				return m, m.setStatus("Password copied to clipboard")
			}
		}
	}

	return m, nil
}

// handleEscape processes the Esc key for different views
func (m model) handleEscape() (tea.Model, tea.Cmd) {
	switch m.view {
	case viewSettings:
		m.view = viewWelcome
		m.lengthInput.Blur()
		return m, m.setStatus("Back to main menu")

	case viewMain:
		m.view = viewWelcome
		return m, m.setStatus("Back to main menu")

	case viewSave:
		m.view = viewMain
		m.siteInput.Blur()
		m.siteInput.SetValue("")
		m.usernameInput.Blur()
		m.usernameInput.SetValue("")
		return m, m.setStatus("Save cancelled")

	case viewList:
		m.view = viewWelcome
		m.filterInput.Blur()
		m.filterInput.SetValue("")
		m.filterText = ""
		m.cursor = 0
		return m, m.setStatus("Back to main menu")
	}

	return m, nil
}

// refreshPassword generates a new password with the current length
func (m model) refreshPassword() (tea.Model, tea.Cmd) {
	newPass, err := generatePassword(m.length)
	if err != nil {
		return m, m.setStatus(fmt.Sprintf("Error: %v", err))
	}
	m.password = newPass
	return m, m.setStatus("Password refreshed")
}

// startSave transitions to the save view
func (m model) startSave() (tea.Model, tea.Cmd) {
	m.view = viewSave
	m.siteInput.SetValue("")
	m.siteInput.Focus()
	m.usernameInput.SetValue("")
	return m, nil
}

// startList transitions to the list view
func (m model) startList() (tea.Model, tea.Cmd) {
	m.view = viewList
	m.filterInput.SetValue("")
	m.filterInput.Focus()
	m.cursor = 0
	m.savedPasswords = nil // Force reload
	return m, nil
}

// startGenerate transitions to password generation (settings view)
func (m model) startGenerate() (tea.Model, tea.Cmd) {
	m.view = viewSettings
	m.lengthInput.SetValue("")
	m.lengthInput.Focus()
	return m, nil
}

// startSettings transitions to the settings view
func (m model) startSettings() (tea.Model, tea.Cmd) {
	m.view = viewSettings
	m.lengthInput.SetValue("")
	m.lengthInput.Focus()
	return m, nil
}

// copyToClipboard copies the current password to clipboard
func (m model) copyToClipboard() (tea.Model, tea.Cmd) {
	if err := clipboard.WriteAll(m.password); err != nil {
		return m, m.setStatus(fmt.Sprintf("Failed to copy: %v", err))
	}
	return m, m.setStatus("Copied to clipboard")
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString("\n")
	switch m.view {
	case viewWelcome:
		s.WriteString(titleStyle.Render("Secure Password Manager"))
		s.WriteString("\n\n")

		options := []string{"Generate New Password", "View Saved Passwords", "Settings", "Quit Application"}
		for i, option := range options {
			if i == m.menuCursor {
				s.WriteString(selectedStyle.Render("> " + option))
			} else {
				s.WriteString("  " + option)
			}
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("Use arrow keys to navigate • [Enter] to select • Or press [G/L/S/Q] for quick access"))
		s.WriteString("\n")

	case viewSettings:
		s.WriteString(titleStyle.Render("Password Generator"))
		s.WriteString("\n\n")
		s.WriteString(m.lengthInput.View())
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render(fmt.Sprintf("Range: %d-%d characters (default: %d)",
			minPasswordLength, maxPasswordLength, defaultPasswordLength)))
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("Press [Enter] to generate • [Esc] for default • [Q] to quit"))

	case viewSave:
		s.WriteString(titleStyle.Render("Save Password"))
		s.WriteString("\n\n")
		s.WriteString(m.siteInput.View())
		s.WriteString("\n")
		s.WriteString(m.usernameInput.View())
		s.WriteString("\n\n")
		s.WriteString(successStyle.Render("Password copied to clipboard"))
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("Press [Tab] to switch fields • [Enter] to save • [Esc] to cancel"))

	case viewMain:
		s.WriteString(titleStyle.Render("Secure Password Generator"))
		s.WriteString("\n\n")

		if m.password != "" {
			s.WriteString("Generated Password:\n")
			passwordBox := boxStyle.Width(66).Render(passwordStyle.Render(m.password))
			s.WriteString(passwordBox)
			s.WriteString("\n\n")
			s.WriteString(infoStyle.Render(fmt.Sprintf("Length: %d characters", m.length)))
			s.WriteString("\n\n")
		}

		s.WriteString(infoStyle.Render("Press [R] to refresh • [C] to copy • [S] to save & copy • [L] to list • [Esc] to menu • [Q] to quit"))

	case viewList:
		s.WriteString(titleStyle.Render("Saved Passwords"))
		s.WriteString("\n\n")
		s.WriteString(m.filterInput.View())
		s.WriteString("\n\n")

		// Headers
		s.WriteString(infoStyle.Render("Site                | Username           | Password\n"))
		s.WriteString(infoStyle.Render(strings.Repeat("-", 60)))
		s.WriteString("\n")

		// Load passwords if not loaded
		if len(m.savedPasswords) == 0 {
			passwords, err := loadPasswordsFromCSV()
			if err != nil {
				s.WriteString(errorStyle.Render(fmt.Sprintf("Error loading passwords: %v", err)))
				s.WriteString("\n")
			} else {
				m.savedPasswords = passwords
			}
		}

		// Filter passwords
		filtered := m.filterPasswords()

		if len(filtered) == 0 {
			s.WriteString(infoStyle.Render("No passwords found."))
			s.WriteString("\n")
		} else {
			for i, record := range filtered {
				if len(record) >= 3 {
					site := record[0]
					username := record[1]
					password := record[2]

					// Truncate long fields
					siteTrunc := site
					if len(siteTrunc) > 18 {
						siteTrunc = siteTrunc[:15] + "..."
					}
					userTrunc := username
					if len(userTrunc) > 16 {
						userTrunc = userTrunc[:13] + "..."
					}
					line := fmt.Sprintf("%-18s | %-16s | %s", siteTrunc, userTrunc, strings.Repeat("•", len(password)))
					if i == m.cursor {
						s.WriteString(selectedStyle.Render("▶ " + line))
					} else {
						s.WriteString("  " + line)
					}
					s.WriteString("\n")
				}
			}
		}

		s.WriteString("\n")
		s.WriteString(infoStyle.Render("Press [↑/↓] to navigate • [Enter] to copy • [Esc] to menu • [Q] to quit"))
	}

	// Add status message if present
	if m.statusMessage != "" && time.Now().Before(m.statusExpiry) {
		s.WriteString("\n")
		if strings.Contains(m.statusMessage, "failed") || strings.Contains(m.statusMessage, "Error") || strings.Contains(m.statusMessage, "cannot") {
			s.WriteString(errorStyle.Render(m.statusMessage))
		} else {
			s.WriteString(successStyle.Render(m.statusMessage))
		}
		s.WriteString("\n")
	}

	return s.String()
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
