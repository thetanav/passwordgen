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
)

// viewState represents the current screen
type viewState string

const (
	viewSettings viewState = "settings"
	viewMain     viewState = "main"
	viewSave     viewState = "save"
)

// model holds the application state
type model struct {
	password       string
	lengthInput    textinput.Model
	saveInput      textinput.Model
	view           viewState
	length         int
	statusMessage  string
	statusExpiry   time.Time
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
func savePasswordToCSV(siteName, password string) error {
	file, err := os.OpenFile(csvFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{siteName, password}
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
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
	saveInput := textinput.New()
	saveInput.Placeholder = "example.com"
	saveInput.Prompt = "Site/Service Name: "
	saveInput.CharLimit = 100
	saveInput.Width = 40

	return model{
		lengthInput: lengthInput,
		saveInput:   saveInput,
		view:        viewSettings,
		length:      defaultPasswordLength,
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

		case "s":
			if m.view == viewMain && m.password != "" {
				return m.startSave()
			}

		case "c":
			if m.view == viewMain && m.password != "" {
				return m.copyToClipboard()
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
		m.saveInput, cmd = m.saveInput.Update(msg)
	}

	return m, cmd
}

// handleEnter processes the Enter key for different views
func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.view {
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
		return m, m.setStatus(fmt.Sprintf("✓ Generated %d-character password", m.length))

	case viewSave:
		siteName := strings.TrimSpace(m.saveInput.Value())
		if siteName == "" {
			return m, m.setStatus("⚠ Site name cannot be empty")
		}

		if err := savePasswordToCSV(siteName, m.password); err != nil {
			return m, m.setStatus(fmt.Sprintf("✗ Save failed: %v", err))
		}

		// Copy to clipboard as well
		if err := clipboard.WriteAll(m.password); err != nil {
			return m, m.setStatus(fmt.Sprintf("✓ Saved to %s (clipboard copy failed)", csvFilename))
		}

		m.view = viewMain
		m.saveInput.Blur()
		m.saveInput.SetValue("")
		return m, m.setStatus(fmt.Sprintf("✓ Saved to %s & copied to clipboard", csvFilename))
	}

	return m, nil
}

// handleEscape processes the Esc key for different views
func (m model) handleEscape() (tea.Model, tea.Cmd) {
	switch m.view {
	case viewSettings:
		m.length = defaultPasswordLength
		newPass, _ := generatePassword(m.length)
		m.password = newPass
		m.view = viewMain
		m.lengthInput.Blur()
		return m, m.setStatus(fmt.Sprintf("Using default length: %d", defaultPasswordLength))

	case viewSave:
		m.view = viewMain
		m.saveInput.Blur()
		m.saveInput.SetValue("")
		return m, m.setStatus("Save cancelled")
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
	return m, m.setStatus("✓ Password refreshed")
}

// startSave transitions to the save view
func (m model) startSave() (tea.Model, tea.Cmd) {
	m.view = viewSave
	m.saveInput.SetValue("")
	m.saveInput.Focus()
	return m, nil
}

// copyToClipboard copies the current password to clipboard
func (m model) copyToClipboard() (tea.Model, tea.Cmd) {
	if err := clipboard.WriteAll(m.password); err != nil {
		return m, m.setStatus(fmt.Sprintf("✗ Failed to copy: %v", err))
	}
	return m, m.setStatus("✓ Copied to clipboard")
}

func (m model) View() string {
	var s strings.Builder
    s.WriteString("\n")
	switch m.view {
	case viewSettings:
		s.WriteString("  ╔══════════════════════════════════════════════════════════╗\n")
		s.WriteString("  ║                     Password Generator                   ║\n")
		s.WriteString("  ╚══════════════════════════════════════════════════════════╝\n\n")
		s.WriteString(m.lengthInput.View())
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("Range: %d-%d characters (default: %d)\n\n", 
			minPasswordLength, maxPasswordLength, defaultPasswordLength))
		s.WriteString("Press [Enter] to generate • [Esc] for default • [Q] to quit\n")

	case viewSave:
		s.WriteString("  ╔══════════════════════════════════════════════════════════╗\n")
		s.WriteString("  ║                      Save Password                       ║\n")
		s.WriteString("  ╚══════════════════════════════════════════════════════════╝\n\n")
		s.WriteString(m.saveInput.View())
		s.WriteString("\n\n")
		s.WriteString("✓ Copied to clipboard\n\n")
		s.WriteString("Press [Enter] to save • [Esc] to cancel\n")

	case viewMain:
		s.WriteString("  ╔══════════════════════════════════════════════════════════╗\n")
		s.WriteString("  ║                Secure Password Generator                 ║\n")
		s.WriteString("  ╚══════════════════════════════════════════════════════════╝\n\n")
		
		if m.password != "" {
			s.WriteString("Generated Password:\n")
			s.WriteString("┌────────────────────────────────────────────────────────┐\n")
			s.WriteString(fmt.Sprintf("│ %s%s │\n", m.password, strings.Repeat(" ", 54-len(m.password))))
			s.WriteString("└────────────────────────────────────────────────────────┘\n\n")
			s.WriteString(fmt.Sprintf("Length: %d characters\n\n", m.length,))
		}

		s.WriteString("Press [R] to refresh • [C] to copy • [S] to save & copy • [Q] to quit\n")
	}

	// Add status message if present
	if m.statusMessage != "" && time.Now().Before(m.statusExpiry) {
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("%s\n", m.statusMessage))
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