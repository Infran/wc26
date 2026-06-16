package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Infran/wc26/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	matchType    string
	matchTeam    string
	matchDate    string
	matchMatchday string
)

var matchesCmd = &cobra.Command{
	Use:   "matches",
	Short: "List matches (optional: filter by type, team, date, matchday)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		resp, err := cl.GetGames()
		if err != nil {
			return err
		}

		games := resp.Games

		if matchType != "" {
			var filtered []api.Game
			t := strings.ToLower(matchType)
			for _, g := range games {
				if strings.ToLower(g.Type) == t {
					filtered = append(filtered, g)
				}
			}
			games = filtered
		}
		if matchTeam != "" {
			var filtered []api.Game
			t := strings.ToLower(matchTeam)
			for _, g := range games {
				if strings.Contains(strings.ToLower(g.HomeTeamNameEN), t) ||
					strings.Contains(strings.ToLower(g.AwayTeamNameEN), t) ||
					strings.Contains(strings.ToLower(g.HomeTeamLabel), t) ||
					strings.Contains(strings.ToLower(g.AwayTeamLabel), t) {
					filtered = append(filtered, g)
				}
			}
			games = filtered
		}
		if matchDate != "" {
			var filtered []api.Game
			for _, g := range games {
				if strings.Contains(g.LocalDate, matchDate) {
					filtered = append(filtered, g)
				}
			}
			games = filtered
		}
		if matchMatchday != "" {
			var filtered []api.Game
			for _, g := range games {
				if g.Matchday == matchMatchday {
					filtered = append(filtered, g)
				}
			}
			games = filtered
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		switch outputFmt {
		case "json":
			return printJSON(games)
		case "plain":
			for _, g := range games {
				home := g.HomeTeamNameEN
				if home == "" {
					home = g.HomeTeamLabel
				}
				away := g.AwayTeamNameEN
				if away == "" {
					away = g.AwayTeamLabel
				}
				status := g.TimeElapsed
				if g.Finished == "TRUE" {
					status = "FT"
				}
				fmt.Printf("%s | %s vs %s | %s-%s | %s | %s\n",
					g.ID, home, away, g.HomeScore, g.AwayScore, g.LocalDate, status)
			}
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tHome\tScore\tAway\tGroup\tDate\tStatus")
			fmt.Fprintln(w, "--\t----\t-----\t----\t-----\t----\t------")
			for _, g := range games {
				home := g.HomeTeamNameEN
				if home == "" {
					home = g.HomeTeamLabel
				}
				away := g.AwayTeamNameEN
				if away == "" {
					away = g.AwayTeamLabel
				}
				score := g.HomeScore + "-" + g.AwayScore
				status := g.TimeElapsed
				if g.Finished == "TRUE" {
					status = "FT"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					g.ID, home, score, away, g.Group, g.LocalDate, status)
			}
			w.Flush()
		}
		return nil
	},
}

func init() {
	matchesCmd.Flags().StringVarP(&matchType, "type", "t", "", "Filter by stage (group, r32, r16, qf, sf, third, final)")
	matchesCmd.Flags().StringVarP(&matchTeam, "team", "", "", "Filter by team name")
	matchesCmd.Flags().StringVarP(&matchDate, "date", "d", "", "Filter by date (e.g. 06/11/2026)")
	matchesCmd.Flags().StringVarP(&matchMatchday, "matchday", "m", "", "Filter by matchday (1-9)")
}
