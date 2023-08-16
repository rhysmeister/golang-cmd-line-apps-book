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
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List todo items",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		active := viper.GetBool("active")
		return listAction(os.Stdout, apiRoot, active)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().Bool("active", false, "List only active tasks")
	viper.BindPFlag("active", listCmd.PersistentFlags().Lookup("active"))
}

func listAction(out io.Writer, apiRoot string, active bool) error {
	items, err := getAll(apiRoot)
	if err != nil {
		return err
	}
	return printAll(out, items, active)
}

func printAll(out io.Writer, items []item, active bool) error {
	w := tabwriter.NewWriter(out, 3, 2, 0, ' ', 0)

	for k, v := range items {
		done := "-"
		if v.Done {
			done = "X"
		}

		if (done == "X" && !active) || (done == "-") {
			fmt.Fprintf(w, "%s\t%d\t%s\t\n", done, k+1, v.Task)
		}
	}
	return w.Flush()
}
