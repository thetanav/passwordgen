package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    password    string
    lengthInput textinput.Model
    view        string // "settings" or "main"
    length      int    // Cached length for refreshes
}

const (
    lowercase = "abcdefghijklmnopqrstuvwxyz"
    uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    numbers   = "0123456789"
    symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

func generatePassword(length int) (string, error) {
    if length < 1 {
        return "", fmt.Errorf("length must be at least 1")
    }

    var charset strings.Builder
    charset.WriteString(lowercase)
    charset.WriteString(uppercase)
    charset.WriteString(numbers)
    charset.WriteString(symbols)

    charsetStr := charset.String()

    password := make([]rune, length)
    for i := range password {
        n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charsetStr))))
        if err != nil {
            return "", err
        }
        password[i] = rune(charsetStr[n.Int64()])
    }

    return string(password), nil
}

func initialModel() model {
    ti := textinput.New()
    ti.Placeholder = "len"
    ti.Focus()
    ti.Prompt = "Length: "
    ti.CharLimit = 3 // Keep it short

    return model{
        lengthInput: ti,
        view:        "settings",
        length:      16, // Default
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
            if m.view == "settings" {
                inputLen, err := strconv.Atoi(m.lengthInput.Value())
                if err != nil || inputLen < 1 {
                    inputLen = 16 // Default on invalid
                }
                m.length = inputLen
                newPass, err := generatePassword(m.length)
                if err != nil {
                    m.password = "Error generating password"
                } else {
                    m.password = newPass
                }
                m.view = "main"
                m.lengthInput.Blur() // Unfocus
                return m, nil
            }

        case "esc":
            if m.view == "settings" {
                // Use default and go to main
                newPass, _ := generatePassword(16)
                m.password = newPass
                m.view = "main"
                m.lengthInput.Blur()
                return m, nil
            }

        case "r":
            if m.view == "main" {
                newPass, err := generatePassword(m.length)
                if err != nil {
                    m.password = "Error generating password"
                } else {
                    m.password = newPass
                }
            }

        case "c":
            if m.view == "main" && m.password != "" {
                err := clipboard.WriteAll(m.password)
                if err != nil {
                    panic(err)
                }
            }
        }
    }

    // Delegate to textinput if in settings view
    if m.view == "settings" {
        m.lengthInput, cmd = m.lengthInput.Update(msg)
        return m, cmd
    }

    return m, cmd
} // <- CRITICAL: This closing brace ends Update() – ensure it's here!

func (m model) View() string {
    var s string

    if m.view == "settings" {
        s = "Password Generator - Set Length\n\n"
        s += m.lengthInput.View() + "\n\n"
        s += "Press Enter to generate, Esc for default (16)."
        return s
    }

    // Main view
    s = "Password Generator\n\n"
    s += "The password is: `"
    s += m.password
    s += "`\n\n"
    s += fmt.Sprintf("Current length: %d\n", m.length)

    s += "\nHelp: q to quit, r to refresh, c to copy.\n"

    return s
}

func main() {
    m := initialModel()
    p := tea.NewProgram(m)
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}