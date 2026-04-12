package commands

import (
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"

	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
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

		resp, err := serviceLayer.GetService(key)
		if err != nil {
			return fmt.Errorf("failed to get value for key '%s': %w", key, err)
		}

		fmt.Println(resp)

		return nil
	},
}
