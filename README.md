# Password Generator CLI

A modern, secure password generator built with Go and Bubbletea. Generate strong passwords, save them with associated usernames, and manage your saved passwords with an intuitive terminal interface.

## Features

- **Secure Password Generation**: Creates cryptographically secure random passwords
- **Customizable Length**: Set password length from 4 to 128 characters
- **Save with Metadata**: Store passwords with site/service name and username
- **Password Management**: View, filter, and copy saved passwords
- **Modern UI**: Beautiful terminal interface with colors and styling
- **Clipboard Integration**: Automatic clipboard copying for convenience

## Usage

### Basic Usage

1. **Welcome Screen**: Choose what you want to do from the main menu
2. **Generate Password**: Select 'G' to create a new password (then set length)
3. **View Saved Passwords**: Select 'L' to browse your saved passwords
4. **Settings**: Select 'S' to configure default settings
5. **Save Passwords**: After generating, use 'S' to save with site name and username

### Key Bindings

#### Welcome Screen
- `G` - Generate new password
- `L` - View saved passwords
- `S` - Open settings
- `Q` - Quit application

#### Settings View
- `Enter` - Generate password with specified length
- `Esc` - Return to main menu
- `Q` - Quit

#### Main View (Password Display)
- `R` - Refresh/Generate new password
- `C` - Copy current password to clipboard
- `S` - Save current password
- `L` - List saved passwords
- `Esc` - Return to main menu
- `Q` - Quit application

#### Save View
- `Tab` - Switch between site name and username fields
- `Enter` - Save password
- `Esc` - Cancel save

#### List View
- `↑/↓` - Navigate through saved passwords
- `Enter` - Copy selected password to clipboard
- `Esc` - Return to main menu
- Type to filter passwords by site or username

## Data Storage

Passwords are stored in `passwords.csv` in the following format:
```
site_name,username,password
example.com,user@example.com,generated_password
```

## Installation

```bash
go build -o passwordgen ./cmd/passwordgen
```

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - Bubbletea components
- [Clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard access
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## Project Structure

```
my-cli-app/
├── cmd/passwordgen/          # Main application entry point
├── internal/
│   ├── app/                  # TUI application logic
│   │   ├── model.go         # Application state and model
│   │   ├── update.go        # Event handling and updates
│   │   └── view.go          # UI rendering
│   └── password/            # Password generation and storage
│       ├── generator.go     # Password generation functions
│       └── csv.go           # CSV file operations
├── go.mod
├── go.sum
├── README.md
├── CONTRIBUTING.md
└── LICENSE
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Development setup
- Code style guidelines
- Pull request process
- Testing requirements

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
