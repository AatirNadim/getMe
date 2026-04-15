package commands

import (
	"encoding/json"
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/cli/utils"
	logger "github.com/AatirNadim/getMe/utils"
	"github.com/spf13/cobra"
)

var BatchGetCmd = &cobra.Command{
	Use:   "batchGet [jsonFilePath]",
	Short: "Batch get values for multiple keys specified in either a JSON file or via --data flag",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dataFlag, _ := cmd.Flags().GetString("data")

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		var respStr string

		if dataFlag != "" {
			// Validate that the provided data is valid JSON
			if !json.Valid([]byte(dataFlag)) {
				return fmt.Errorf("provided data flag is not valid JSON")
			}
			var err error
			respStr, err = serviceLayer.BatchGetServiceFromData([]byte(dataFlag))
			if err != nil {
				return fmt.Errorf("failed to perform batch get from data: %w", err)
			}
		} else {
			if len(args) == 0 {
				return fmt.Errorf("must provide either a JSON file path as a positional argument or the --data flag")
			}
			jsonFilePath := args[0]
			err := utils.ValidateJSONAndFilePath(jsonFilePath)
			if err != nil {
				return fmt.Errorf("failed to validate JSON file '%s': %w", jsonFilePath, err)
			}
			var serviceErr error
			respStr, serviceErr = serviceLayer.BatchGetService(jsonFilePath)
			if serviceErr != nil {
				return fmt.Errorf("failed to perform batch get: %w", serviceErr)
			}
		}

		outPath, _ := cmd.Flags().GetString("out")
		if outPath != "" {
			err := utils.StoreJSONInFile([]byte(respStr), outPath)
			if err != nil {
				logger.Error("Error occurred while storing JSON value in file:", err)
			}
		}

		fmt.Println(respStr)

		return nil
	},
}

func init() {
	BatchGetCmd.Flags().StringP("out", "o", "", "Optional path to write JSON value to")
	BatchGetCmd.Flags().StringP("data", "d", "", "JSON data payload")
}
