/*
Copyright © 2023 Rhys Campbell
Copyrights appy to this source code.
Check LICENSE for details.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/rhysmeister/pomo/app"
	"github.com/rhysmeister/pomo/pomodoro"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Interactive Pomodoro Timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := getRepo()
		if err != nil {
			return err
		}

		config := pomodoro.NewConfig(
			repo,
			viper.GetDuration("pomo"),
			viper.GetDuration("short"),
			viper.GetDuration("long"),
		)
		return rootAction(os.Stdout, config)
	},
}

var cfgFile string

func rootAction(out io.Writer, config *pomodoro.IntervalConfig) error {
	a, err := app.New(config)
	if err != nil {
		return err
	}
	return a.Run()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.pomo.yaml)")

	rootCmd.Flags().DurationP("pomo", "p", 25*time.Minute,
		"Pomodoro duration")
	rootCmd.Flags().DurationP("short", "s", 5*time.Minute,
		"Short break duration")
	rootCmd.Flags().DurationP("long", "l", 15*time.Minute,
		"Long break duration")

	viper.BindPFlag("pomo", rootCmd.Flags().Lookup("pomo"))
	viper.BindPFlag("short", rootCmd.Flags().Lookup("short"))
	viper.BindPFlag("long", rootCmd.Flags().Lookup("long"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pomo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pomo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
