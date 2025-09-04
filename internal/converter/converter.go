package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/aarontaylor/converthtml2pdf/internal/browser"
)

type Converter struct {
	options Options
}

func New(options Options) *Converter {
	return &Converter{
		options: options,
	}
}

func (c *Converter) Convert(url, outputPath string, verbose, quiet bool) error {
	ctx, cancel := browser.NewContext(c.options.Timeout, verbose)
	defer cancel()

	if c.options.UserAgent != "" {
		ctx = browser.WithUserAgent(ctx, c.options.UserAgent)
	}

	if c.options.CookiesFile != "" {
		cookies, err := c.loadCookies()
		if err != nil {
			return &FileError{Message: fmt.Sprintf("failed to load cookies: %v", err)}
		}
		if err := browser.SetCookies(ctx, cookies); err != nil {
			return &BrowserError{Message: fmt.Sprintf("failed to set cookies: %v", err)}
		}
	}

	var buf []byte
	err := chromedp.Run(ctx,
		c.navigateAndWait(url, verbose, quiet),
		c.printToPDF(&buf),
	)

	if err != nil {
		if err == context.DeadlineExceeded {
			return &ConversionError{Message: fmt.Sprintf("timeout after %d seconds", c.options.Timeout)}
		}
		return &ConversionError{Message: fmt.Sprintf("conversion failed: %v", err)}
	}

	if err := ioutil.WriteFile(outputPath, buf, 0644); err != nil {
		return &FileError{Message: fmt.Sprintf("failed to write PDF: %v", err)}
	}

	return nil
}

func (c *Converter) navigateAndWait(url string, verbose, quiet bool) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		if !quiet {
			fmt.Printf("✓ Loading page: %s\n", url)
		}

		if err := chromedp.Navigate(url).Do(ctx); err != nil {
			return err
		}

		if c.options.WaitSelector != "" {
			if verbose {
				log.Printf("Waiting for selector: %s", c.options.WaitSelector)
			}
			if err := chromedp.WaitVisible(c.options.WaitSelector).Do(ctx); err != nil {
				return fmt.Errorf("selector '%s' not found", c.options.WaitSelector)
			}
		} else {
			if err := chromedp.WaitReady("body").Do(ctx); err != nil {
				return err
			}
		}

		if c.options.WaitTime > 0 {
			if verbose {
				log.Printf("Waiting %d seconds for page to settle", c.options.WaitTime)
			}
			time.Sleep(time.Duration(c.options.WaitTime) * time.Second)
		}

		if !quiet {
			fmt.Println("✓ Generating PDF...")
		}

		return nil
	}
}

func (c *Converter) printToPDF(buf *[]byte) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		margins := c.options.ParseMargins()
		
		params := page.PrintToPDF()
		params.Landscape = c.options.Landscape
		params.PrintBackground = c.options.PrintBackground
		params.Scale = c.options.Scale
		params.MarginTop = margins.Top
		params.MarginBottom = margins.Bottom
		params.MarginLeft = margins.Left
		params.MarginRight = margins.Right

		switch c.options.Format {
		case "A3":
			params.PaperWidth = 11.7
			params.PaperHeight = 16.5
		case "Letter":
			params.PaperWidth = 8.5
			params.PaperHeight = 11
		case "Legal":
			params.PaperWidth = 8.5
			params.PaperHeight = 14
		case "Tabloid":
			params.PaperWidth = 11
			params.PaperHeight = 17
		default:
			params.PaperWidth = 8.27
			params.PaperHeight = 11.7
		}

		if c.options.HeaderTemplate != "" {
			params.DisplayHeaderFooter = true
			params.HeaderTemplate = c.options.HeaderTemplate
		}
		if c.options.FooterTemplate != "" {
			params.DisplayHeaderFooter = true
			params.FooterTemplate = c.options.FooterTemplate
		}

		data, _, err := params.Do(ctx)
		if err != nil {
			return err
		}
		*buf = data
		return nil
	})
}

func (c *Converter) loadCookies() ([]browser.Cookie, error) {
	data, err := ioutil.ReadFile(c.options.CookiesFile)
	if err != nil {
		return nil, err
	}

	var cookies []browser.Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, err
	}

	return cookies, nil
}