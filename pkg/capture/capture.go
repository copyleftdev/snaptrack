// Package capture provides functions to capture DOM content or other assets from a webpage
// using a headless Chrome (via chromedp).
package capture

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

// Browser holds a chromedp context, allowing you to reuse it for multiple captures if desired.
type Browser struct {
	Ctx       context.Context
	Cancel    context.CancelFunc
	allocated bool
}

// NewBrowser creates a new browser context that can be used for multiple captures.
// If you'd rather create a fresh Chrome instance for each capture, you can skip this
// and just use CaptureHTML directly.
func NewBrowser() (*Browser, error) {
	// Create a root context; cancelFunc frees resources when done.
	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		// Below are typical flags to enable headless mode, disable GPU, etc.
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		// You can add more flags if necessary...
	)

	// Create a new browser context from the allocator context
	browserCtx, browserCancel := chromedp.NewContext(ctx)
	b := &Browser{
		Ctx: browserCtx,
		Cancel: func() {
			// Cancel both contexts when shutting down
			browserCancel()
			cancel()
		},
		allocated: true,
	}

	// You can run a quick check to ensure Chrome starts
	if err := chromedp.Run(browserCtx); err != nil {
		b.Cancel()
		return nil, fmt.Errorf("failed to start chromedp: %w", err)
	}
	return b, nil
}

// CaptureHTML retrieves the full HTML content of the given URL.
// By default, this function creates a *temporary* browser context that times out
// after the specified duration. If you want to reuse a browser (e.g. for multiple URLs),
// pass a shared context from a Browser instance created by NewBrowser().
func CaptureHTML(url string, timeout time.Duration) (string, error) {
	// Create a timed context so we don't hang indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a fresh context from the parent
	ctx, cancelBrowser := chromedp.NewContext(ctx)
	defer cancelBrowser()

	var html string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.NodeVisible),
		chromedp.OuterHTML("html", &html),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return "", fmt.Errorf("chromedp run error for %s: %w", url, err)
	}
	return html, nil
}

// CaptureHTMLWithBrowser uses an existing Browser context to navigate to the URL
// and retrieve the HTML. This is helpful if you want to batch multiple captures
// without spawning a fresh headless instance each time.
func CaptureHTMLWithBrowser(b *Browser, url string, timeout time.Duration) (string, error) {
	if b == nil || !b.allocated {
		return "", fmt.Errorf("invalid or uninitialized browser context")
	}

	// Create a timed child context from the existing browser root context
	ctx, cancel := context.WithTimeout(b.Ctx, timeout)
	defer cancel()

	var html string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.NodeVisible),
		chromedp.OuterHTML("html", &html),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return "", fmt.Errorf("chromedp run error for %s: %w", url, err)
	}
	return html, nil
}

// Example of a specialized capture function, e.g., simulating a specific device or screen size.
func CaptureHTMLAsMobile(b *Browser, url string, timeout time.Duration) (string, error) {
	if b == nil || !b.allocated {
		return "", fmt.Errorf("invalid or uninitialized browser context")
	}

	// iPhone 12 viewport example
	mobileMetrics := &emulation.SetDeviceMetricsOverrideParams{
		Width:             390,
		Height:            844,
		DeviceScaleFactor: 3.0,
		Mobile:            true,
	}

	ctx, cancel := context.WithTimeout(b.Ctx, timeout)
	defer cancel()

	var html string
	tasks := chromedp.Tasks{
		mobileMetrics,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.NodeVisible),
		chromedp.OuterHTML("html", &html),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return "", fmt.Errorf("chromedp run error for %s: %w", url, err)
	}
	return html, nil
}
