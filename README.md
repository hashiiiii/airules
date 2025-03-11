# airules

`airules` is a command-line tool for installing and managing rules-for-ai files for AI-powered editors like Cursor and Windsurf.

## Features

- Install local and global rules-for-ai files for Cursor and Windsurf editors
- Install rule sets from remote repositories like awesome-cursorrules
- Support for both English and Japanese templates
- Interactive mode for selecting rule sets from remote repositories

## Installation

### Prerequisites

- Go 1.24 or later

### Building from source

1. Clone the repository:

```bash
git clone https://github.com/hashiiiii/airules.git
cd airules
```

2. Build the project:

```bash
make build
```

3. (Optional) Install the binary to your PATH:

```bash
make install
```

## Usage

### Installing Cursor rules

```bash
# Install local Cursor rules
airules cursor --type local

# Install global Cursor rules
airules cursor --type global

# Install both local and global Cursor rules
airules cursor --type all

# Install Japanese templates
airules cursor --language ja
```

### Installing Windsurf rules

```bash
# Install local Windsurf rules
airules windsurf --type local

# Install global Windsurf rules
airules windsurf --type global

# Install both local and global Windsurf rules
airules windsurf --type all

# Install Japanese templates
airules windsurf --language ja
```

### Installing rule sets from remote repositories

```bash
# List available rule sets
airules remote --list

# Install a specific rule set locally
airules remote --install <rule-set-name> --type local

# Install a specific rule set globally
airules remote --install <rule-set-name> --type global

# Interactive mode (no flags)
airules remote
```

### Displaying version information

```bash
airules version
```

## Development

### Running tests

```bash
make test
```

### Building for different platforms

```bash
make build-all
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [awesome-cursorrules](https://github.com/PatrickJS/awesome-cursorrules) - A collection of Cursor rules
