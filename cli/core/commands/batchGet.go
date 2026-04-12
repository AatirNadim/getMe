package commands

import (
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"

	"github.com/AatirNadim/getMe/cli/utils"
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

		fmt.Println(respStr)

		return nil
	},
}
