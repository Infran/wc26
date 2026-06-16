package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups with standings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		resp, err := cl.GetGroups()
		if err != nil {
			return err
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		switch outputFmt {
		case "json":
			return printJSON(resp.Groups)
		case "plain":
			for _, g := range resp.Groups {
				fmt.Printf("Group %s\n", g.Name)
				for _, t := range g.Teams {
					fmt.Printf("  %s\t%d pts\n", t.TeamID, t.Pts)
				}
				fmt.Println()
			}
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			for _, g := range resp.Groups {
				fmt.Fprintf(w, "=== Group %s ===\n", g.Name)
				fmt.Fprintln(w, "Team ID\tPts\tMP\tW\tD\tL\tGF\tGA\tGD")
				fmt.Fprintln(w, "--------\t---\t--\t-\t-\t-\t--\t--\t--")
				for _, t := range g.Teams {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
						t.TeamID, t.Pts, t.MP, t.W, t.D, t.L, t.GF, t.GA, t.GD)
				}
				fmt.Fprintln(w)
			}
			w.Flush()
		}
		return nil
	},
}
