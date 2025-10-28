package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"my-cli-app/internal/password"
)

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "esc":
			return m.handleEscape()

		case "r":
			if m.CurrentView == ViewMain {
				return m.refreshPassword()
			}

		case "1", "g":
			if m.CurrentView == ViewWelcome {
				return m.startGenerate()
			}

		case "2", "l":
			if m.CurrentView == ViewMain {
				return m.startList()
			}
			if m.CurrentView == ViewWelcome {
				return m.startList()
			}

		case "3", "s":
			if m.CurrentView == ViewMain && m.Password != "" {
				return m.startSave()
			}
			if m.CurrentView == ViewWelcome {
				return m.startSettings()
			}

		case "4", "q":
			if m.CurrentView == ViewWelcome {
				return m.confirmQuit()
			}

		case "c":
			if m.CurrentView == ViewMain && m.Password != "" {
				return m.copyToClipboard()
			}

		case "tab":
			if m.CurrentView == ViewSave {
				if m.SiteInput.Focused() {
					m.SiteInput.Blur()
					m.UsernameInput.Focus()
				} else {
					m.UsernameInput.Blur()
					m.SiteInput.Focus()
				}
				return m, nil
			}

		case "up":
			if m.CurrentView == ViewWelcome && m.MenuCursor > 0 {
				m.MenuCursor--
				return m, nil
			} else if m.CurrentView == ViewList && m.Cursor > 0 {
				m.Cursor--
				return m, nil
			}

		case "down":
			if m.CurrentView == ViewWelcome && m.MenuCursor < 3 {
				m.MenuCursor++
				return m, nil
			} else if m.CurrentView == ViewList {
				filtered := m.FilterPasswords()
				if m.Cursor < len(filtered)-1 {
					m.Cursor++
				}
				return m, nil
			}
		}

	case clearStatusMsg:
		if time.Now().After(m.StatusExpiry) {
			m.StatusMessage = ""
		}
		return m, nil
	}

	// Update the appropriate text input based on current view
	switch m.CurrentView {
	case ViewSettings:
		m.LengthInput, cmd = m.LengthInput.Update(msg)
	case ViewSave:
		if m.SiteInput.Focused() {
			m.SiteInput, cmd = m.SiteInput.Update(msg)
		} else {
			m.UsernameInput, cmd = m.UsernameInput.Update(msg)
		}
	case ViewList:
		oldFilter := m.FilterText
		m.FilterInput, cmd = m.FilterInput.Update(msg)
		m.FilterText = m.FilterInput.Value()
		if m.FilterText != oldFilter {
			m.Cursor = 0 // Reset cursor when filter changes
		}
	}

	return m, cmd
}

// handleEnter processes the Enter key for different views
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.CurrentView {
	case ViewWelcome:
		switch m.MenuCursor {
		case 0:
			return m.startGenerate()
		case 1:
			return m.startList()
		case 2:
			return m.startSettings()
		case 3:
			return m, tea.Quit
		}

	case ViewSettings:
		inputLen, err := strconv.Atoi(m.LengthInput.Value())
		if err != nil || inputLen < 4 || inputLen > 128 {
			inputLen = 16
			m.StatusMessage = "Invalid length. Using default: 16"
		}
		m.Length = inputLen

		newPass, err := password.GeneratePassword(m.Length)
		if err != nil {
			m.Password = ""
			return m, m.SetStatus("Error: " + err.Error())
		}

		m.Password = newPass
		m.CurrentView = ViewMain
		m.LengthInput.Blur()
		return m, m.SetStatus("Generated " + strconv.Itoa(m.Length) + "-character password")

	case ViewSave:
		siteName := strings.TrimSpace(m.SiteInput.Value())
		username := strings.TrimSpace(m.UsernameInput.Value())
		if siteName == "" {
			return m, m.SetStatus("Site name cannot be empty")
		}

		if err := password.SavePasswordToCSV(siteName, username, m.Password); err != nil {
			return m, m.SetStatus("Save failed: " + err.Error())
		}

		// Copy to clipboard as well
		if err := clipboard.WriteAll(m.Password); err != nil {
			return m, m.SetStatus("✓ Saved to passwords.csv (clipboard copy failed)")
		}

		m.CurrentView = ViewMain
		m.SiteInput.Blur()
		m.SiteInput.SetValue("")
		m.UsernameInput.Blur()
		m.UsernameInput.SetValue("")
		return m, m.SetStatus("Saved to passwords.csv & copied to clipboard")

	case ViewConfirmQuit:
		return m, tea.Quit

	case ViewList:
		filtered := m.FilterPasswords()
		if len(filtered) > 0 && m.Cursor >= 0 && m.Cursor < len(filtered) {
			record := filtered[m.Cursor]
			if len(record) >= 3 {
				password := record[2]
				if err := clipboard.WriteAll(password); err != nil {
					return m, m.SetStatus("Failed to copy: " + err.Error())
				}
				return m, m.SetStatus("Password copied to clipboard")
			}
		}
	}

	return m, nil
}

// handleEscape processes the Esc key for different views
func (m Model) handleEscape() (tea.Model, tea.Cmd) {
	switch m.CurrentView {
	case ViewSettings:
		m.CurrentView = ViewWelcome
		m.LengthInput.Blur()
		return m, m.SetStatus("Back to main menu")

	case ViewMain:
		m.CurrentView = ViewWelcome
		return m, m.SetStatus("Back to main menu")

	case ViewSave:
		m.CurrentView = ViewMain
		m.SiteInput.Blur()
		m.SiteInput.SetValue("")
		m.UsernameInput.Blur()
		m.UsernameInput.SetValue("")
		return m, m.SetStatus("Save cancelled")

	case ViewConfirmQuit:
		m.CurrentView = ViewWelcome
		return m, nil

	case ViewList:
		m.CurrentView = ViewWelcome
		m.FilterInput.Blur()
		m.FilterInput.SetValue("")
		m.FilterText = ""
		m.Cursor = 0
		return m, m.SetStatus("Back to main menu")
	}

	return m, nil
}

// refreshPassword generates a new password with the current length
func (m Model) refreshPassword() (tea.Model, tea.Cmd) {
	newPass, err := password.GeneratePassword(m.Length)
	if err != nil {
		return m, m.SetStatus("Error: " + err.Error())
	}
	m.Password = newPass
	return m, m.SetStatus("Password refreshed")
}

// startSave transitions to the save view
func (m Model) startSave() (tea.Model, tea.Cmd) {
	m.CurrentView = ViewSave
	m.SiteInput.SetValue("")
	m.SiteInput.Focus()
	m.UsernameInput.SetValue("")
	return m, nil
}

// startList transitions to the list view
func (m Model) startList() (tea.Model, tea.Cmd) {
	m.CurrentView = ViewList
	m.FilterInput.SetValue("")
	m.FilterInput.Focus()
	m.Cursor = 0
	passwords, err := password.LoadPasswordsFromCSV()
	if err != nil {
		m.SavedPasswords = [][]string{}
	} else {
		m.SavedPasswords = passwords
	}
	return m, nil
}

// startGenerate transitions to password generation (settings view)
func (m Model) startGenerate() (tea.Model, tea.Cmd) {
	m.CurrentView = ViewSettings
	m.LengthInput.SetValue("")
	m.LengthInput.Focus()
	return m, nil
}

// startSettings transitions to the settings view
func (m Model) startSettings() (tea.Model, tea.Cmd) {
	m.CurrentView = ViewSettings
	m.LengthInput.SetValue("")
	m.LengthInput.Focus()
	return m, nil
}

// copyToClipboard copies the current password to clipboard
func (m Model) copyToClipboard() (tea.Model, tea.Cmd) {
	if err := clipboard.WriteAll(m.Password); err != nil {
		return m, m.SetStatus("Failed to copy: " + err.Error())
	}
	return m, m.SetStatus("Copied to clipboard")
}

// confirmQuit transitions to quit confirmation
func (m Model) confirmQuit() (tea.Model, tea.Cmd) {
	m.CurrentView = ViewConfirmQuit
	return m, nil
}
