package browser

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path,omitempty"`
}

func NewContext(timeout int, verbose bool) (context.Context, context.CancelFunc) {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	}

	if chromePath := os.Getenv("CONVERTHTML2PDF_CHROME_PATH"); chromePath != "" {
		opts = append(opts, chromedp.ExecPath(chromePath))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(
		context.Background(),
		opts...,
	)

	var ctx context.Context
	var cancel context.CancelFunc

	if verbose {
		ctx, cancel = chromedp.NewContext(
			allocCtx,
			chromedp.WithDebugf(log.Printf),
		)
	} else {
		ctx, cancel = chromedp.NewContext(allocCtx)
	}

	ctx, timeoutCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)

	return ctx, func() {
		timeoutCancel()
		cancel()
		allocCancel()
	}
}

func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return ctx
}

func SetCookies(ctx context.Context, cookies []Cookie) error {
	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, cookie := range cookies {
				expr := `document.cookie = "` + cookie.Name + `=` + cookie.Value
				if cookie.Domain != "" {
					expr += `; domain=` + cookie.Domain
				}
				if cookie.Path != "" {
					expr += `; path=` + cookie.Path
				}
				expr += `"`

				if err := chromedp.Evaluate(expr, nil).Do(ctx); err != nil {
					return err
				}
			}
			return nil
		}),
	)
}