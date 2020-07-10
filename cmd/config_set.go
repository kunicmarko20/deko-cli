package cmd

import (
	"fmt"
	"github.com/kunicmarko20/deko-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configSetCmd = &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Creates or Updates a config.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set(args[0], args[1])
			err := viper.WriteConfig()

			if err != nil {
				util.Exit(err)
			}

			fmt.Println("successfully set config with key: " + args[0])
		},
	}
)
