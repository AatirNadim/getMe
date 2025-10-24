package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"getMeMod/cli/core"
	"getMeMod/cli/utils"
	"getMeMod/utils/logger"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

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
			return fmt.Errorf("server returned non-OK status: %s", resp.Status)
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
			return fmt.Errorf("server returned non-OK status: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

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
			return fmt.Errorf("server returned non-OK status: %s", resp.Status)
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
			return fmt.Errorf("server returned non-OK status: %s", resp.Status)
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
	rootCmd.AddCommand(putCmd)
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
