# airules

`airules` is a CLI tool for installing configuration files for AI-powered editors like Windsurf and Cursor to appropriate locations.

## Features

- Install Windsurf configuration files
- Install Cursor configuration files
- Selective installation of local and global configuration files

## Installation

### Using Go

```bash
go install github.com/hashiiiii/airules@latest
```

### Building from source

```bash
git clone https://github.com/hashiiiii/airules.git
cd airules
go build -o bin/airules ./cmd/airules
```

## Usage

### Basic commands

```bash
# Install Windsurf configuration files (both local and global)
airules windsurf

# Install Cursor configuration files (both local and global)
airules cursor

# Install only Windsurf local configuration file
airules windsurf -l
# or
airules windsurf --local

# Install only Cursor global configuration file
airules cursor -g
# or
airules cursor --global

# Display version information
airules version

# Display help
airules -h
```

## Configuration File Locations

### Windsurf

- Local configuration file: `cascade.local.json` in the current directory
- Global configuration file:
  - macOS/Linux: `~/.config/windsurf/cascade.global.json`
  - Windows: `%APPDATA%\Windsurf\cascade.global.json`

### Cursor

- Local configuration file: `prompt_library.local.json` in the current directory
- Global configuration file:
  - macOS: `~/Library/Application Support/Cursor/prompt_library.global.json`
  - Linux: `~/.config/cursor/prompt_library.global.json`
  - Windows: `%APPDATA%\Cursor\prompt_library.global.json`

## Customizing Templates

Template files are located in the `vendor/rules-for-ai` directory. You can customize the installed configurations by editing these files.

## License

MIT
