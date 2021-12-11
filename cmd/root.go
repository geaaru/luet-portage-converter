/*
Copyright (C) 2020-2021  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Luet-lab/luet-portage-converter/pkg/converter"

	. "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	cliName = `Copyright (c) 2020-2021 - Daniele Rondina

Portage/Overlay converter for Luet specs.`

	version = "0.5.1"
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
		Use:     "luet-portage-converter --",
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
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
