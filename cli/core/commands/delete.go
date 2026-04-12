package commands

import (
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a key-value pair from the store",
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

		err := serviceLayer.DeleteService(key)

		logger.Debug("Performed delete request for key:", key)

		if err != nil {
			logger.Error("Error occurred while making DELETE request:", err)
			return fmt.Errorf("failed to delete key '%s': %w", key, err)
		}

		fmt.Println("Key-value pair successfully deleted from the store")

		return nil
	},
}
