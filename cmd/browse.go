/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/LeeroyDing/seerr-cli/pkg/seerr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse [trending|movies|tv]",
	Short: "Browse trending or popular media",
	Long: `Discover new content by browsing trending items or popular movies and TV shows.
	
Example:
  seerr browse trending
  seerr browse movies`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		category := args[0]
		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Printf("Browsing %s...\n\n", category)

		var resp *seerr.SearchResponse
		var err error

		switch category {
		case "trending":
			resp, err = client.GetTrending()
		case "movies":
			resp, err = client.GetPopularMovies()
		case "tv":
			resp, err = client.GetPopularTV()
		default:
			fmt.Printf("Error: invalid category '%s'. Use trending, movies, or tv.\n", category)
			return
		}

		if err != nil {
			fmt.Printf("Error fetching content: %v\n", err)
			return
		}

		if len(resp.Results) == 0 {
			fmt.Println("No items found.")
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
	rootCmd.AddCommand(browseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// browseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// browseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
