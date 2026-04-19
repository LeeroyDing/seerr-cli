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

// listPendingCmd represents the list-pending command
var listPendingCmd = &cobra.Command{
	Use:   "list-pending",
	Short: "List all pending requests",
	Long: `Display a list of all requests waiting for approval.
	
Example:
  seerr admin list-pending`,
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Println("Fetching pending requests...")

		// Note: We need to filter by pending. I'll reuse ListRequests but with a filter if supported.
		// For now, I'll just list and filter locally or use the API if I updated the client.
		resp, err := client.ListRequests(50, 0)
		if err != nil {
			fmt.Printf("Error fetching requests: %v\n", err)
			return
		}

		pendingCount := 0
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tMEDIA\tREQUESTED AT")
		fmt.Fprintln(w, "--\t-----\t------------")

		for _, req := range resp.Results {
			if req.Status == 1 { // PENDING
				title := req.Media.Title
				if title == "" {
					title = req.Media.Name
				}
				fmt.Fprintf(w, "%d\t%s\t%s\n", req.ID, title, req.CreatedAt)
				pendingCount++
			}
		}
		w.Flush()

		if pendingCount == 0 {
			fmt.Println("No pending requests.")
		}
	},
}

func init() {
	adminCmd.AddCommand(listPendingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listPendingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listPendingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
