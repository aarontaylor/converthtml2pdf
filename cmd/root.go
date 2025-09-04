package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aarontaylor/converthtml2pdf/internal/converter"
	"github.com/aarontaylor/converthtml2pdf/internal/utils"
)

const version = "1.0.0"

type Config struct {
	URL      string
	Output   string
	Options  converter.Options
	Verbose  bool
	Quiet    bool
	ShowHelp bool
	Version  bool
}

func Execute() error {
	config := parseFlags()

	if config.Version {
		fmt.Printf("converthtml2pdf version %s\n", version)
		return nil
	}

	if config.ShowHelp || len(flag.Args()) < 2 {
		printHelp()
		if !config.ShowHelp {
			return fmt.Errorf("missing required arguments")
		}
		return nil
	}

	config.URL = flag.Arg(0)
	config.Output = flag.Arg(1)

	if !strings.HasSuffix(config.Output, ".pdf") {
		config.Output += ".pdf"
	}

	if err := utils.ValidateURL(config.URL); err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	outputDir := filepath.Dir(config.Output)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	conv := converter.New(config.Options)
	
	if !config.Quiet {
		fmt.Println("✓ Initializing browser context...")
	}

	if err := conv.Convert(config.URL, config.Output, config.Verbose, config.Quiet); err != nil {
		return handleError(err)
	}

	if !config.Quiet {
		fileInfo, _ := os.Stat(config.Output)
		size := float64(fileInfo.Size()) / (1024 * 1024)
		fmt.Printf("✓ PDF saved successfully: %s (%.1f MB)\n", config.Output, size)
	}

	return nil
}

func parseFlags() *Config {
	config := &Config{
		Options: converter.DefaultOptions(),
	}

	flag.StringVar(&config.Options.Format, "format", "A4", "Page format (A4, A3, Letter, Legal, Tabloid)")
	flag.StringVar(&config.Options.Format, "f", "A4", "Page format (shorthand)")
	
	flag.BoolVar(&config.Options.Landscape, "landscape", false, "Use landscape orientation")
	flag.BoolVar(&config.Options.Landscape, "l", false, "Use landscape orientation (shorthand)")
	
	flag.StringVar(&config.Options.Margin, "margin", "0.4", "Page margins in inches")
	flag.StringVar(&config.Options.Margin, "m", "0.4", "Page margins in inches (shorthand)")
	
	flag.Float64Var(&config.Options.Scale, "scale", 1.0, "Scale of webpage rendering (0.1-2.0)")
	flag.Float64Var(&config.Options.Scale, "s", 1.0, "Scale of webpage rendering (shorthand)")
	
	flag.IntVar(&config.Options.WaitTime, "wait", 2, "Wait time in seconds after page load")
	flag.IntVar(&config.Options.WaitTime, "w", 2, "Wait time in seconds (shorthand)")
	
	flag.IntVar(&config.Options.Timeout, "timeout", 30, "Maximum time in seconds for page load")
	flag.IntVar(&config.Options.Timeout, "t", 30, "Maximum time in seconds (shorthand)")
	
	flag.BoolVar(&config.Options.PrintBackground, "background", true, "Print background graphics")
	flag.BoolVar(&config.Options.PrintBackground, "b", true, "Print background graphics (shorthand)")
	
	flag.StringVar(&config.Options.HeaderTemplate, "header", "", "Custom header HTML")
	flag.StringVar(&config.Options.FooterTemplate, "footer", "", "Custom footer HTML")
	flag.StringVar(&config.Options.WaitSelector, "selector", "", "CSS selector to wait for")
	flag.StringVar(&config.Options.CookiesFile, "cookies", "", "Path to JSON file with cookies")
	flag.StringVar(&config.Options.UserAgent, "user-agent", "", "Custom user agent string")
	flag.StringVar(&config.Options.UserAgent, "u", "", "Custom user agent (shorthand)")
	
	flag.IntVar(&config.Options.ViewportWidth, "width", 1920, "Viewport width in pixels")
	flag.IntVar(&config.Options.ViewportHeight, "height", 1080, "Viewport height in pixels")
	
	flag.BoolVar(&config.Options.NoImages, "no-images", false, "Skip loading images")
	flag.BoolVar(&config.Options.Grayscale, "grayscale", false, "Convert to grayscale PDF")
	
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&config.Quiet, "quiet", false, "Suppress all output except errors")
	flag.BoolVar(&config.Quiet, "q", false, "Suppress all output (shorthand)")
	
	flag.BoolVar(&config.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&config.ShowHelp, "h", false, "Show help message (shorthand)")
	
	flag.BoolVar(&config.Version, "version", false, "Show version information")
	flag.BoolVar(&config.Version, "v", false, "Show version (shorthand)")

	flag.Parse()
	
	return config
}

func printHelp() {
	fmt.Println("converthtml2pdf - Convert HTML web pages to PDF")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  converthtml2pdf [flags] <url> <output.pdf>")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  converthtml2pdf https://example.com output.pdf")
	fmt.Println("  converthtml2pdf -l -f Letter https://example.com report.pdf")
	fmt.Println("  converthtml2pdf --wait 5 --selector \"#content\" https://app.com doc.pdf")
	fmt.Println()
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func handleError(err error) error {
	switch err.(type) {
	case *converter.URLError:
		os.Exit(2)
	case *converter.ConversionError:
		os.Exit(3)
	case *converter.FileError:
		os.Exit(4)
	case *converter.BrowserError:
		os.Exit(5)
	}
	return err
}