/*
	Copyright © 2021-2023 Macaroni OS Linux
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
	cliName = `Copyright (c) 2020-2023 - Daniele Rondina

Portage/Overlay converter for Anise specs.`

	version = "0.14.0"
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

	rootCmd.PersistentFlags().StringArrayP("tree", "t", []string{}, "Path of the tree to use.")
	rootCmd.PersistentFlags().String("rules", "", "Rules file.")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug verbosity.")
	rootCmd.PersistentFlags().String("backend", "reposcan", "Select backend resolver: qdepends|reposcan.")
	rootCmd.PersistentFlags().Bool("disable-stage2", false, "Disable stage2 phase.")
	rootCmd.PersistentFlags().Bool("disable-stage3", false, "Disable stage3 phase.")
	rootCmd.PersistentFlags().Bool("enable-stage4", false, "Enable experimental stage4 phase.")
	rootCmd.PersistentFlags().StringArray("reposcan-files", []string{}, "Append additional reposcan files. Only for reposcan backend.")
	rootCmd.PersistentFlags().StringArray("disable-use-flag", []string{}, "Append additional use flags to disable.")
	rootCmd.PersistentFlags().Bool("ignore-missing-deps", false, "Ignore missing deps on resolver.")

	rootCmd.AddCommand(
		newGenerateCommand(),
		newReposcanResolveCommand(),
		newSyncCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
