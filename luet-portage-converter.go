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

package main

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

	version = "0.1.0"
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
		PreRun: func(cmd *cobra.Command, args []string) {
			to, _ := cmd.Flags().GetString("to")
			if to == "" {
				fmt.Println("Missing --to argument")
				os.Exit(1)
			}

			err := initConfig()
			if err != nil {
				fmt.Println("Error on setup config/logger: " + err.Error())
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			treePath, _ := cmd.Flags().GetStringArray("tree")
			reposcanSources, _ := cmd.Flags().GetStringArray("reposcan-files")
			disableUseFlags, _ := cmd.Flags().GetStringArray("disable-use-flag")
			to, _ := cmd.Flags().GetString("to")
			rulesFile, _ := cmd.Flags().GetString("rules")
			override, _ := cmd.Flags().GetBool("override")
			backend, _ := cmd.Flags().GetString("backend")
			ignoreMissingDeps, _ := cmd.Flags().GetBool("ignore-missing-deps")
			pkgs, _ := cmd.Flags().GetStringArray("pkg")

			converter := converter.NewPortageConverter(to, backend)
			converter.Override = override
			converter.IgnoreMissingDeps = ignoreMissingDeps
			converter.TreePaths = treePath

			if len(pkgs) > 0 {
				converter.SetFilteredPackages(pkgs)
			}

			err := converter.LoadRules(rulesFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = converter.LoadTrees(treePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if len(reposcanSources) > 0 {
				for _, source := range reposcanSources {
					converter.Specs.AddReposcanSource(source)
				}
			}

			if len(disableUseFlags) > 0 {
				converter.Specs.ReposcanDisabledUseFlags = append(converter.Specs.ReposcanDisabledUseFlags, disableUseFlags...)
			}

			err = converter.Generate()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		},
	}

	rootCmd.Flags().StringArrayP("tree", "t", []string{}, "Path of the tree to use.")
	rootCmd.Flags().String("to", "", "Targer tree where bump new specs.")
	rootCmd.Flags().String("rules", "", "Rules file.")
	rootCmd.Flags().Bool("override", false, "Override existing specs if already present.")
	rootCmd.Flags().String("backend", "reposcan", "Select backend resolver: qdepends|reposcan.")
	rootCmd.Flags().StringArray("reposcan-files", []string{}, "Append additional reposcan files. Only for reposcan backend.")
	rootCmd.Flags().StringArray("disable-use-flag", []string{}, "Append additional use flags to disable.")
	rootCmd.Flags().Bool("ignore-missing-deps", false, "Ignore missing deps on resolver.")
	rootCmd.Flags().StringArrayP("pkg", "p", []string{},
		"Define the list of the packages to generate instead of the full list defined in rules file.")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
