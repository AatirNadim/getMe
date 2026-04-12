package commands

import (
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"

	"github.com/AatirNadim/getMe/cli/utils"
	"github.com/AatirNadim/getMe/utils/logger"
	"github.com/spf13/cobra"
)

var BatchGetCmd = &cobra.Command{
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
}
