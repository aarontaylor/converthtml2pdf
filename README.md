# converthtml2pdf

[![Go Version](https://img.shields.io/badge/Go-1.24.4-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A command-line tool written in Go that converts HTML web pages to PDF files with high fidelity rendering using headless Chrome/Chromium.

## Features

- Convert any web page to PDF with a simple command
- Support for various page formats (A4, A3, Letter, Legal, Tabloid)
- Customizable margins, orientation, and scale
- Wait for specific CSS selectors or JavaScript execution
- Cookie authentication support
- Custom headers and footers
- Viewport size configuration
- Grayscale conversion option

## Installation

### Prerequisites

converthtml2pdf requires Chrome or Chromium to be installed on your system:

- **macOS**: Chrome is usually pre-installed, or download from [Google Chrome](https://www.google.com/chrome/)
- **Linux**: `sudo apt-get install chromium-browser` or `sudo yum install chromium`
- **Windows**: Download from [Google Chrome](https://www.google.com/chrome/)

### From Source

```bash
git clone https://github.com/aarontaylor/converthtml2pdf.git
cd converthtml2pdf
make build
sudo make install
```

### Using Go Install

```bash
go install github.com/aarontaylor/converthtml2pdf@latest
```

## Usage

### Basic Usage

```bash
converthtml2pdf <url> <output.pdf>
```

### Examples

Convert a simple web page:
```bash
converthtml2pdf https://example.com output.pdf
```

Convert with landscape orientation and Letter format:
```bash
converthtml2pdf -l -f Letter https://example.com report.pdf
```

Wait for specific content to load:
```bash
converthtml2pdf --wait 5 --selector "#content" https://app.com doc.pdf
```

Custom margins (in inches):
```bash
converthtml2pdf --margin "1,1,1,1" https://example.com output.pdf
```

### Available Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--format` | `-f` | Page format (A4, A3, Letter, Legal, Tabloid) | A4 |
| `--landscape` | `-l` | Use landscape orientation | false |
| `--margin` | `-m` | Page margins in inches (single value or "top,right,bottom,left") | 0.4 |
| `--scale` | `-s` | Scale of webpage rendering (0.1-2.0) | 1.0 |
| `--wait` | `-w` | Wait time in seconds after page load | 2 |
| `--timeout` | `-t` | Maximum time in seconds for page load | 30 |
| `--background` | `-b` | Print background graphics | true |
| `--header` | | Custom header HTML | |
| `--footer` | | Custom footer HTML | |
| `--selector` | | CSS selector to wait for before capturing | |
| `--cookies` | | Path to JSON file with cookies | |
| `--user-agent` | `-u` | Custom user agent string | |
| `--width` | | Viewport width in pixels | 1920 |
| `--height` | | Viewport height in pixels | 1080 |
| `--no-images` | | Skip loading images | false |
| `--grayscale` | | Convert to grayscale PDF | false |
| `--verbose` | | Enable verbose logging | false |
| `--quiet` | `-q` | Suppress all output except errors | false |
| `--help` | `-h` | Show help message | |
| `--version` | `-v` | Show version information | |

### Cookie Authentication

To use cookie authentication, create a JSON file with your cookies:

```json
[
  {
    "name": "session_id",
    "value": "abc123",
    "domain": ".example.com",
    "path": "/"
  }
]
```

Then use:
```bash
converthtml2pdf --cookies cookies.json https://app.example.com output.pdf
```

### Custom Headers and Footers

Add page numbers and custom text:
```bash
converthtml2pdf --footer "<div style='font-size: 10px; text-align: center;'>Page <span class='pageNumber'></span> of <span class='totalPages'></span></div>" https://example.com output.pdf
```

## Environment Variables

- `CONVERTHTML2PDF_CHROME_PATH`: Path to Chrome/Chromium executable
- `CONVERTHTML2PDF_TIMEOUT`: Default timeout override
- `CONVERTHTML2PDF_USER_AGENT`: Default user agent
- `CONVERTHTML2PDF_DEBUG`: Enable debug logging

## Building from Source

### Requirements

- Go 1.21 or higher
- Chrome or Chromium installed

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt

# Run linters
make lint
```

## Exit Codes

The tool uses specific exit codes for different error types:

- `0`: Success
- `1`: Invalid arguments
- `2`: URL unreachable
- `3`: Conversion failed
- `4`: File write error
- `5`: Browser error

## Troubleshooting

### Chrome Not Found

If Chrome is not found, set the path explicitly:
```bash
export CONVERTHTML2PDF_CHROME_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
converthtml2pdf https://example.com output.pdf
```

### Page Not Loading Correctly

For JavaScript-heavy sites, increase the wait time:
```bash
converthtml2pdf --wait 10 https://spa-app.com output.pdf
```

Or wait for a specific element:
```bash
converthtml2pdf --selector ".content-loaded" https://app.com output.pdf
```

### Large Pages Timing Out

Increase the timeout:
```bash
converthtml2pdf --timeout 60 https://heavy-page.com output.pdf
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with [chromedp](https://github.com/chromedp/chromedp) for Chrome automation
- Inspired by tools like Puppeteer and wkhtmltopdf
