# Tuido

Tuido (_pronounced to-do_) is a simple, minimalist TUI designed to manage a local todo list.

## Installation

### Via Go

```bash
go install github.com/jackroberts-gh/tuido@latest
```

### Via GitHub Releases

Download the latest binary for your platform from the [releases page](https://github.com/jackroberts-gh/tuido/releases).

### Build from Source

```bash
git clone https://github.com/jackroberts-gh/tuido.git
cd tuido
go build
```

## Usage

Run `tuido` to launch the interactive interface. Tasks are stored in `~/.tuido/tasks.json`.

Colors adapt to your system theme. Restart the app if you switch between light/dark mode.

## Keyboard Shortcuts

**Navigation**
- `↑`/`k` - Move up
- `↓`/`j` - Move down

**Tasks**
- `space` - Toggle completion
- `a` - Add task
- `e` - Edit task
- `d` - Delete task
- `p` - Change priority (1=Low, 2=Medium, 3=High)
- `Shift+E` - Edit due date

**View**
- `t` - Toggle completed tasks
- `?` - Help
- `q`/`Ctrl+C` - Quit


## License

MIT
