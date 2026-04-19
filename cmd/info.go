/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/LeeroyDing/seerr-cli/pkg/seerr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Show detailed information about a movie or TV show",
	Long: `Display summary, ratings, and other metadata for a specific media item.
	
Example:
  seerr info 12345 --type movie`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mediaID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: invalid media ID '%s'\n", args[0])
			return
		}

		mediaType, _ := cmd.Flags().GetString("type")
		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Printf("Fetching details for %s ID %d...\n\n", mediaType, mediaID)

		var details *seerr.MediaDetails
		if mediaType == "movie" {
			details, err = client.GetMovieDetails(mediaID)
		} else {
			details, err = client.GetTVDetails(mediaID)
		}

		if err != nil {
			fmt.Printf("Error fetching details: %v\n", err)
			return
		}

		title := details.Title
		if title == "" {
			title = details.Name
		}
		date := details.ReleaseDate
		if date == "" {
			date = details.FirstAirDate
		}

		fmt.Printf("%s (%s)\n", title, strings.Split(date, "-")[0])
		fmt.Printf("Rating: %.1f/10\n", details.VoteAverage)
		
		var genres []string
		for _, g := range details.Genres {
			genres = append(genres, g.Name)
		}
		fmt.Printf("Genres: %s\n", strings.Join(genres, ", "))
		
		fmt.Printf("\nOverview:\n%s\n", details.Overview)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringP("type", "t", "movie", "Media type (movie or tv)")
}
