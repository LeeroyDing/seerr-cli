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

// denyCmd represents the deny command
var denyCmd = &cobra.Command{
	Use:   "deny [requestID]",
	Short: "Deny a media request",
	Long: `Deny a specific pending request by its ID.
The request will be rejected and removed from the pending list.

Example:
  seerr admin deny 123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requestID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: invalid request ID '%s'\n", args[0])
			return
		}

		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Printf("Denying request ID %d...\n", requestID)

		err = client.DeclineRequest(requestID)
		if err != nil {
			fmt.Printf("Error denying request: %v\n", err)
			return
		}

		fmt.Println("Successfully denied!")
	},
}

func init() {
	adminCmd.AddCommand(denyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// denyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// denyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
