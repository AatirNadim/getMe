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
)

var getJsonOutputPath string

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

			// Find target command
			var targetCmd *cobra.Command
			for _, c := range rootCmd.Commands() {
				if c.Name() == cmdName {
					targetCmd = c
					break
				}
			}

			if targetCmd == nil || targetCmd.Name() == "getMe_repl" || cmdName == "help" {
				if cmdName != "help" {
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

			// Special handling for getJson output flag in REPL
			// To keep it simple, we just parse it if it exists inside cmdArgs
			if cmdName == "getJson" {
				var finalArgs []string
				var outPath string
				for i := 0; i < len(cmdArgs); i++ {
					if (cmdArgs[i] == "-o" || cmdArgs[i] == "--out") && i+1 < len(cmdArgs) {
						outPath = cmdArgs[i+1]
						i++ // skip next arg
					} else {
						finalArgs = append(finalArgs, cmdArgs[i])
					}
				}
				getJsonOutputPath = outPath
				cmdArgs = finalArgs
			}

			// We bypass standard flag parsing and args validation in Cobra REPL
			// because Cobra is not designed for continuous re-execution on the same tree.
			// We will manually validate Args if possible.
			if targetCmd.Args != nil {
				if err := targetCmd.Args(targetCmd, cmdArgs); err != nil {
					fmt.Println("Error:", err)
					continue
				}
			}

			if targetCmd.RunE != nil {
				if err := targetCmd.RunE(targetCmd, cmdArgs); err != nil {
					fmt.Println("Error:", err)
				}
			} else if targetCmd.Run != nil {
				targetCmd.Run(targetCmd, cmdArgs)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commands.GetCmd)
	commands.GetJsonCmd.Flags().StringVarP(&getJsonOutputPath, "out", "o", "", "Optional path to write JSON value to")
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
