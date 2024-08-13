/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/raspiantoro/goms/internal/app"
	"github.com/spf13/cobra"
)

// migrationCmd represents the migration command
var migrationCmd = &cobra.Command{
	Use:   "migration [migration name]",
	Short: "Add new migration file",
	Long: `Create and add a new migration file to your Golang project. 
The migration name will be prefixed with the current datetime in the format 'YYYYMMDDhhmm' as a migration version/identifier in the generated file name.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print("Error: add migration requires a name for the migration.\n\n")
			cmd.Help()
			return
		} else if len(args) > 1 {
			fmt.Print("Error: add migration only accept one arguments [migration name]\n\n")
			cmd.Help()
			return
		}

		cfg := app.Config{Name: args[0]}

		if flagDir != "" {
			cfg.Path = flagDir
		}

		app.CreateMigration(cfg)
	},
}

func init() {
	addCmd.AddCommand(migrationCmd)

	// Here you will define your flags and configuration settings.
	migrationCmd.Flags().StringVarP(&flagDir, "dir", "d", "", "goms files directory")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
