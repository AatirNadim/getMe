package commands

import (
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/cli/utils"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
)

var GetJsonCmd = &cobra.Command{
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

		outPath, _ := cmd.Flags().GetString("out")
		if outPath != "" {
			err := utils.StoreJSONInFile(resp, outPath)
			if err != nil {
				logger.Error("Error occurred while storing JSON value in file:", err)
			}
		}

		fmt.Println(string(resp))
		return nil
	},
}

func init() {
	GetJsonCmd.Flags().StringP("out", "o", "", "Optional path to write JSON value to")
}
