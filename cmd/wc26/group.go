package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var groupCmd = &cobra.Command{
	Use:   "group <letter>",
	Short: "Get a single group with standings",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		resp, err := cl.GetGroup(args[0])
		if err != nil {
			return err
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		switch outputFmt {
		case "json":
			return printJSON(resp)
		case "plain":
			fmt.Printf("Group %s\n\n", resp.Group.Name)
			for _, t := range resp.Group.Teams {
				fmt.Printf("%s\tPts:%s\tMP:%s\tW:%s\tD:%s\tL:%s\tGF:%s\tGA:%s\tGD:%s\n",
					t.TeamID, t.Pts, t.MP, t.W, t.D, t.L, t.GF, t.GA, t.GD)
			}
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintf(w, "=== Group %s ===\n", resp.Group.Name)
			fmt.Fprintln(w, "Team ID\tPts\tMP\tW\tD\tL\tGF\tGA\tGD")
			fmt.Fprintln(w, "--------\t---\t--\t-\t-\t-\t--\t--\t--")
			for _, t := range resp.Group.Teams {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					t.TeamID, t.Pts, t.MP, t.W, t.D, t.L, t.GF, t.GA, t.GD)
			}
			w.Flush()
		}
		return nil
	},
}
