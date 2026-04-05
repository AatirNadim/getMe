package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/AatirNadim/getMe/utils/logger"
)

// methods valid for the scope of our CLI interactions with the getMe service
var validMethods = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodPut:    true,
	http.MethodDelete: true,
}

// requestOptions holds all the possible configurations for creating an HTTP request
type RequestOptions struct {
	Method      string
	URL         string
	Path        string
	Body        io.Reader
	Headers     map[string]string // Optional
	QueryParams map[string]string // Optional
	PathParams  map[string]string // Optional
}

// createHTTPRequest is a universal function to build an HTTP request
func CreateHTTPRequest(opts RequestOptions) (*http.Request, error) {
	method := strings.ToUpper(opts.Method)

	// method has to be one of the valid HTTP methods we support for our CLI interactions
	if !validMethods[method] {
		return nil, fmt.Errorf("invalid HTTP method: %s", method)
	}

	// handle path params by replacing {key} or :key with the actual value securely
	finalPath := opts.Path
	if len(opts.PathParams) > 0 {
		for key, val := range opts.PathParams {
			escapedVal := url.PathEscape(val)
			finalPath = strings.ReplaceAll(finalPath, "{"+key+"}", escapedVal)
			finalPath = strings.ReplaceAll(finalPath, ":"+key, escapedVal)
		}
	}

	// construct the full URL seamlessly
	// handling cases where BaseURL or Path may or may not have slashes to avoid malformed URLs
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

	// adding optioanl headers
	if len(opts.Headers) > 0 {
		for k, v := range opts.Headers {
			req.Header.Add(k, v)
		}
	}

	return req, nil
}

func exetuteHTTPRequestUtil(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)

	logger.Debug("Received response from server for request:", resp)

	if err != nil {
		logger.Error("Error occurred while making request:", err)
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
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

	body, err := exetuteHTTPRequestUtil(client, req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ExecuteHTTPRequestAndReturnBuffer(client *http.Client, req *http.Request) ([]byte, error) {
	return exetuteHTTPRequestUtil(client, req)
}

func ValidateJSONAndFilePath(jsonFilePath string) error {
	info, err := os.Stat(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to stat JSON file '%s': %w", jsonFilePath, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("JSON path '%s' is not a regular file", jsonFilePath)
	}
	if info.Size() == 0 {
		return fmt.Errorf("JSON file '%s' is empty", jsonFilePath)
	}
	if info.Size() > MaxJSONFileSizeBytes {
		return fmt.Errorf("JSON file '%s' size %d bytes exceeds the limit of %d bytes", jsonFilePath, info.Size(), MaxJSONFileSizeBytes)
	}

	return nil
}

func GetStringFromJSONFile(jsonFilePath string) (string, error) {
	err := ValidateJSONAndFilePath(jsonFilePath)
	if err != nil {
		return "", fmt.Errorf("JSON file validation failed for file '%s': %w", jsonFilePath, err)
	}

	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read JSON file '%s': %w", jsonFilePath, err)
	}

	if !json.Valid(fileContent) {
		return "", fmt.Errorf("file '%s' does not contain valid JSON", jsonFilePath)
	}

	var compacted bytes.Buffer
	if err := json.Compact(&compacted, fileContent); err != nil {
		return "", fmt.Errorf("failed to compact JSON from file '%s': %w", jsonFilePath, err)
	}
	value := compacted.String()
	logger.Info("Compacted JSON value: ", value)

	return value, nil
}

func StoreJSONInFile(data []byte, outputPath string) error {
	if !json.Valid(data) {
		return fmt.Errorf("data is not valid JSON, cannot store in file '%s'", outputPath)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		return fmt.Errorf("failed to format JSON data for file '%s': %w", outputPath, err)
	}

	if err := os.WriteFile(outputPath, pretty.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write JSON to file '%s': %w", outputPath, err)
	}
	fmt.Println("JSON value written to", outputPath)
	return nil

}
