package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/cli/utils"
	logger "github.com/AatirNadim/getMe/utils"

	"github.com/spf13/cobra"
)

var BatchPutCmd = &cobra.Command{
	Use:   "batchPut [jsonFilePath]",
	Short: "Batch put key-value pairs from either a JSON file or via --data flag into the store",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dataFlag, _ := cmd.Flags().GetString("data")

		serviceLayer, ok := cmd.Context().Value(utils.ServiceLayerKey).(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		var keyValuePairs map[string]string
		var sourceName string

		if dataFlag != "" {
			if !json.Valid([]byte(dataFlag)) {
				return fmt.Errorf("provided data flag is not valid JSON")
			}
			if err := json.Unmarshal([]byte(dataFlag), &keyValuePairs); err != nil {
				return fmt.Errorf("failed to parse JSON data: %w", err)
			}
			sourceName = "provided data flag"
		} else {
			if len(args) == 0 {
				return fmt.Errorf("must provide either a JSON file path as a positional argument or the --data flag")
			}
			jsonFilePath := args[0]
			err := utils.ValidateJSONAndFilePath(jsonFilePath)
			if err != nil {
				return fmt.Errorf("failed to validate JSON file '%s': %w", jsonFilePath, err)
			}

			fileContent, err := os.ReadFile(jsonFilePath)
			if err != nil {
				return fmt.Errorf("failed to read file '%s': %w", jsonFilePath, err)
			}

			if err := json.Unmarshal(fileContent, &keyValuePairs); err != nil {
				return fmt.Errorf("failed to parse JSON file '%s': %w", jsonFilePath, err)
			}
			sourceName = fmt.Sprintf("file '%s'", jsonFilePath)
		}

		err := serviceLayer.BatchPutService(keyValuePairs)

		logger.Debug("Performed BATCH PUT request with data from:", sourceName)

		if err != nil {
			logger.Error("Error occurred while making BATCH PUT request:", err)
			return fmt.Errorf("failed to perform batch put with data from %s: %w", sourceName, err)
		}

		fmt.Println("Batch put operation successful for key-value pairs from:", sourceName)

		return nil
	},
}

func init() {
	BatchPutCmd.Flags().StringP("data", "d", "", "JSON data payload")
}
