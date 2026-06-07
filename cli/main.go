package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var apiURL string

func main() {
	rootCmd := &cobra.Command{
		Use:   "ledger",
		Short: "Grievance Ledger CLI",
	}

	rootCmd.PersistentFlags().StringVar(&apiURL, "url", "http://localhost:8000", "API Base URL")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List incidents",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := http.Get(apiURL + "/incidents")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			fmt.Println(string(body))
		},
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new incident",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Println("Usage: create <reporter_id> <subject> <description>")
				return
			}
			payload := map[string]interface{}{
				"reporter_id":               args[0],
				"subject":                   args[1],
				"description":               args[2],
				"assumed_good_intentions":   true,
				"promised_to_be_kind_to_yourself": true,
				"requires_accommodation":    false,
			}
			jsonData, _ := json.Marshal(payload)
			resp, err := http.Post(apiURL+"/incidents", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			fmt.Println(string(body))
		},
	}

	complimentCmd := &cobra.Command{
		Use:   "compliment",
		Short: "Get a wholesome compliment",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := http.Get(apiURL + "/compliments")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			fmt.Println(string(body))
		},
	}

	rootCmd.AddCommand(listCmd, createCmd, complimentCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
