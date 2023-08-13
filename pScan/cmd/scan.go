/*
Copyright © 2023 Rhys Campbell
Copyrights appy to this source code.
Check LICENSE for details.
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rhysmeister/pScan/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		ports := viper.GetString("ports")
		var port_list []int
		timeout := viper.GetInt("timeout")
		protocol := viper.GetString("protocol")
		if protocol != "tcp" && protocol != "udp" {
			return fmt.Errorf("invalid protocol value supplied: %s", protocol)
		}
		filter := viper.GetString("filter")
		if filter != "open" && filter != "closed" && filter != "both" {
			return fmt.Errorf("invalid filter value supplied: %s", filter)
		}

		if strings.Contains(ports, ",") {
			for _, p := range strings.Split(ports, ",") {
				p, err := strconv.Atoi(p)
				if err != nil {
					return errors.New("invalid value passed for --ports")
				}
				if p < 1 || p > 65535 {
					return fmt.Errorf("invalid port number: %d", p)
				}
				port_list = append(port_list, p)
			}
		} else if strings.Contains(ports, "-") {
			start, err := strconv.Atoi(strings.Split(ports, "-")[0])
			if err != nil {
				return errors.New("invalid value passed for --ports")
			}
			end, err := strconv.Atoi(strings.Split(ports, "-")[1])
			if err != nil {
				return errors.New("invalid value passed for --ports")
			}
			for i := start; i <= end; i++ {
				if i < 1 || i > 65535 {
					return fmt.Errorf("invalid port number: %d", i)
				}
				port_list = append(port_list, i)
			}
		} else if p, err := strconv.Atoi(ports); err == nil {
			if p < 1 || p > 65535 {
				return fmt.Errorf("invalid port number: %d", p)
			}
			port_list = append(port_list, p)
		} else {
			return errors.New("invalid value passed for --ports")
		}
		return scanAction(os.Stdout, hostsFile, port_list, timeout, protocol, filter)
	},
}

func scanAction(out io.Writer, hostsFile string, ports []int, timeout int, protocol string, filter string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	results := scan.Run(hl, ports, timeout, protocol)
	return printResults(out, results, filter)
}

func printResults(out io.Writer, results []scan.Results, filter string) error {
	message := ""

	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)

		if r.NotFound {
			message += fmt.Sprintf(" Host not found\n\n")
			continue
		}
		message += fmt.Sprintln()
		for _, p := range r.PortStates {
			if filter == "both" || (filter == "open" && p.Open) || (filter == "closed" && !p.Open) {
				message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
			}
		}
		message += fmt.Sprintln()
	}
	_, err := fmt.Fprint(out, message)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.PersistentFlags().StringP("ports", "p", "22,80,443",
		"ports to scan")
	viper.BindPFlag("ports", scanCmd.PersistentFlags().Lookup("ports"))

	scanCmd.PersistentFlags().StringP("protocol", "x", "tcp",
		"Protocol to scan")
	viper.BindPFlag("protocol", scanCmd.PersistentFlags().Lookup("protocol"))

	scanCmd.PersistentFlags().IntP("timeout", "t", 1,
		"timeout for scan in seconds")
	viper.BindPFlag("timeout", scanCmd.PersistentFlags().Lookup("timeout"))

	scanCmd.PersistentFlags().StringP("filter", "b", "both",
		"Display open ports, closed ports or both")
	viper.BindPFlag("filter", scanCmd.PersistentFlags().Lookup("filter"))
}
