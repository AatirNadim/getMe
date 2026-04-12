package commands

import (
	"fmt"

	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
)

var PutCmd = &cobra.Command{
	Use:   "put [key] [value]",
	Short: "Put a key-value pair into the store",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		key := args[0]
		value := args[1]

		if key == "" {
			return fmt.Errorf("invalid key: %s", key)
		}
		if value == "" {
			return fmt.Errorf("invalid value: %s", value)
		}

		serviceLayer, ok := cmd.Context().Value("serviceLayer").(*service.ServiceLayer)

		if !ok {
			return fmt.Errorf("service layer not found in context")
		}

		err := serviceLayer.PutService(key, value)

		logger.Debug("Performed put request for key:", key, " and value:", value)

		if err != nil {
			logger.Error("Error occurred while making PUT request:", err)
			return fmt.Errorf("failed to put key '%s': %w", key, err)
		}

		fmt.Println("Key-value pair successfully put into the store")

		return nil
	},
}
