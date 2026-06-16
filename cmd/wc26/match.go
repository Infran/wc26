package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var matchCmd = &cobra.Command{
	Use:   "match <id>",
	Short: "Get a single match by ID (1-104)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		game, err := cl.GetGameByID(args[0])
		if err != nil {
			return err
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		home := game.HomeTeamNameEN
		if home == "" {
			home = game.HomeTeamLabel
		}
		away := game.AwayTeamNameEN
		if away == "" {
			away = game.AwayTeamLabel
		}

		switch outputFmt {
		case "json":
			return printJSON(game)
		case "plain":
			fmt.Printf("ID:\t%s\n", game.ID)
			fmt.Printf("Home:\t%s\n", home)
			fmt.Printf("Away:\t%s\n", away)
			fmt.Printf("Score:\t%s-%s\n", game.HomeScore, game.AwayScore)
			fmt.Printf("Home Scorers:\t%s\n", game.HomeScorers)
			fmt.Printf("Away Scorers:\t%s\n", game.AwayScorers)
			fmt.Printf("Group:\t%s\n", game.Group)
			fmt.Printf("Matchday:\t%s\n", game.Matchday)
			fmt.Printf("Date:\t%s\n", game.LocalDate)
			fmt.Printf("Stadium ID:\t%s\n", game.StadiumID)
			fmt.Printf("Finished:\t%s\n", game.Finished)
			fmt.Printf("Status:\t%s\n", game.TimeElapsed)
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "ID:\t%s\n", game.ID)
			fmt.Fprintf(w, "Home:\t%s\n", home)
			fmt.Fprintf(w, "Away:\t%s\n", away)
			fmt.Fprintf(w, "Score:\t%s-%s\n", game.HomeScore, game.AwayScore)
			fmt.Fprintf(w, "Home Scorers:\t%s\n", game.HomeScorers)
			fmt.Fprintf(w, "Away Scorers:\t%s\n", game.AwayScorers)
			fmt.Fprintf(w, "Group:\t%s\n", game.Group)
			fmt.Fprintf(w, "Matchday:\t%s\n", game.Matchday)
			fmt.Fprintf(w, "Date:\t%s\n", game.LocalDate)
			fmt.Fprintf(w, "Stadium ID:\t%s\n", game.StadiumID)
			fmt.Fprintf(w, "Finished:\t%s\n", game.Finished)
			fmt.Fprintf(w, "Status:\t%s\n", game.TimeElapsed)
			w.Flush()
		}
		return nil
	},
}
