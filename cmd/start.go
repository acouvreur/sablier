package cmd

import (
	"github.com/acouvreur/sablier/v2/app"
	"github.com/acouvreur/sablier/v2/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Sablier server",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.NewConfig()
		viper.Unmarshal(&conf)

		err := app.Start(conf)
		if err != nil {
			panic(err)
		}
	},
}
