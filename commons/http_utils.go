package commons

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/AatirNadim/getMe/utils"
)

var validMethods = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodPut:    true,
	http.MethodDelete: true,
}

type RequestOptions struct {
	Method      string
	URL         string
	Path        string
	Body        io.Reader
	Headers     map[string]string // Optional
	QueryParams map[string]string // Optional
	PathParams  map[string]string // Optional
}

func CreateHTTPRequest(opts RequestOptions) (*http.Request, error) {
	method := strings.ToUpper(opts.Method)

	if !validMethods[method] {
		return nil, fmt.Errorf("invalid HTTP method: %s", method)
	}

	finalPath := opts.Path
	if len(opts.PathParams) > 0 {
		for key, val := range opts.PathParams {
			escapedVal := url.PathEscape(val)
			finalPath = strings.ReplaceAll(finalPath, "{"+key+"}", escapedVal)
			finalPath = strings.ReplaceAll(finalPath, ":"+key, escapedVal)
		}
	}

	fullURL := opts.URL
	if finalPath != "" {
		if strings.HasSuffix(fullURL, "/") && strings.HasPrefix(finalPath, "/") {
			fullURL += finalPath[1:]
		} else if !strings.HasSuffix(fullURL, "/") && !strings.HasPrefix(finalPath, "/") {
			fullURL += "/" + finalPath
		} else {
			fullURL += finalPath
		}
	}

	req, err := http.NewRequest(method, fullURL, opts.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if len(opts.QueryParams) > 0 {
		q := req.URL.Query()
		for k, v := range opts.QueryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(opts.Headers) > 0 {
		for k, v := range opts.Headers {
			req.Header.Add(k, v)
		}
	}

	return req, nil
}

func ExecuteHTTPRequestUtil(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)

	if err != nil {
		utils.Error("Error occurred while making request: ", err)
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	utils.Debug("Received response from server for request: ", resp)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status: %s, \nbody: %v", resp.Status, string(body))
	}
	return body, nil
}

func ExecuteHTTPRequest(client *http.Client, req *http.Request) (string, error) {
	body, err := ExecuteHTTPRequestUtil(client, req)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func ExecuteHTTPRequestAndReturnBuffer(client *http.Client, req *http.Request) ([]byte, error) {
	return ExecuteHTTPRequestUtil(client, req)
}
