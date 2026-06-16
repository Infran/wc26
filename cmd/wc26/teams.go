package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var teamGroup string

var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "List all teams (optional: filter by group)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		resp, err := cl.GetTeams(teamGroup)
		if err != nil {
			return err
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		switch outputFmt {
		case "json":
			return printJSON(resp.Teams)
		case "plain":
			for _, t := range resp.Teams {
				fmt.Printf("%s\t%s\t%s\t%s\n", t.ID, t.NameEN, t.FifaCode, t.Groups)
			}
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tName\tFIFA Code\tGroup\tFlag")
			fmt.Fprintln(w, "--\t----\t---------\t-----\t----")
			for _, t := range resp.Teams {
				flag := t.Flag
				if flag == "" {
					flag = "-"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t.ID, t.NameEN, t.FifaCode, t.Groups, flag)
			}
			w.Flush()
		}
		return nil
	},
}

func init() {
	teamsCmd.Flags().StringVarP(&teamGroup, "group", "g", "", "Filter by group letter (A-L)")
}

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON marshal: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func jsonMarshalOrPanic(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func printSingleJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON marshal: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
