# sprite

A tool for combining multiple SVG files into a single sprite sheet.

## Features

- Combine multiple SVG files into one sprite sheet
- Automatically generate 1x, 2x, and 3x resolution PNG files
- Generate JSON map files with position and size information for each sprite
- High-quality SVG → PNG conversion using resvg
- Efficient bin packing using Shelf First Fit Decreasing Height (FFDH) algorithm

## Prerequisites

### Installing resvg

This project uses [resvg](https://github.com/RazrFalcon/resvg) to convert SVG files to PNG.

#### Option 1: Install via Cargo (Recommended)

```bash
# If Rust is not installed
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Install resvg
cargo install resvg
```

#### Option 2: Download Pre-built Binary

Download the appropriate binary for your operating system from the [resvg releases](https://github.com/RazrFalcon/resvg/releases) page and add it to your PATH.

#### Option 3: Use Makefile

```bash
make install-deps
```

## Installation

```bash
# Install dependencies and build
make all

# Or manually
go build -o sprite .
```

## Usage

```bash
./sprite <svg-file-1> <svg-file-2> ...
```

### Example

```bash
./sprite testdata/*.svg
```

Generated files:
- `sprites.svg` - Combined SVG file
- `sprites.png` - 1x PNG sprite sheet
- `sprites@2x.png` - 2x PNG sprite sheet
- `sprites@3x.png` - 3x PNG sprite sheet
- `sprites.json` - 1x sprite map
- `sprites@2x.json` - 2x sprite map
- `sprites@3x.json` - 3x sprite map

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Run with test data

```bash
make run
```

### Clean

```bash
make clean
```

## License

MIT
