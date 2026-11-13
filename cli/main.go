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

const requiredAcceptMediaType = "application/vnd.ledger-grievance.v1+json; charset=utf-8"

var (
	apiURL          string
	nanosecondEpoch string
	slaMode         bool
)

func buildRequest(method, path string, payload any) (*http.Request, error) {
	var bodyReader io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, apiURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if payload != nil || method == http.MethodPost || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}
	if method == http.MethodPost || method == http.MethodPatch || method == http.MethodDelete {
		req.Header.Set("Accept", requiredAcceptMediaType)
	}
	if nanosecondEpoch != "" {
		req.Header.Set("X-Nanosecond-Epoch", nanosecondEpoch)
	}
	if slaMode {
		req.Header.Set("X-SLA-Mode", "100-percent")
	}
	return req, nil
}

func request(method, path string, payload any) {
	req, err := buildRequest(method, path, payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(string(body))
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "ledger",
		Short: "Grievance Ledger CLI",
	}

	rootCmd.PersistentFlags().StringVar(&apiURL, "url", "http://localhost:8000", "API Base URL")
	rootCmd.PersistentFlags().StringVar(&nanosecondEpoch, "nanosecond-epoch", "", "X-Nanosecond-Epoch header for PATCH/DELETE")
	rootCmd.PersistentFlags().BoolVar(&slaMode, "sla-mode", false, "Send X-SLA-Mode: 100-percent")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List incidents",
		Run: func(cmd *cobra.Command, args []string) {
			request(http.MethodGet, "/incidents", nil)
		},
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new incident",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				fmt.Println("Usage: create <reporter_id> <subject> <description>")
				return
			}
			payload := map[string]any{
				"reporter_id":                     args[0],
				"subject":                         args[1],
				"description":                     args[2],
				"assumed_good_intentions":         true,
				"promised_to_be_kind_to_yourself": true,
				"requires_accommodation":          false,
			}
			request(http.MethodPost, "/incidents", payload)
		},
	}

	complimentCmd := &cobra.Command{
		Use:   "compliment",
		Short: "Get a wholesome compliment",
		Run: func(cmd *cobra.Command, args []string) {
			request(http.MethodGet, "/compliments", nil)
		},
	}

	wisdomCmd := &cobra.Command{
		Use:   "wisdom",
		Short: "Get Gopher Wisdom",
		Run: func(cmd *cobra.Command, args []string) {
			request(http.MethodGet, "/wisdom", nil)
		},
	}

	bouquetCmd := &cobra.Command{
		Use:   "bouquet",
		Short: "Get a wholesome package bouquet",
		Run: func(cmd *cobra.Command, args []string) {
			request(http.MethodGet, "/bouquet", nil)
		},
	}

	vouchCmd := &cobra.Command{
		Use:   "vouch",
		Short: "Vouch for an incident",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Usage: vouch <incident_id>")
				return
			}
			request(http.MethodPost, "/incidents/"+args[0]+"/vouch", nil)
		},
	}

	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Check system deep health",
		Run: func(cmd *cobra.Command, args []string) {
			request(http.MethodGet, "/health/deep", nil)
		},
	}

	rootCmd.AddCommand(listCmd, createCmd, complimentCmd, wisdomCmd, bouquetCmd, vouchCmd, healthCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
