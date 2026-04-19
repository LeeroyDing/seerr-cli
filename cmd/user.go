/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/LeeroyDing/seerr-cli/pkg/seerr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Show your Seerr user profile",
	Long: `Display information about your Seerr user account, including your
username, email, and permissions.

Example:
  seerr user`,
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		fmt.Println("Fetching user info...")

		user, err := client.GetMe()
		if err != nil {
			fmt.Printf("Error fetching user info: %v\n", err)
			return
		}

		fmt.Printf("User Profile:\n")
		fmt.Printf("  Display Name:  %s\n", user.DisplayName)
		fmt.Printf("  User ID:       %d\n", user.ID)
		fmt.Printf("  Total Requests: %d\n", user.RequestCount)
		fmt.Printf("  Permissions:    0x%x\n", user.Permissions)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
