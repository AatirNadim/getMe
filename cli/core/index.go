package core

import (
	"fmt"
	"getMeMod/server/store"
	"getMeMod/utils/logger"
	"getMeMod/server/store/utils/constants"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var storeInstance *store.Store

var rootCmd = &cobra.Command{
	Use:   "getMe",
	Short: "A simple file-based key-value store.",
	Long: `getMe is a CLI application that provides a persistent key-value store
backed by an append-only log on your local disk.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (*http.Client, error) {
		// Determine default store path in user's home directory: ~/.getMeStore
		return createHttpClient(constants.SocketPath)
	},
}

var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a value by its key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// key := args[0]
		// value, found, err := storeInstance.Get(key)
		// if err != nil {
		// 	return fmt.Errorf("error getting value for key '%s': %w", key, err)
		// }
		// if !found {
		// 	return fmt.Errorf("key '%s' not found", key)
		// }
		// fmt.Println(value)

		logger.Info("Get command is not implemented yet")

		return nil
	},
}

var putCmd = &cobra.Command{
	Use:   "put [key] [value]",
	Short: "Put a key-value pair into the store",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// key := args[0]
		// value := args[1]
		// if err := storeInstance.Put(key, value); err != nil {
		// 	return fmt.Errorf("error putting value for key '%s': %w", key, err)
		// }
		// fmt.Printf("Successfully set value for key '%s'\n", key)

		logger.Info("Put command is not implemented yet")
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a key-value pair from the store",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// key := args[0]
		// if err := storeInstance.Delete(key); err != nil {
		// 	return fmt.Errorf("error deleting key '%s': %w", key, err)
		// }
		// fmt.Printf("Successfully deleted key '%s'\n", key)

		logger.Info("Delete command is not implemented yet")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(putCmd)
	rootCmd.AddCommand(deleteCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}