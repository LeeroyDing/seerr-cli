/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/LeeroyDing/seerr-cli/pkg/seerr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "request [id]",
	Short: "Request a movie or TV show",
	Long: `Submit a request for a specific media item by its ID. 
Requires the item ID and type flag.

Example:
  seerr request 27205 --type movie
  seerr request 1399 --type tv`,
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
		fmt.Printf("Requesting %s ID %d...\n", mediaType, mediaID)

		err = client.Request(mediaID, mediaType)
		if err != nil {
			fmt.Printf("Error requesting media: %v\n", err)
			return
		}

		fmt.Println("Successfully requested!")
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringP("type", "t", "movie", "Media type (movie or tv)")
}
