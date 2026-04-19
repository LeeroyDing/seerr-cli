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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your media requests",
	Long: `Display a list of your pending and approved requests.
	
Example:
  seerr list`,
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Println("Fetching requests...")

		resp, err := client.ListRequests(20, 0)
		if err != nil {
			fmt.Printf("Error fetching requests: %v\n", err)
			return
		}

		if len(resp.Results) == 0 {
			fmt.Println("No requests found.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tMEDIA\tREQUESTED AT")
		fmt.Fprintln(w, "--\t------\t-----\t------------")

		for _, req := range resp.Results {
			status := formatStatus(req.Status)
			title := req.Media.Title
			if title == "" {
				title = req.Media.Name
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", req.ID, status, title, req.CreatedAt)
		}
		w.Flush()
	},
}

func formatStatus(s int) string {
	switch s {
	case 1:
		return "PENDING"
	case 2:
		return "APPROVED"
	case 3:
		return "DECLINED"
	default:
		return "UNKNOWN"
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
