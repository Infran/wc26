package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stadiumsCmd = &cobra.Command{
	Use:   "stadiums",
	Short: "List all stadiums",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		resp, err := cl.GetStadiums()
		if err != nil {
			return err
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		switch outputFmt {
		case "json":
			return printJSON(resp.Stadiums)
		case "plain":
			for _, s := range resp.Stadiums {
				fmt.Printf("%s\t%s\t%s\t%d\n", s.ID, s.NameEN, s.CityEN, s.Capacity)
			}
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tName\tCity\tCountry\tCapacity")
			fmt.Fprintln(w, "--\t----\t----\t-------\t--------")
			for _, s := range resp.Stadiums {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
					s.ID, s.NameEN, s.CityEN, s.CountryEN, s.Capacity)
			}
			w.Flush()
		}
		return nil
	},
}
