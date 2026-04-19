/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/LeeroyDing/seerr-cli/pkg/seerr"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// interactiveCmd represents the interactive command
var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Launch interactive mode",
	Long: `Start a menu-driven session to browse and manage your Seerr instance.
	
Example:
  seerr interactive`,
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("url")
		apiKey := viper.GetString("api-key")

		if url == "" || apiKey == "" {
			fmt.Println("Error: Seerr URL and API Key must be configured first.")
			return
		}

		client := seerr.NewClient(url, apiKey)
		
		for {
			prompt := promptui.Select{
				Label: "Seerr CLI - Main Menu",
				Items: []string{"Search Media", "Browse Trending", "My Requests", "My Profile", "Exit"},
			}

			_, result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			switch result {
			case "Search Media":
				runSearchFlow(client)
			case "Browse Trending":
				runBrowseFlow(client, "trending")
			case "My Requests":
				runListFlow(client)
			case "My Profile":
				runUserFlow(client)
			case "Exit":
				fmt.Println("Goodbye!")
				os.Exit(0)
			}
		}
	},
}

func runSearchFlow(c *seerr.Client) {
	prompt := promptui.Prompt{
		Label: "Enter search query",
	}

	query, err := prompt.Run()
	if err != nil {
		return
	}

	resp, err := c.Search(query)
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
		return
	}

	if len(resp.Results) == 0 {
		fmt.Println("No results found.")
		return
	}

	var options []string
	for _, res := range resp.Results {
		title := res.Title
		if title == "" {
			title = res.Name
		}
		options = append(options, fmt.Sprintf("[%s] %s (%s)", res.MediaType, title, res.ReleaseDate))
	}
	options = append(options, "Back")

	sel := promptui.Select{
		Label: "Select item to request",
		Items: options,
	}

	i, result, err := sel.Run()
	if err != nil || result == "Back" {
		return
	}

	item := resp.Results[i]
	fmt.Printf("Requesting %s...\n", item.Title)
	err = c.Request(item.ID, item.MediaType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Successfully requested!")
	}
}

func runBrowseFlow(c *seerr.Client, cat string) {
	resp, err := c.GetTrending()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Similar to search flow but with results
	var options []string
	for _, res := range resp.Results {
		title := res.Title
		if title == "" {
			title = res.Name
		}
		options = append(options, fmt.Sprintf("[%s] %s", res.MediaType, title))
	}
	options = append(options, "Back")

	sel := promptui.Select{
		Label: "Select item to request",
		Items: options,
	}

	i, result, err := sel.Run()
	if err != nil || result == "Back" {
		return
	}

	item := resp.Results[i]
	err = c.Request(item.ID, item.MediaType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Successfully requested!")
	}
}

func runListFlow(c *seerr.Client) {
	resp, err := c.ListRequests(20, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	for _, req := range resp.Results {
		title := req.Media.Title
		if title == "" {
			title = req.Media.Name
		}
		fmt.Printf("- %s (ID: %d)\n", title, req.ID)
	}
	fmt.Println("\nPress Enter to continue")
	fmt.Scanln()
}

func runUserFlow(c *seerr.Client) {
	user, err := c.GetMe()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("\nUser: %s (Requests: %d)\n", user.DisplayName, user.RequestCount)
	fmt.Println("Press Enter to continue")
	fmt.Scanln()
}

func init() {
	rootCmd.AddCommand(interactiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// interactiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// interactiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
