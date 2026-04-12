package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AatirNadim/getMe/cli/core"
	"github.com/AatirNadim/getMe/cli/core/commands"
	"github.com/AatirNadim/getMe/cli/core/service"
	"github.com/AatirNadim/getMe/cli/utils"
	"github.com/AatirNadim/getMe/utils/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "getMe",
	Short: "A simple file-based key-value store.",
	Long: `getMe is a CLI application that provides a persistent key-value store
backed by an append-only log on your local disk.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		httpClient, err := core.CreateHttpClient(utils.SocketPath)

		serviceLayer := &service.ServiceLayer{
			HttpClient: httpClient,
		}

		logger.Info("HTTP client created with socket path:", utils.SocketPath)

		logger.Info("Http client set as context to the command")
		ctx := context.WithValue(cmd.Context(), "serviceLayer", serviceLayer)
		cmd.SetContext(ctx)

		return err
	},
}

var replCmd = &cobra.Command{
	Use:   "getMe_repl",
	Short: "Start an interactive REPL session for getMe",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Welcome to the getMe REPL. Type 'exit' or 'quit' to close.")

		for {
			fmt.Print("getMe> ")
			if !scanner.Scan() {
				break
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if line == "exit" || line == "quit" {
				break
			}

			inputArgs := utils.ParseCommandLine(line)
			if len(inputArgs) == 0 {
				continue
			}

			cmdName := inputArgs[0]
			cmdArgs := inputArgs[1:]

			// to clear the interactive terminal screen
			if cmdName == "clear" || cmdName == "cls" {
				fmt.Print("\033[H\033[2J")
				fmt.Printf("getMe> \n")
				continue
			}

			// Find target command
			var targetCmd *cobra.Command
			for _, c := range rootCmd.Commands() {
				if c.Name() == cmdName {
					targetCmd = c
					break
				}
			}

			if targetCmd == nil || targetCmd.Name() == "getMe_repl" || cmdName == "help" {
				if targetCmd != nil && targetCmd.Name() == "getMe_repl" {
					fmt.Println("you sly mofo, you just tried to run REPL inside REPL. that's cute..!")
				} else if cmdName != "help" {
					fmt.Printf("Unknown command: %s\n", cmdName)
				}
				// print available commands
				fmt.Println("Available commands:")
				for _, c := range rootCmd.Commands() {
					if c.Name() != "getMe_repl" {
						fmt.Printf("  %-15s %s\n", c.Name(), c.Short)
					}
				}
				continue
			}

			// Set the context from our repl command to the target command
			targetCmd.SetContext(cmd.Context())

			// Reset all flags to default values to avoid state bleeding across REPL iterations
			targetCmd.Flags().VisitAll(func(f *pflag.Flag) {
				_ = f.Value.Set(f.DefValue)
				f.Changed = false
			})

			// Parse flags out of the raw input
			if err := targetCmd.ParseFlags(cmdArgs); err != nil {
				if err == pflag.ErrHelp {
					targetCmd.Help()
				} else {
					fmt.Println("Error parsing flags:", err)
				}
				continue
			}

			// Retrieve the positional arguments (without the flags)
			cleanArgs := targetCmd.Flags().Args()

			// We bypass standard flag parsing and args validation in Cobra REPL
			// because Cobra is not designed for continuous re-execution on the same tree.
			// We will manually validate Args if possible.
			if targetCmd.Args != nil {
				if err := targetCmd.Args(targetCmd, cleanArgs); err != nil {
					fmt.Println("Error:", err)
					continue
				}
			}

			if targetCmd.RunE != nil {
				if err := targetCmd.RunE(targetCmd, cleanArgs); err != nil {
					fmt.Println("Error:", err)
				}
			} else if targetCmd.Run != nil {
				targetCmd.Run(targetCmd, cleanArgs)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commands.GetCmd)
	rootCmd.AddCommand(commands.GetJsonCmd)
	rootCmd.AddCommand(commands.BatchGetCmd)
	rootCmd.AddCommand(commands.PutCmd)
	rootCmd.AddCommand(commands.PutJsonCmd)
	rootCmd.AddCommand(commands.BatchPutCmd)
	rootCmd.AddCommand(commands.DeleteCmd)
	rootCmd.AddCommand(commands.BatchDeleteCmd)
	rootCmd.AddCommand(commands.ClearCmd)
	rootCmd.AddCommand(replCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
