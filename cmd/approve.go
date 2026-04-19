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

// approveCmd represents the approve command
var approveCmd = &cobra.Command{
	Use:   "approve [requestID]",
	Short: "Approve a media request",
	Long: `Approve a pending request so it can be processed.
	
Example:
  seerr admin approve 1`,
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
		fmt.Printf("Approving request ID %d...\n", requestID)

		err = client.ApproveRequest(requestID)
		if err != nil {
			fmt.Printf("Error approving request: %v\n", err)
			return
		}

		fmt.Println("Successfully approved!")
	},
}

func init() {
	adminCmd.AddCommand(approveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// approveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// approveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
