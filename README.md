# Tuido

Terminal TODO list manager.

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
go build -o tuido main.go
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

## Due Dates

Supports relative (`3d`, `1w`, `2m`), absolute (`2026-03-15`), natural (`today`, `tomorrow`), and short (`Jan 15`) formats. Press Enter with empty input to remove a due date.

## License

MIT
