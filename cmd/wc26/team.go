package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Infran/wc26/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var teamCmd = &cobra.Command{
	Use:   "team <id|name>",
	Short: "Get a single team by ID or name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		arg := args[0]

		var team *api.Team
		var err error

		team, err = cl.GetTeamByID(arg)
		if err != nil {
			team, err = cl.GetTeamByName(arg)
			if err != nil {
				return fmt.Errorf("team not found: %w", err)
			}
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		switch outputFmt {
		case "json":
			return printJSON(team)
		case "plain":
			fmt.Printf("ID:\t%s\nName (EN):\t%s\nName (FA):\t%s\nFIFA Code:\t%s\nISO2:\t%s\nGroup:\t%s\nFlag:\t%s\n",
				team.ID, team.NameEN, team.NameFA, team.FifaCode, team.Iso2, team.Groups, team.Flag)
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "ID:\t%s\n", team.ID)
			fmt.Fprintf(w, "Name (EN):\t%s\n", team.NameEN)
			fmt.Fprintf(w, "Name (FA):\t%s\n", team.NameFA)
			fmt.Fprintf(w, "FIFA Code:\t%s\n", team.FifaCode)
			fmt.Fprintf(w, "ISO2:\t%s\n", team.Iso2)
			fmt.Fprintf(w, "Group:\t%s\n", team.Groups)
			fmt.Fprintf(w, "Flag:\t%s\n", team.Flag)
			w.Flush()
		}
		return nil
	},
}
