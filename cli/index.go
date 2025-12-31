package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AatirNadim/getMe/cli/core"
	"github.com/AatirNadim/getMe/cli/utils"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
)

var getJsonOutputPath string

var rootCmd = &cobra.Command{
	Use:   "getMe",
	Short: "A simple file-based key-value store.",
	Long: `getMe is a CLI application that provides a persistent key-value store
backed by an append-only log on your local disk.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		httpClient, err := core.CreateHttpClient(utils.SocketPath)

		logger.Info("HTTP client created with socket path:", utils.SocketPath)

		logger.Info("Http client set as context to the command")
		ctx := context.WithValue(cmd.Context(), "httpClientKey", httpClient)
		cmd.SetContext(ctx)

		return err
	},
}

var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a value by its key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		key := args[0]

		if key == "" {
			return fmt.Errorf("invalid key: %s", key)
		}
		httpClient, ok := cmd.Context().Value("httpClientKey").(*http.Client)

		if !ok {
			return fmt.Errorf("http client not found in context")
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", utils.BaseUrl, utils.GetRoute), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		q := req.URL.Query()
		q.Add("key", key)
		req.URL.RawQuery = q.Encode()

		logger.Info("Sending GET request for key:", key, "encoded the key as request query parameter")

		resp, err := httpClient.Do(req)

		logger.Debug("Received response from server for GET request for key:", resp)

		if err != nil {
			logger.Error("Error occurred while making GET request:", err)
			return fmt.Errorf("failed to get key '%s': %w", key, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned non-OK status: %s, \nbody: %v", resp.Status, string(body))
		}

		fmt.Println(string(body))

		return nil
	},
}

var getJsonCmd = &cobra.Command{
	Use:   "getJson [key]",
	Short: "Get a JSON value by its key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		key := args[0]

		if key == "" {
			return fmt.Errorf("invalid key: %s", key)
		}
		httpClient, ok := cmd.Context().Value("httpClientKey").(*http.Client)

		if !ok {
			return fmt.Errorf("http client not found in context")
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", utils.BaseUrl, utils.GetRoute), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		q := req.URL.Query()
		q.Add("key", key)
		req.URL.RawQuery = q.Encode()

		logger.Info("Sending GET (JSON) request for key:", key, "encoded the key as request query parameter")

		resp, err := httpClient.Do(req)
		if err != nil {
			logger.Error("Error occurred while making GET (JSON) request:", err)
			return fmt.Errorf("failed to get key '%s': %w", key, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned non-OK status: %s, \nbody: %v", resp.Status, string(body))
		}

		if !json.Valid(body) {
			return fmt.Errorf("value for key '%s' is not valid JSON; use 'get' instead", key)
		}

		var pretty bytes.Buffer
		if err := json.Indent(&pretty, body, "", "  "); err != nil {
			return fmt.Errorf("failed to pretty-print JSON for key '%s': %w", key, err)
		}

		if getJsonOutputPath != "" {
			if err := os.WriteFile(getJsonOutputPath, pretty.Bytes(), 0o644); err != nil {
				return fmt.Errorf("failed to write JSON to file '%s': %w", getJsonOutputPath, err)
			}
			fmt.Println("JSON value written to", getJsonOutputPath)
			return nil
		} else {
			fmt.Println("No output path specified, simply printing the value on screen")
		}

		fmt.Println(pretty.String())
		return nil
	},
}

var putCmd = &cobra.Command{
	Use:   "put [key] [value]",
	Short: "Put a key-value pair into the store",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		key := args[0]
		value := args[1]

		if key == "" {
			return fmt.Errorf("invalid key: %s", key)
		}
		if value == "" {
			return fmt.Errorf("invalid value: %s", value)
		}

		httpClient, ok := cmd.Context().Value("httpClientKey").(*http.Client)

		if !ok {
			return fmt.Errorf("http client not found in context")
		}

		logger.Debug("Preparing JSON payload for PUT request with key:", key, " and value:", value)
		jsonPayload, err := json.Marshal(utils.PutRequestBody{
			Key:   key,
			Value: value,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal JSON payload: %w", err)
		}

		logger.Debug("preparing io reader payload with:", jsonPayload)
		readerPayload := bytes.NewReader(jsonPayload)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", utils.BaseUrl, utils.PutRoute), readerPayload)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		logger.Info("Sending PUT request for key:", key, " and value:", value, " encoded as request query parameters")

		resp, err := httpClient.Do(req)

		logger.Debug("Received response from server for PUT request for key:", resp)

		if err != nil {
			logger.Error("Error occurred while making PUT request:", err)
			return fmt.Errorf("failed to put key '%s': %w", key, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("server returned non-OK status: ", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		fmt.Println(string(body))

		return nil
	},
}

var putJsonCmd = &cobra.Command{
	Use:   "putJson [key] [jsonFilePath]",
	Short: "Put a key with a JSON value loaded from a file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		key := args[0]
		jsonFilePath := args[1]

		if key == "" {
			return fmt.Errorf("invalid key: %s", key)
		}
		if jsonFilePath == "" {
			return fmt.Errorf("invalid JSON file path: %s", jsonFilePath)
		}

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
		if info.Size() > utils.MaxJSONFileSizeBytes {
			return fmt.Errorf("JSON file '%s' size %d bytes exceeds the limit of %d bytes", jsonFilePath, info.Size(), utils.MaxJSONFileSizeBytes)
		}

		fileContent, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to read JSON file '%s': %w", jsonFilePath, err)
		}

		if !json.Valid(fileContent) {
			return fmt.Errorf("file '%s' does not contain valid JSON", jsonFilePath)
		}

		var compacted bytes.Buffer
		if err := json.Compact(&compacted, fileContent); err != nil {
			return fmt.Errorf("failed to compact JSON from file '%s': %w", jsonFilePath, err)
		}
		value := compacted.String()
		fmt.Println("compacted json value: ", value)

		httpClient, ok := cmd.Context().Value("httpClientKey").(*http.Client)
		if !ok {
			return fmt.Errorf("http client not found in context")
		}

		logger.Debug("Preparing JSON payload for PUT (JSON) request with key:", key, " and value loaded from file:", jsonFilePath)
		jsonPayload, err := json.Marshal(utils.PutRequestBody{
			Key:   key,
			Value: value,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal JSON payload: %w", err)
		}

		readerPayload := bytes.NewReader(jsonPayload)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", utils.BaseUrl, utils.PutRoute), readerPayload)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		logger.Info("Sending PUT (JSON) request for key:", key, " with JSON value from file:", jsonFilePath)

		resp, err := httpClient.Do(req)
		if err != nil {
			logger.Error("Error occurred while making PUT (JSON) request:", err)
			return fmt.Errorf("failed to put key '%s' from JSON file '%s': %w", key, jsonFilePath, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("server returned non-OK status: ", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		fmt.Println(string(body))
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a key-value pair from the store",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		key := args[0]

		if key == "" {
			return fmt.Errorf("invalid key: %s", key)
		}
		httpClient, ok := cmd.Context().Value("httpClientKey").(*http.Client)

		if !ok {
			return fmt.Errorf("http client not found in context")
		}

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s", utils.BaseUrl, utils.DeleteRoute), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		q := req.URL.Query()
		q.Add("key", key)
		req.URL.RawQuery = q.Encode()
		logger.Info("Sending DELETE request for key:", key, " encoded as request query parameter")

		resp, err := httpClient.Do(req)
		logger.Debug("Received response from server for DELETE request for key:", resp)

		if err != nil {
			logger.Error("Error occurred while making DELETE request:", err)
			return fmt.Errorf("failed to delete key '%s': %w", key, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("server returned non-OK status:", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// fmt.Println("failed to read response body:", err)
			return fmt.Errorf("failed to read response body: %w", err)
		}

		// fmt.Println("printing the response body:")
		fmt.Println(string(body))

		return nil
	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all key-value pairs from the store",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		httpClient, ok := cmd.Context().Value("httpClientKey").(*http.Client)

		if !ok {
			return fmt.Errorf("http client not found in context")
		}

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s", utils.BaseUrl, utils.ClearStoreRoute), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		logger.Info("Sending CLEAR request to clear all key-value pairs from the store")

		resp, err := httpClient.Do(req)

		logger.Debug("Received response from server for CLEAR request:", resp)

		if err != nil {
			logger.Error("Error occurred while making CLEAR request:", err)
			return fmt.Errorf("failed to clear store: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("server returned non-OK status:", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		fmt.Println(string(body))

		return nil
	},
}

var batchPutCmd = &cobra.Command{
	Use:   "batchPut [jsonFilePath]",
	Short: "Batch put key-value pairs from a JSON file into the store",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		jsonFilePath := args[0]

		if jsonFilePath == "" {
			return fmt.Errorf("invalid file path: %s", jsonFilePath)
		}

		fileContent, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %w", jsonFilePath, err)
		}

		// being able to parse the JSON input file in the desired format is important!
		var keyValuePairs map[string]string
		if err := json.Unmarshal(fileContent, &keyValuePairs); err != nil {
			return fmt.Errorf("failed to parse JSON file '%s': %w", jsonFilePath, err)
		}

		httpClient, ok := cmd.Context().Value("httpClientKey").(*http.Client)

		if !ok {
			return fmt.Errorf("http client not found in context")
		}

		jsonPayload, err := json.Marshal(keyValuePairs)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON payload: %w", err)
		}

		readerPayload := bytes.NewReader(jsonPayload)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", utils.BaseUrl, utils.BatchPutRoute), readerPayload)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		logger.Info("Sending BATCH PUT request with data from file:", jsonFilePath)

		resp, err := httpClient.Do(req)

		logger.Debug("Received response from server for BATCH PUT request:", resp)

		if err != nil {
			logger.Error("Error occurred while making BATCH PUT request:", err)
			return fmt.Errorf("failed to perform batch put: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("server returned non-OK status:", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		fmt.Println(string(body))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getJsonCmd.Flags().StringVarP(&getJsonOutputPath, "out", "o", "", "Optional path to write JSON value to")
	rootCmd.AddCommand(getJsonCmd)
	rootCmd.AddCommand(putCmd)
	rootCmd.AddCommand(putJsonCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(batchPutCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
