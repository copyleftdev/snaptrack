package capture

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func CaptureHTML(url string, timeout time.Duration) (string,
	map[string][]string,
	map[string][]string,
	int,
	error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", nil, nil, 0, fmt.Errorf("creating request for %s: %w", url, err)
	}

	req.Header.Set("User-Agent", "Snaptrack/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, nil, 0, fmt.Errorf("http get error for %s: %w", url, err)
	}
	defer resp.Body.Close()

	finalReqHeaders := cloneHeader(req.Header)

	finalRespHeaders := cloneHeader(resp.Header)
	statusCode := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, 0, fmt.Errorf("reading body from %s: %w", url, err)
	}
	return string(body), finalReqHeaders, finalRespHeaders, statusCode, nil
}

func cloneHeader(h http.Header) map[string][]string {
	m := make(map[string][]string)
	for k, vals := range h {
		tmp := make([]string, len(vals))
		copy(tmp, vals)
		m[k] = tmp
	}
	return m
}
