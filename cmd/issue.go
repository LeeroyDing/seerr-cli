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

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue [id]",
	Short: "Report an issue with a media item",
	Long: `Report a problem (e.g., missing episodes, bad quality) for a specific 
media item by its ID. Requires the item ID and type flag.

Example:
  seerr issue 27205 --type movie --message "The video is buffering too much"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mediaID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: invalid media ID '%s'\n", args[0])
			return
		}

		issueTypeStr, _ := cmd.Flags().GetString("type")
		message, _ := cmd.Flags().GetString("message")
		
		issueType := parseIssueType(issueTypeStr)
		
		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Printf("Reporting %s issue for media ID %d...\n", issueTypeStr, mediaID)

		err = client.CreateIssue(mediaID, issueType, message)
		if err != nil {
			fmt.Printf("Error reporting issue: %v\n", err)
			return
		}

		fmt.Println("Successfully reported!")
	},
}

func parseIssueType(s string) int {
	switch strings.ToLower(s) {
	case "video":
		return 1
	case "audio":
		return 2
	case "subs", "subtitles":
		return 3
	default:
		return 4 // Other
	}
}

func init() {
	rootCmd.AddCommand(issueCmd)

	issueCmd.Flags().StringP("type", "t", "video", "Issue type (video, audio, subs, other)")
	issueCmd.Flags().StringP("message", "m", "", "Detailed description of the issue")
	issueCmd.MarkFlagRequired("message")
}
