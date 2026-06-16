package main

import (
	"fmt"
	"os"

	"github.com/Infran/wc26/internal/api"
	"github.com/Infran/wc26/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	apiURL   string
	outputFmt string
	noColor  bool
)

func newAPIClient() *api.Client {
	baseURL := viper.GetString("api.base_url")
	if apiURL != "" {
		baseURL = apiURL
	}
	cl := api.NewClient(baseURL)
	if token := viper.GetString("auth.token"); token != "" {
		cl.SetToken(token)
	}
	return cl
}

var rootCmd = &cobra.Command{
	Use:   "wc26",
	Short: "FIFA World Cup 2026 CLI — Live scores, teams, groups, matches & stadiums",
	Long: `A command-line interface for the FIFA World Cup 2026 API.

Provides access to live scores, team data, group standings, match schedules,
stadium information, and health status. Supports both the public hosted API
and self-hosted instances.

Documentation: https://worldcup26.ir/api-docs
GitHub: https://github.com/Infran/wc26`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use == "completion" {
			return nil
		}
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(func() {
		if err := config.InitConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: config init: %v\n", err)
		}
	})

	rootCmd.PersistentFlags().StringVarP(&apiURL, "api-url", "", "", "API base URL (default: https://worldcup26.ir)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "", "Output format: table, json, plain")
	rootCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "", false, "Disable colored output")

	viper.BindPFlag("api.base_url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(teamsCmd)
	rootCmd.AddCommand(teamCmd)
	rootCmd.AddCommand(groupsCmd)
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(matchesCmd)
	rootCmd.AddCommand(matchCmd)
	rootCmd.AddCommand(stadiumsCmd)
	rootCmd.AddCommand(stadiumCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(updateCmd)
}
