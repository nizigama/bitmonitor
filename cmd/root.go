package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/martinlindhe/notify"
	"github.com/spf13/cobra"
)

var numericInput int

var rootCmd = &cobra.Command{
	Use:   "BitMonitor",
	Short: "BitMonitor is a bitcoin price monitor app",
	Long: `BitMonitor helps you monitor the price of the bitcoin
by getting a notification as often as per minute, you can configure this on the launch 
or leave it by default to the 5 minutes notification`,
	Run: func(cmd *cobra.Command, args []string) {
		interval, err := getNotificationsInterval()
		if err != nil {
			return
		}
		intervalListener(interval)
	},
}

// Execute executes the root command
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getNotificationsInterval() (int, error) {
	fmt.Print("How often in minutes: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	input = strings.Split(input, " ")[0]

	numericInput, err := strconv.Atoi(input)

	if err != nil {
		fmt.Println("You must enter a valid number representing minutes!")
		return 0, err
	}
	return numericInput, nil
}

func intervalListener(minutes int) {
	ticker := time.NewTicker(time.Minute * time.Duration(minutes))

	defer ticker.Stop()

	fetchPriceUpdates()

	for {
		select {
		case _ = <-ticker.C:
			fetchPriceUpdates()
		}
	}
}

func fetchPriceUpdates() {
	fmt.Println("Fetching updates ...")

	resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")

	if err != nil {
		fmt.Println("Error getting updates!\nCheck your internet connection")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("BitMonitor error\nContact Support!!!")
		os.Exit(1)
	}
	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("BitMonitor error\nContact Support!!!")
		os.Exit(1)
	}

	var responseOBJ map[string]interface{}

	err = json.Unmarshal(response, &responseOBJ)

	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		os.Exit(1)
	}

	for _, v := range responseOBJ {
		switch v.(type) {
		case map[string]interface{}:
			for k, v2 := range v.(map[string]interface{}) {
				if k == "USD" {
					for k2, v3 := range v2.(map[string]interface{}) {
						if k2 == "rate" {
							message := fmt.Sprintf("Current bitcoin value is: %s USD", v3)
							log.Println(message)
							notify.Alert("BitMonitor", "Bitcoin price update", message, "path/to/icon.png")
						}
					}
				}
			}
		default:

		}
	}

}
