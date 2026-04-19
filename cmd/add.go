/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"regexp"

	"github.com/LeeroyDing/seerr-cli/pkg/seerr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [imdbLinkOrID]",
	Short: "Add a movie or show using an IMDb link or ID",
	Long: `Quickly add media to Seerr by providing its IMDb URL or ID.
	
Example:
  seerr add https://www.imdb.com/title/tt0111161/
  seerr add tt0111161`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		
		// Extract IMDb ID (tt followed by digits)
		re := regexp.MustCompile(`tt\d+`)
		imdbID := re.FindString(input)
		
		if imdbID == "" {
			fmt.Println("Error: could not find a valid IMDb ID in the input.")
			return
		}

		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Printf("Looking up IMDb ID %s...\n", imdbID)

		resp, err := client.Search(imdbID)
		if err != nil {
			fmt.Printf("Error searching: %v\n", err)
			return
		}

		if len(resp.Results) == 0 {
			fmt.Printf("No items found in Seerr for IMDb ID %s.\n", imdbID)
			return
		}

		// Usually the first result for an IMDb ID is the exact match
		item := resp.Results[0]
		title := item.Title
		if title == "" {
			title = item.Name
		}

		fmt.Printf("Found: %s (%s) [%s]\n", title, item.ReleaseDate, item.MediaType)
		fmt.Printf("Requesting %s...\n", title)

		err = client.Request(item.ID, item.MediaType)
		if err != nil {
			fmt.Printf("Error requesting: %v\n", err)
			return
		}

		fmt.Println("Successfully added!")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
