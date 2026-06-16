package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check API server health",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		health, err := cl.Health()
		if err != nil {
			return err
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		switch outputFmt {
		case "json":
			return printJSON(health)
		case "plain":
			fmt.Printf("Status:\t%s\n", health.Status)
			fmt.Printf("Version:\t%s\n", health.Version)
			fmt.Printf("Environment:\t%s\n", health.Environment)
			fmt.Printf("Uptime:\t%d s\n", health.Uptime)
			fmt.Printf("DB Status:\t%s\n", health.Database.Status)
			fmt.Printf("DB Name:\t%s\n", health.Database.Name)
			fmt.Printf("Memory Used:\t%s\n", health.Memory.Used)
			fmt.Printf("Memory Total:\t%s\n", health.Memory.Total)
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "Status:\t%s\n", health.Status)
			fmt.Fprintf(w, "Version:\t%s\n", health.Version)
			fmt.Fprintf(w, "Environment:\t%s\n", health.Environment)
			fmt.Fprintf(w, "Uptime:\t%d s\n", health.Uptime)
			fmt.Fprintf(w, "Database Status:\t%s\n", health.Database.Status)
			fmt.Fprintf(w, "Database Name:\t%s\n", health.Database.Name)
			fmt.Fprintf(w, "Memory Used:\t%s\n", health.Memory.Used)
			fmt.Fprintf(w, "Memory Total:\t%s\n", health.Memory.Total)
			w.Flush()
		}
		return nil
	},
}
