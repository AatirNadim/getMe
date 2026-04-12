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

var BatchDeleteCmd = &cobra.Command{
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
