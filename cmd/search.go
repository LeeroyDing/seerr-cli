/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/LeeroyDing/seerr-cli/pkg/seerr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for movies or TV shows",
	Long: `Search for media items on your Seerr instance by title or keyword.
Displays a list of results with their IDs, which can be used for 'info' or 'request' commands.

Example:
  seerr search "Inception"
  seerr search "The Boys"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			fmt.Println("Use 'seerr config --url <url> --api-key <key>' to set them.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Printf("Searching for '%s'...\n\n", query)

		resp, err := client.Search(query)
		if err != nil {
			fmt.Printf("Error searching: %v\n", err)
			return
		}

		if len(resp.Results) == 0 {
			fmt.Println("No results found.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tTITLE/NAME\tRELEASE DATE")
		fmt.Fprintln(w, "--\t----\t----------\t------------")

		for _, res := range resp.Results {
			title := res.Title
			if title == "" {
				title = res.Name
			}
			date := res.ReleaseDate
			if date == "" {
				date = res.FirstAirDate
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", res.ID, res.MediaType, title, date)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
