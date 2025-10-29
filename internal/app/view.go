package app

import (
	"fmt"
	"strings"
	"time"
)

func (m Model) View() string {
	var s strings.Builder
	s.WriteString("\n")
	switch m.CurrentView {
	case ViewWelcome:
		s.WriteString(titleStyle.Render("Secure Password Manager"))
		s.WriteString("\n\n")

		options := []string{"1. Generate New Password", "2. View Saved Passwords", "3. Quit Application"}
		for i, option := range options {
			if i == m.MenuCursor {
				s.WriteString(selectedStyle.Render("> " + option))
			} else {
				s.WriteString("  " + option)
			}
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("Use arrow keys or numbers to navigate • [Enter] to select"))
		s.WriteString("\n")

	case ViewConfirmQuit:
		s.WriteString(titleStyle.Render("Confirm Quit"))
		s.WriteString("\n\n")
		s.WriteString("Are you sure you want to quit?\n\n")
		s.WriteString(infoStyle.Render("Press [Y] to quit • [N] or [Esc] to cancel"))
		s.WriteString("\n")

	case ViewSave:
		s.WriteString(titleStyle.Render("Save Password"))
		s.WriteString("\n\n")
		s.WriteString(m.SiteInput.View())
		s.WriteString("\n")
		s.WriteString(m.UsernameInput.View())
		s.WriteString("\n\n")
		s.WriteString(successStyle.Render("Password copied to clipboard"))
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("Press [Tab] to switch fields • [Enter] to save • [Esc] to cancel"))

	case ViewMain:
		s.WriteString(titleStyle.Render("Secure Password Generator"))
		s.WriteString("\n\n")

		if m.Password != "" {
			s.WriteString("Generated Password:\n")
			passwordBox := boxStyle.Width(66).Render(passwordStyle.Render(m.Password))
			s.WriteString(passwordBox)
			s.WriteString("\n\n")
			s.WriteString(infoStyle.Render(fmt.Sprintf("Length: %d characters", m.Length)))
			s.WriteString("\n\n")
		}

		s.WriteString(infoStyle.Render("Press [R] to refresh • [C] to copy • [S] to save & copy • [L] to list • [Esc] to menu • [Q] to quit"))

	case ViewList:
		s.WriteString(titleStyle.Render("Saved Passwords"))
		s.WriteString("\n\n")
		s.WriteString(m.FilterInput.View())
		s.WriteString("\n\n")

		var table strings.Builder

		// Headers
		headerLine := fmt.Sprintf("%-18s | %-18s | %-20s", "Site", "Username", "Password")
		table.WriteString(infoStyle.Render("  " + headerLine + "\n"))
		table.WriteString(infoStyle.Render("  " + strings.Repeat("-", 62) + "\n"))

		// Filter passwords
		filtered := m.filterPasswords()

		if len(filtered) == 0 {
			table.WriteString(infoStyle.Render("No passwords found."))
			table.WriteString("\n")
		} else {
			for i, record := range filtered {
				if len(record) >= 3 {
					site := record[0]
					username := record[1]

					// Truncate long fields
					siteTrunc := site
					if len(siteTrunc) > 18 {
						siteTrunc = siteTrunc[:15] + "..."
					}
					userTrunc := username
					if len(userTrunc) > 18 {
						userTrunc = userTrunc[:15] + "..."
					}
					line := fmt.Sprintf("%-18s | %-18s | %-20s", siteTrunc, userTrunc, strings.Repeat("•", 20))
					if i == m.Cursor {
						table.WriteString(selectedStyle.Render("▶ " + line))
					} else {
						table.WriteString("  " + line)
					}
					table.WriteString("\n")
				}
			}
		}

		s.WriteString(boxStyle.Width(70).Render(table.String()))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("Press [↑/↓] to navigate • [Enter] to copy • [Esc] to menu • [Q] to quit"))
	}

	// Add status message if present
	if m.StatusMessage != "" && time.Now().Before(m.StatusExpiry) {
		s.WriteString("\n")
		if strings.Contains(m.StatusMessage, "failed") || strings.Contains(m.StatusMessage, "Error") || strings.Contains(m.StatusMessage, "cannot") {
			s.WriteString(errorStyle.Render(m.StatusMessage))
		} else {
			s.WriteString(successStyle.Render(m.StatusMessage))
		}
		s.WriteString("\n")
	}

	return s.String()
}
