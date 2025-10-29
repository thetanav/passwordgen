package app

import (
	"fmt"
	"strings"
	"time"
)

func (m Model) View() string {
	var s strings.Builder
	s.WriteString("\n")

	// Helper function to create header with title and status
	createHeader := func(title string) string {
		if m.StatusMessage != "" && time.Now().Before(m.StatusExpiry) {
			// Use normal text styling (no color, normal weight)
			styledMsg := "[" + m.StatusMessage + "]"
			// Create a line with title left-aligned and status right-aligned
			line := fmt.Sprintf("%-50s%s", title, styledMsg)
			return line
		}
		return title
	}

	switch m.CurrentView {
	case ViewWelcome:
		s.WriteString(titleStyle.Render(createHeader("Secure Password Manager")))
		s.WriteString("\n\n")

		options := []string{"1. Generate New Password", "2. View Saved Passwords", "3. Settings", "4. Quit Application"}
		for i, option := range options {
			if i == m.MenuCursor {
				s.WriteString(selectedStyle.Render("> " + option))
			} else {
				s.WriteString("  " + option)
			}
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("[↑/↓] - navigate • [Enter] to select"))
		s.WriteString("\n")

	case ViewSave:
		s.WriteString(titleStyle.Render(createHeader("Save Password")))
		s.WriteString("\n\n")
		s.WriteString(m.SiteInput.View())
		s.WriteString("\n")
		s.WriteString(m.UsernameInput.View())
		s.WriteString("\n\n")
		s.WriteString(successStyle.Render("Password copied to clipboard"))
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("[Tab] - switch fields • [Enter] - save"))

	case ViewMain:
		s.WriteString(titleStyle.Render(createHeader("Secure Password Generator")))
		s.WriteString("\n\n")

		if m.Password != "" {
			s.WriteString(fmt.Sprintf("Generated Password: %s", m.Password))
			s.WriteString("\n\n")
			s.WriteString(fmt.Sprintf("Length: %d characters", m.Length))
			s.WriteString("\n\n")
		}

		s.WriteString(infoStyle.Render("[R] - refresh • [C] - copy • [S] - save & copy"))

	case ViewList:
		s.WriteString(titleStyle.Render(createHeader("Saved Passwords")))
		s.WriteString("\n\n")
		s.WriteString(m.FilterInput.View())
		s.WriteString("\n\n")

		var table strings.Builder

		// Headers
		table.WriteString(infoStyle.Render(fmt.Sprintf("  %-18s   %-18s   %-20s", "Site", "Username", "Password")))
		table.WriteString("\n")

		// Filter passwords
		filtered := m.filterPasswords()

		if len(filtered) == 0 {
			table.WriteString(infoStyle.Render("\n  No passwords found."))
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
					line := fmt.Sprintf("%-18s   %-18s   %-5s", siteTrunc, userTrunc, strings.Repeat("•", 5))
					if i == m.Cursor {
						table.WriteString(selectedStyle.Render("> " + line))
					} else {
						table.WriteString("  " + line)
					}
					table.WriteString("\n")
				}
			}
		}

		s.WriteString(boxStyle.Width(70).Render(table.String()))
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("[↑/↓] - navigate • [Enter] - copy"))

	case ViewSettings:
		s.WriteString(titleStyle.Render(createHeader("Password Settings")))
		s.WriteString("\n\n")
		s.WriteString(m.LengthInput.View())
		s.WriteString("\n\n")

		options := []string{
			fmt.Sprintf("[%s] Include Lowercase (a-z)", checkbox(m.IncludeLower)),
			fmt.Sprintf("[%s] Include Uppercase (A-Z)", checkbox(m.IncludeUpper)),
			fmt.Sprintf("[%s] Include Numbers (0-9)", checkbox(m.IncludeNumbers)),
			fmt.Sprintf("[%s] Include Symbols (!@#$%%^&*)", checkbox(m.IncludeSymbols)),
		}

		for i, option := range options {
			if i == m.SettingsCursor {
				s.WriteString(selectedStyle.Render("> " + option))
			} else {
				s.WriteString("  " + option)
			}
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("[↑/↓] - navigate • [Space] - toggle • [Enter] - save"))
	}

	return s.String()
}
