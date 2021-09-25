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

	"github.com/Luet-lab/luet-portage-converter/pkg/converter"

	. "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"

	"github.com/spf13/cobra"
)

func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate luet specs.",
		PreRun: func(cmd *cobra.Command, args []string) {
			to, _ := cmd.Flags().GetString("to")
			if to == "" {
				fmt.Println("Missing --to argument")
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
			stage2, _ := cmd.Flags().GetBool("disable-stage2")
			stage3, _ := cmd.Flags().GetBool("disable-stage3")
			stage4, _ := cmd.Flags().GetBool("enable-stage4")
			debug, _ := cmd.Flags().GetBool("debug")
			backend, _ := cmd.Flags().GetString("backend")
			ignoreMissingDeps, _ := cmd.Flags().GetBool("ignore-missing-deps")
			pkgs, _ := cmd.Flags().GetStringArray("pkg")
			withPortagePkgs, _ := cmd.Flags().GetBool("with-portage-pkg")
			disableConflicts, _ := cmd.Flags().GetBool("disable-conflicts")
			layer4Rdepends, _ := cmd.Flags().GetBool("layer4rdepends")

			converter := converter.NewPortageConverter(to, backend)
			converter.Override = override
			converter.IgnoreMissingDeps = ignoreMissingDeps
			converter.TreePaths = treePath
			converter.WithPortagePkgs = withPortagePkgs
			converter.DisableStage2 = stage2
			converter.DisableStage3 = stage3
			converter.DisableStage4 = !stage4
			converter.DisableConflicts = disableConflicts
			converter.UsingLayerForRuntime = layer4Rdepends

			if debug {
				LuetCfg.GetGeneral().Debug = debug
			}

			if len(pkgs) > 0 {
				converter.SetFilteredPackages(pkgs)
			}

			if len(treePath) == 0 {
				DebugC(GetAurora().Bold("ATTENTION! No trees defined."))
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

	cmd.Flags().String("to", "", "Targer tree where bump new specs.")
	cmd.Flags().Bool("override", false, "Override existing specs if already present.")
	cmd.Flags().StringArrayP("pkg", "p", []string{},
		"Define the list of the packages to generate instead of the full list defined in rules file.")
	cmd.Flags().Bool("with-portage-pkg", false, "Generate portage packages for every required package.")
	cmd.Flags().Bool("disable-conflicts", false, "Disable elaboration of runtime and buildtime conflicts.")
	cmd.Flags().Bool("layer4rdepends", false, "Check layer for runtime deps and skip generation.")

	return cmd
}
