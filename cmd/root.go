/*
	Copyright © 2021-2024 Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/macaroni-os/anise-portage-converter/pkg/converter"

	. "github.com/geaaru/luet/pkg/config"
	. "github.com/geaaru/luet/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	cliName = `Copyright (c) 2020-2024 - Daniele Rondina

Portage/Overlay bridge for Macaroni OS Sambuca stack.`

	version = "0.16.0"
)

func initConfig() error {
	LuetCfg.Viper.SetEnvPrefix("LUET")
	LuetCfg.Viper.AutomaticEnv() // read in environment variables that match

	// Create EnvKey Replacer for handle complex structure
	replacer := strings.NewReplacer(".", "__")
	LuetCfg.Viper.SetEnvKeyReplacer(replacer)
	LuetCfg.Viper.SetTypeByDefaultValue(true)

	err := LuetCfg.Viper.Unmarshal(&LuetCfg)
	if err != nil {
		return err
	}

	InitAurora()
	NewSpinner()

	return nil
}

func Execute() {

	var rootCmd = &cobra.Command{
		Use:     "anise-portage-converter --",
		Short:   cliName,
		Version: fmt.Sprintf("%s-g%s %s", version, converter.BuildCommit, converter.BuildTime),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := initConfig()
			if err != nil {
				fmt.Println("Error on setup config/logger: " + err.Error())
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug verbosity.")

	rootCmd.AddCommand(
		newGenerateCommand(),
		newReposcanResolveCommand(),
		newReposcanMetadataCommand(),
		newReposcanGenerateCommand(),
		newSyncCommand(),
		newPortageCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
