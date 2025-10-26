package main

import (
    "fmt"
    "math/big"
    "os"
    "crypto/rand"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    password string
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

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {

        case "ctrl+c", "q":
            return m, tea.Quit

        case "r":
            newPass, err := generatePassword(16)
            if err == nil {
                m.password = newPass
            } else {
                m.password = "Error generating password"
            }

        case "c":
            // TODO: Implement copy to clipboard (requires external library like github.com/atotto/clipboard)
            // For now, just print it
            fmt.Println("Password (for manual copy):", m.password)
        }
    }

    return m, nil
}

func (m model) View() string {
    // The header
    s := "Password Generator\n\n"

    // Iterate over our choices
    s += "The password is: `"
    s += m.password
    s += "`\n"

    // The footer
    s += "\nHelp q to quit, r to refresh, c to copy.\n"

    // Send the UI for rendering
    return s
}

func main() {
    initialPass, err := generatePassword(16)
    if err != nil {
        initialPass = "Error generating initial password"
    }
    m := model{
        password: initialPass,
    }
    p := tea.NewProgram(m)
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}