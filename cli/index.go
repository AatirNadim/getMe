package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/AatirNadim/getMe/cli/core"
	"github.com/AatirNadim/getMe/cli/core/service"
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

		serviceLayer := &service.ServiceLayer{
			HttpClient: httpClient,
		}

		logger.Info("HTTP client created with socket path:", utils.SocketPath)

		logger.Info("Http client set as context to the command")
		ctx := context.WithValue(cmd.Context(), "serviceLayer", serviceLayer)
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
		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		var resp, err = serviceLayer.GetService(key)
		if err != nil {
			return fmt.Errorf("failed to get value for key '%s': %w", key, err)
		}

		fmt.Println(resp)

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
		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		resp, err := serviceLayer.GetJsonValueService(key)

		if err != nil {
			return fmt.Errorf("failed to get JSON value for key '%s': %w", key, err)
		}

		if getJsonOutputPath != "" {
			err := utils.StoreJSONInFile(resp, getJsonOutputPath)
			if err != nil {
				logger.Error("Error occurred while storing JSON value in file:", err)
			}
		}

		fmt.Println(string(resp))
		return nil
	},
}

var batchGetCmd = &cobra.Command{
	Use:   "batchGet [jsonFilePath]",
	Short: "Batch get values for multiple keys specified in a JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonFilePath := args[0]

		err := utils.ValidateJSONAndFilePath(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to validate JSON file '%s': %w", jsonFilePath, err)
		}

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		respStr, err := serviceLayer.BatchGetService(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to perform batch get: %w", err)
		}

		fmt.Println(respStr)

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

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		err := serviceLayer.PutService(key, value)

		logger.Debug("Performed put request for key:", key, " and value:", value)

		if err != nil {
			logger.Error("Error occurred while making PUT request:", err)
			return fmt.Errorf("failed to put key '%s': %w", key, err)
		}

		fmt.Println("Key-value pair successfully put into the store")

		return nil
	},
}

var putJsonCmd = &cobra.Command{
	Use:   "putJson [key] [jsonFilePath]",
	Short: "Put a key with a JSON value loaded from a file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)
		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		key := args[0]
		jsonFilePath := args[1]

		if key == "" {
			return fmt.Errorf("invalid key: %s", key)
		}
		if jsonFilePath == "" {
			return fmt.Errorf("invalid JSON file path: %s", jsonFilePath)
		}

		value, err := utils.GetStringFromJSONFile(jsonFilePath)

		if err != nil {
			return fmt.Errorf("Failed to extract value from JSON file '%s': %w", jsonFilePath, err)
		}

		err = serviceLayer.PutService(key, value)

		logger.Debug("Performed PUT (JSON) request for key:", key, " and JSON value from file:", jsonFilePath)

		if err != nil {
			logger.Error("Error occurred while making PUT (JSON) request:", err)
			return fmt.Errorf("failed to put key '%s' with JSON value from file '%s': %w", key, jsonFilePath, err)
		}

		fmt.Println("Key and JSON value successfully put into the store")

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
		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		err := serviceLayer.DeleteService(key)

		logger.Debug("Performed delete request for key:", key)

		if err != nil {
			logger.Error("Error occurred while making DELETE request:", err)
			return fmt.Errorf("failed to delete key '%s': %w", key, err)
		}

		fmt.Println("Key-value pair successfully deleted from the store")

		return nil
	},
}

var batchDeleteCmd = &cobra.Command{
	Use:   "batchDelete [jsonFilePath]",
	Short: "Batch delete multiple keys specified in a JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonFilePath := args[0]

		err := utils.ValidateJSONAndFilePath(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to validate JSON file '%s': %w", jsonFilePath, err)
		}

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		fileContent, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %w", jsonFilePath, err)
		}

		var batchDeleteReq utils.BatchGetRequestBody
		if err := json.Unmarshal(fileContent, &batchDeleteReq); err != nil {
			return fmt.Errorf("failed to parse JSON file '%s': %w", jsonFilePath, err)
		}

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		err = serviceLayer.BatchDeleteService(batchDeleteReq)

		logger.Debug("Performed batch delete request for keys from file:", jsonFilePath)

		if err != nil {
			logger.Error("Error occurred while making batch DELETE request:", err)
			return fmt.Errorf("failed to perform batch delete for keys from file '%s': %w", jsonFilePath, err)
		}

		fmt.Println("Batch delete operation successful for keys from file:", jsonFilePath)

		return nil
	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all key-value pairs from the store",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		err := serviceLayer.ClearStoreService()

		logger.Debug("Performed clear store request")

		if err != nil {
			logger.Error("Error occurred while making clear store request:", err)
			return fmt.Errorf("failed to clear the store: %w", err)
		}

		fmt.Println("Store successfully cleared")

		return nil
	},
}

var batchPutCmd = &cobra.Command{
	Use:   "batchPut [jsonFilePath]",
	Short: "Batch put key-value pairs from a JSON file into the store",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		jsonFilePath := args[0]

		err := utils.ValidateJSONAndFilePath(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to validate JSON file '%s': %w", jsonFilePath, err)
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

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		err = serviceLayer.BatchPutService(keyValuePairs)

		logger.Debug("Performed BATCH PUT request with data from file:", jsonFilePath)

		if err != nil {
			logger.Error("Error occurred while making BATCH PUT request:", err)
			return fmt.Errorf("failed to perform batch put with data from file '%s': %w", jsonFilePath, err)
		}

		fmt.Println("Batch put operation successful for key-value pairs from file:", jsonFilePath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getJsonCmd.Flags().StringVarP(&getJsonOutputPath, "out", "o", "", "Optional path to write JSON value to")
	rootCmd.AddCommand(getJsonCmd)
	rootCmd.AddCommand(batchGetCmd)
	rootCmd.AddCommand(putCmd)
	rootCmd.AddCommand(putJsonCmd)
	rootCmd.AddCommand(batchPutCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(batchDeleteCmd)
	rootCmd.AddCommand(clearCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
