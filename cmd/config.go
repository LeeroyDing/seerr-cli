/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Seerr CLI",
	Long: `Set the Seerr instance URL and API Key.
	
Example:
  seerr config --url https://seerr.example.com --api-key your-api-key`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if url == "" && apiKey == "" {
			fmt.Printf("Current Configuration:\n")
			fmt.Printf("  URL: %s\n", viper.GetString("url"))
			fmt.Printf("  API Key: %s\n", maskApiKey(viper.GetString("api-key")))
			return
		}

		if url != "" {
			viper.Set("url", url)
		}
		if apiKey != "" {
			viper.Set("api-key", apiKey)
		}

		err := viper.WriteConfig()
		if err != nil {
			// If config file doesn't exist, try creating it
			err = viper.SafeWriteConfig()
			if err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}
		}
		fmt.Println("Configuration updated successfully.")
	},
}

func maskApiKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "...." + key[len(key)-4:]
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringP("url", "u", "", "Seerr instance URL")
	configCmd.Flags().StringP("api-key", "k", "", "Seerr API Key")
}
