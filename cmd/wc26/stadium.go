package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stadiumCmd = &cobra.Command{
	Use:   "stadium <id>",
	Short: "Get a single stadium by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		stadium, err := cl.GetStadiumByID(args[0])
		if err != nil {
			return err
		}

		outputFmt := viper.GetString("output")
		if outputFmt == "" {
			outputFmt = "table"
		}

		switch outputFmt {
		case "json":
			return printJSON(stadium)
		case "plain":
			fmt.Printf("ID:\t%s\n", stadium.ID)
			fmt.Printf("Name (EN):\t%s\n", stadium.NameEN)
			fmt.Printf("Name (FA):\t%s\n", stadium.NameFA)
			fmt.Printf("FIFA Name:\t%s\n", stadium.FifaName)
			fmt.Printf("City:\t%s\n", stadium.CityEN)
			fmt.Printf("Country:\t%s\n", stadium.CountryEN)
			fmt.Printf("Capacity:\t%d\n", stadium.Capacity)
		default:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "ID:\t%s\n", stadium.ID)
			fmt.Fprintf(w, "Name (EN):\t%s\n", stadium.NameEN)
			fmt.Fprintf(w, "Name (FA):\t%s\n", stadium.NameFA)
			fmt.Fprintf(w, "FIFA Name:\t%s\n", stadium.FifaName)
			fmt.Fprintf(w, "City:\t%s\n", stadium.CityEN)
			fmt.Fprintf(w, "Country:\t%s\n", stadium.CountryEN)
			fmt.Fprintf(w, "Capacity:\t%d\n", stadium.Capacity)
			w.Flush()
		}
		return nil
	},
}
