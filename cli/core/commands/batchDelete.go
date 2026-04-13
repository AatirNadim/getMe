package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AatirNadim/getMe/cli/core/service"

	"github.com/AatirNadim/getMe/cli/utils"
	"github.com/AatirNadim/getMe/commons"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
)

var BatchDeleteCmd = &cobra.Command{
	Use:   "batchDelete [jsonFilePath]",
	Short: "Batch delete multiple keys specified in either a JSON file or via --data flag",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dataFlag, _ := cmd.Flags().GetString("data")

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)
		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		var batchDeleteReq commons.BatchGetRequestBody
		var sourceName string

		if dataFlag != "" {
			if !json.Valid([]byte(dataFlag)) {
				return fmt.Errorf("provided data flag is not valid JSON")
			}
			if err := json.Unmarshal([]byte(dataFlag), &batchDeleteReq); err != nil {
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

			if err := json.Unmarshal(fileContent, &batchDeleteReq); err != nil {
				return fmt.Errorf("failed to parse JSON file '%s': %w", jsonFilePath, err)
			}
			sourceName = fmt.Sprintf("file '%s'", jsonFilePath)
		}

		err := serviceLayer.BatchDeleteService(batchDeleteReq)

		logger.Debug("Performed batch delete request for keys from:", sourceName)

		if err != nil {
			logger.Error("Error occurred while making batch DELETE request:", err)
			return fmt.Errorf("failed to perform batch delete for keys from %s: %w", sourceName, err)
		}

		fmt.Println("Batch delete operation successful for keys from:", sourceName)

		return nil
	},
}

func init() {
	BatchDeleteCmd.Flags().StringP("data", "d", "", "JSON data payload")
}
