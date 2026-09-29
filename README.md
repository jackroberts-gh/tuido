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

Colors adapt to your system theme, and follow it live - switch between light and dark mode
while Tuido is running and the palette updates to match, no restart needed.

## Keyboard Shortcuts

**Navigation**
- `↑`/`k` - Move up
- `↓`/`j` - Move down

**Task status**

Each task moves through three states: todo, in-progress, then done.

- `→`/`l` - Advance status (todo to in-progress to done)
- `←`/`h` - Undo status (done to in-progress to todo)

**Tasks**
- `a` - Add task
- `e` - Edit task
- `d` - Delete task
- `sp` - Sort by priority (3 modes - highest, lowest, unsorted)
- `sd` - Sort by date (3 modes - soonest, latest, unsorted)

**View**
- `t` - Toggle completed tasks
- `?` - Help
- `q`/`Ctrl+C` - Quit


## License

MIT
