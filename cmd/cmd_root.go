package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = `athena`
)

var (
	cfgFile string
	debug   bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "athena",
	Short: "A service to backup Cassandra keyspaces to MS Azure",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig, enableDebug)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.athena.yaml)")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "if set, this enables debugging")
}

func initConfig() {
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("athena")
		viper.AddConfigPath("/etc/athena")
		viper.AddConfigPath("$HOME/.athena")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix(appName)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Loaded config file:", viper.ConfigFileUsed())
	}
}
