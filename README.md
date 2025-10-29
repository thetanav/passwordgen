# Password Generator CLI

A modern, secure password generator built with Go and Bubbletea. Generate strong passwords, save them with associated usernames, and manage your saved passwords with an intuitive terminal interface.

## Features

- **Secure Password Generation**: Creates cryptographically secure random passwords
- **Customizable Length**: Set password length from 4 to 128 characters
- **Save with Metadata**: Store passwords with site/service name and username
- **Password Management**: View, filter, and copy saved passwords
- **Modern UI**: Beautiful terminal interface with colors and styling
- **Clipboard Integration**: Automatic clipboard copying for convenience

## Data Storage

Passwords are stored in `passwords.csv` in the following format:

```
site_name,username,password
example.com,user@example.com,generated_password
```

## Installation

### Option 1: One-liner install (recommended)

```bash
curl -sSL https://raw.githubusercontent.com/thetanav/passwordgen/main/install.sh | bash
```

This will download the source code, build the application, and install it to a directory in your PATH.

### Option 2: Using the install script locally

If you have the source code:

```bash
./install.sh
```

### Option 3: Manual build

```bash
go build -o passwordgen .
# Then move passwordgen to a directory in your PATH, e.g.:
# sudo mv passwordgen /usr/local/bin/
```

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - Bubbletea components
- [Clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard access
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## Project Structure

```
passwordgen/
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
