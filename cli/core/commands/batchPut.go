package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/cli/utils"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
)

var BatchPutCmd = &cobra.Command{
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
