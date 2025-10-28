# Contributing to Secure Password Manager

Thank you for your interest in contributing to the Secure Password Manager! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.19 or later
- Git

### Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/my-cli-app.git`
3. Change to the project directory: `cd my-cli-app`
4. Install dependencies: `go mod download`
5. Build the project: `go build ./cmd/passwordgen`

### Project Structure

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

## Code Style

This project follows standard Go conventions:

- Use `gofmt` to format your code
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused on a single responsibility

### Linting

Run the following to check code quality:

```bash
go vet ./...
gofmt -d .
```

## Making Changes

1. Create a new branch for your feature/fix: `git checkout -b feature/your-feature-name`
2. Make your changes
3. Test your changes thoroughly
4. Ensure all tests pass (if applicable)
5. Format your code: `gofmt -w .`
6. Commit your changes with a clear, descriptive message

### Commit Messages

Follow conventional commit format:

- `feat: add new feature`
- `fix: resolve bug`
- `docs: update documentation`
- `style: format code`
- `refactor: restructure code`

## Pull Request Process

1. Ensure your branch is up to date with main
2. Push your branch to your fork
3. Create a Pull Request with a clear title and description
4. Reference any related issues
5. Wait for review and address any feedback

### PR Requirements

- All CI checks must pass
- Code is properly formatted
- Changes are tested
- Documentation is updated if needed
- No breaking changes without discussion

## Testing

Currently, this project doesn't have automated tests. When contributing:

- Test your changes manually
- Ensure the application builds without errors
- Verify functionality works as expected

Future contributions should include tests for new features.

## Reporting Issues

When reporting bugs:

- Use a clear, descriptive title
- Provide steps to reproduce
- Include your Go version and OS
- Attach relevant logs or screenshots

## Code of Conduct

This project follows a code of conduct to ensure a welcoming environment for all contributors. Be respectful and constructive in all interactions.

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project (MIT).