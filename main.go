// Copyright (c) 2024 Youssouf Drif
// Licensed under the MIT License: https://opensource.org/licenses/MIT

package main

import (
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// init initializes the command-line flags for the application.
// It sets up persistent flags for LTE interface, WiFi interface, WiFi SSID,
// WiFi password, and failure detection delay. It also marks the LTE interface,
// WiFi interface, WiFi SSID, and WiFi password flags as required.
func init() {
	rootCmd.PersistentFlags().StringP("lte-if", "l", "", "LTE interface (required)")
	rootCmd.PersistentFlags().StringP("wifi-if", "w", "", "WiFi interface (required)")
	rootCmd.PersistentFlags().StringP("wifi-ssid", "s", "", "WiFi SSID (required)")
	rootCmd.PersistentFlags().StringP("wifi-password", "p", "", "WiFi password (required)")
	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "Probe server endpoint (required)")
	rootCmd.PersistentFlags().IntP("retry", "r", 5, "Retry count before switching to WiFi (default: 5)")
	rootCmd.MarkPersistentFlagRequired("lte-if")
	rootCmd.MarkPersistentFlagRequired("wifi-if")
	rootCmd.MarkPersistentFlagRequired("wifi-ssid")
	rootCmd.MarkPersistentFlagRequired("wifi-password")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// pingIP uses ICMP to ping an IP address and returns the response time in milliseconds.
// Returns -1 if there is an error or if the ping fails.
func pingIP(ip string) int {
	cmd := exec.Command("ping", "-c", "1", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return -1
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, "1 packets received") {
		return -1
	}

	// Extract response time
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "time=") {
			parts := strings.Split(line, " ")
			for _, part := range parts {
				if strings.HasPrefix(part, "time=") {
					timeStr := strings.TrimPrefix(part, "time=")
					timeStr = strings.TrimSuffix(timeStr, " ms")
					responseTime, err := strconv.ParseFloat(timeStr, 32)
					if err != nil {
						return -1
					}
					return int(responseTime)
				}
			}
		}
	}

	return -1
}

// pingInterface pings an interface and when the retry-count is met with consecutive failures, it returns -1.
func pingInterface(endpoint string, retry int) int {
	log.Info().Msgf("Pinging endpoint %s", endpoint)
	failures := 0
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		log.Warn().Msgf("Stopping ping due to user interrupt...")
		log.Info().Msg("Exiting the program...")
		os.Exit(0)
	}()
	for {
		time.Sleep(time.Second)
		responseTime := pingIP(endpoint)
		if responseTime != -1 {
			failures = 0
		} else {
			failures++
			log.Warn().Msgf("Failed to ping %s. Attempt %d out of %d. Retrying...", endpoint, failures, retry)
			if failures >= retry {
				return -1
			}
		}
	}
}

// connectToWiFi connects to the given wifi bssid with the given password.
func connectToWiFi(ifwifi string, bssid string, password string) error {
	cmd := exec.Command("iwconfig", ifwifi, "essid", bssid)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Info().Msg(string(output))
	cmd = exec.Command("iwconfig", ifwifi, "key", password)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Info().Msg(string(output))
	cmd = exec.Command("dhclient", ifwifi)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Info().Msg(string(output))
	// ping the default router to check if the connection is successful
	for {
		time.Sleep(time.Second)
		output, err := exec.Command("ip", "route", "show", "default", "dev", ifwifi).CombinedOutput()
		if err != nil {
			log.Error().Msgf("Error getting default route after connecting to WiFi: %s", err)
			return nil
		}
		route := strings.Split(string(output), " ")[2]
		responseTime := pingIP(route)
		if responseTime != -1 {
			break
		}
	}
	return nil
}

// changeDefaultRoute sets the default route to the interface set at wifi-if.
func changeDefaultRoute(wifiIF string) {
	cmd := exec.Command("ip", "route", "change", "default", "dev", wifiIF)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Msg(err.Error())
	}
	log.Info().Msg(string(output))
}

var rootCmd = &cobra.Command{
	Use:   "if-reliability",
	Short: "Interface Reliability tool",
	Long:  "Interface Reliability tool is a tool to check the reliability of an interface.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting Interface Reliability tool...")
		lteIF, _ := cmd.Flags().GetString("lte-if")
		wifiIF, _ := cmd.Flags().GetString("wifi-if")
		wifiSSID, _ := cmd.Flags().GetString("wifi-ssid")
		wifiPassword, _ := cmd.Flags().GetString("wifi-password")
		endPoint, _ := cmd.Flags().GetString("endpoint")
		failureDelay, _ := cmd.Flags().GetInt("failure-delay")

		log.Info().Msgf("Starting Interface Reliability tool with:")
		log.Info().Msgf("- LTE interface: %s", lteIF)
		log.Info().Msgf("- WiFi interface: %s", wifiIF)
		log.Info().Msgf("- WiFi SSID: %s", wifiSSID)
		log.Info().Msgf("- WiFi password: %s", wifiPassword)
		log.Info().Msgf("- Failure detection delay: %dms", failureDelay)

		time.Sleep(time.Duration(failureDelay) * time.Millisecond)

		pingInterface(endPoint, 5)
		log.Error().Msgf("Ping toward %s endpoint failed", endPoint)
		connectToWiFi(wifiIF, wifiSSID, wifiPassword)
		log.Info().Msgf("Successfully connected to WiFi with SSID %s", wifiSSID)
		changeDefaultRoute(wifiIF)
		log.Info().Msgf("Successfully changed default route to %s", wifiIF)
		pingInterface(endPoint, 5)
	},
}

// Main starts the command-line interface.
func main() {
	rootCmd.Execute()
}
