package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"getMeMod/cli/core"
	"getMeMod/server/store/utils/constants"
	"getMeMod/utils/logger"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)


type PutRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}


var rootCmd = &cobra.Command{
	Use:   "getMe",
	Short: "A simple file-based key-value store.",
	Long: `getMe is a CLI application that provides a persistent key-value store
backed by an append-only log on your local disk.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Determine default store path in user's home directory: ~/.getMeStore
		httpClient, err := core.CreateHttpClient(constants.SocketPath)

		logger.Info("HTTP client created with socket path:", constants.SocketPath)

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

		req, err := http.NewRequest("GET", "http://unix/get", nil)
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

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned non-OK status: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		fmt.Println(string(body))
		httpClient.Get("")

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
		jsonPayload, err := json.Marshal(PutRequestBody{
			Key:   key,
			Value: value,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal JSON payload: %w", err)
		}


		logger.Debug("preparing io reader payload with:", jsonPayload)
		readerPayload := bytes.NewReader(jsonPayload)


		req, err := http.NewRequest("POST", "http://unix/put", readerPayload)
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

		logger.Info("Put command is not implemented yet")
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

		req, err := http.NewRequest("DELETE", "http://unix/delete", nil)
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

		logger.Info("Delete command is not implemented yet")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(putCmd)
	rootCmd.AddCommand(deleteCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
