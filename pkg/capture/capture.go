package capture

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CaptureHTML preserves the old signature but internally uses a basic HTTP GET
// (no JavaScript execution). That way, any code referencing capture.CaptureHTML
// continues to work with minimal changes.
func CaptureHTML(url string, timeout time.Duration) (string, error) {
	// Create a timed context so we don’t hang indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request for %s: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http get error for %s: %w", url, err)
	}
	defer resp.Body.Close()

	// (Optional) handle non-200 statuses:
	// if resp.StatusCode != http.StatusOK {
	//     return "", fmt.Errorf("got status %d for %s", resp.StatusCode, url)
	// }

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body from %s: %w", url, err)
	}
	return string(body), nil
}
