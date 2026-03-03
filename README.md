# Tuido

A beautiful, interactive terminal TODO list manager built with Go and BubbleTea.

## Features

- **Interactive TUI**: Full-featured terminal user interface with keyboard navigation
- **Task Management**: Add, delete, and toggle task completion
- **Priority Levels**: Organize tasks with high, medium, and low priorities
- **Due Dates**: Set deadlines with flexible date formats
- **Persistent Storage**: Tasks stored locally in `~/.tuido/tasks.json`
- **Color-coded Display**: Visual indicators for priorities and due dates

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/jackroberts-gh/tuido.git
cd tuido

# Build the binary
go build -o tuido main.go

# Optionally, install to your $GOPATH/bin
go install
```

## Usage

Simply run the `tuido` command to launch the interactive interface:

```bash
./tuido
```

Or if installed:

```bash
tuido
```

## Keyboard Shortcuts

### Navigation
- `↑` or `k` - Move cursor up
- `↓` or `j` - Move cursor down

### Task Operations
- `space` - Toggle task completion
- `a` - Add new task
- `d` - Delete selected task
- `p` - Change priority (1=Low, 2=Medium, 3=High)
- `e` - Edit due date

### View Options
- `t` - Toggle show/hide completed tasks
- `?` - Show help screen
- `q` or `Ctrl+C` - Quit application

## Due Date Formats

Tuido supports flexible date input formats:

- **Relative**: `3d` (3 days), `1w` (1 week), `2m` (2 months)
- **Absolute**: `2026-03-15`, `03/15/2026`
- **Natural**: `today`, `tomorrow`, `next week`
- **Short**: `Jan 15`, `01/15`

To remove a due date, simply press Enter with an empty input.

## Priority System

Tasks can be assigned three priority levels:

- **High** (Red ■) - Urgent and important tasks
- **Medium** (Yellow ■) - Default priority for new tasks
- **Low** (Green ■) - Less urgent tasks

## Data Storage

All tasks are stored in `~/.tuido/tasks.json` in a human-readable JSON format. You can manually edit this file if needed, though it's recommended to use the TUI interface.

### Example Storage Format

```json
{
  "tasks": [
    {
      "id": "uuid-here",
      "text": "Complete project documentation",
      "completed": false,
      "priority": 2,
      "due_date": "2026-03-10T00:00:00Z",
      "created_at": "2026-03-03T10:00:00Z"
    }
  ]
}
```

## Project Structure

```
tuido/
├── main.go                    # Application entry point
├── internal/
│   ├── model/                 # Data models
│   │   ├── task.go           # Task structure
│   │   └── tasklist.go       # Task list operations
│   ├── storage/              # Data persistence
│   │   └── json.go           # JSON file I/O
│   ├── tui/                  # Terminal UI
│   │   ├── model.go          # BubbleTea model
│   │   ├── update.go         # Event handling
│   │   ├── view.go           # Rendering
│   │   └── styles.go         # UI styling
│   └── config/               # Configuration
│       └── config.go         # Data directory paths
└── README.md
```

## Dependencies

- [BubbleTea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal output
- [UUID](https://github.com/google/uuid) - Unique task identifiers

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o tuido main.go
```

### Running in Development

```bash
go run main.go
```

## License

MIT License - feel free to use and modify as needed.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Future Enhancements

Potential features for future versions:
- Task categories and tags
- Search and filter functionality
- Task notes and descriptions
- Recurring tasks
- Export to Markdown/CSV
- Multiple task lists
- Custom themes
- Configurable keyboard shortcuts

## Author

Built with ❤️ using Go and BubbleTea
