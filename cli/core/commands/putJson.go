package commands

import (
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/cli/utils"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
)

var PutJsonCmd = &cobra.Command{
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
