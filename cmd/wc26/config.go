package main

import (
	"fmt"
	"os"

	"github.com/Infran/wc26/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.InitConfig(); err != nil {
			return fmt.Errorf("initializing config: %w", err)
		}
		cfgPath, _ := config.ConfigFilePath()
		fmt.Fprintf(os.Stderr, "Config initialized at: %s\n", cfgPath)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current config values",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := config.Show()
		if err != nil {
			return err
		}
		fmt.Print(data)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value (e.g. api.base_url, auth.token, output)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Set(args[0], args[1]); err != nil {
			return fmt.Errorf("setting config: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Config updated: %s = %s\n", args[0], args[1])
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show config file location",
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := config.ConfigFilePath()
		if err != nil {
			return err
		}
		fmt.Println(p)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPathCmd)
}
