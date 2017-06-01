package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// appName and version are fed through the go linker
	appName = `athena`
	version = ``

	cfgFile string

	baseCmd = &cobra.Command{
		Use:   appName,
		Short: "A tool to backup Cassandra keyspaces",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	baseCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.athena.yaml)")
}

func main() {
	if err := baseCmd.Execute(); err != nil {
		log.Fatal(err)
	}
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
		log.Print("Loaded config file:", viper.ConfigFileUsed())
	} else {
		log.Print("Could not load configuration:", err)
	}
}
