package commands

import (
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
)

var ClearCmd = &cobra.Command{
	Use:   "clearStoreConfirm",
	Short: "Clear all key-value pairs from the store",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		err := serviceLayer.ClearStoreService()

		logger.Debug("Performed clear store request")

		if err != nil {
			logger.Error("Error occurred while making clear store request:", err)
			return fmt.Errorf("failed to clear the store: %w", err)
		}

		fmt.Println("Store successfully cleared")

		return nil
	},
}
