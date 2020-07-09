package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configListCmd = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			for key, value := range viper.AllSettings() {
				fmt.Println(key + ": " + value.(string))
			}
		},
	}
)
