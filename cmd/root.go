package cmd

import (
	"github.com/kunicmarko20/deko-cli/util"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use: "deko-cli",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(releaseCmd)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		util.Exit(err)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".deko-cli")
	viper.SetConfigType("yaml")

	if _, err := os.Stat(viper.ConfigFileUsed()); err != nil {
		viper.SafeWriteConfig()
	}

	if err := viper.ReadInConfig(); err != nil {
		util.Exit(err)
	}
}
