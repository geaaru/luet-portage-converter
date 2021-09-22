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
	"github.com/Luet-lab/luet-portage-converter/pkg/specs"

	. "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"

	"github.com/spf13/cobra"
)

func newReposcanResolveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reposcan-resolve",
		Aliases: []string{"rr"},
		Short:   "Resolve a package from reposcan tree.",
		PreRun: func(cmd *cobra.Command, args []string) {
			pkg, _ := cmd.Flags().GetString("pkg")
			if pkg == "" {
				fmt.Println("Missing --pkg argument")
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
			jsonOutput, _ := cmd.Flags().GetBool("json")
			backend, _ := cmd.Flags().GetString("backend")
			ignoreMissingDeps, _ := cmd.Flags().GetBool("ignore-missing-deps")
			pkg, _ := cmd.Flags().GetString("pkg")

			converter := converter.NewPortageConverter(to, backend)
			converter.Override = override
			converter.IgnoreMissingDeps = ignoreMissingDeps
			converter.TreePaths = treePath
			converter.DisableStage2 = stage2
			converter.DisableStage3 = stage3
			converter.DisableStage4 = !stage4

			if debug {
				LuetCfg.GetGeneral().Debug = debug
			}

			err := converter.LoadRules(rulesFile)
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

			err = converter.InitConverter(false)
			if err != nil {
				Error(fmt.Sprintf("Error on init converter: %s", err.Error()))
				os.Exit(1)
			}

			artefact := *specs.NewPortageConverterArtefact(pkg)

			// Check if it's present artefact from map
			art, err := converter.Specs.GetArtefactByPackage(pkg)
			if err == nil {
				DebugC(fmt.Sprintf("[%s] Using artefact from map. Uses disabled: %s, enabled: %s",
					pkg, art.Uses.Disabled, art.Uses.Enabled))
				// POST: use artefact from map.
				artefact = *art
			}

			opts := specs.PortageResolverOpts{
				EnableUseFlags:   artefact.Uses.Enabled,
				DisabledUseFlags: artefact.Uses.Disabled,
			}

			solution, err := converter.Resolver.Resolve(pkg, opts)
			if err != nil {
				Error(fmt.Sprintf("Error on resolve %s: %s", pkg, err.Error()))
				os.Exit(1)
			}

			if jsonOutput {
				fmt.Println(solution)
			} else {

				fmt.Println(fmt.Sprintf("Package: %s-%s",
					solution.Package.GetPackageNameWithSlot(),
					solution.Package.GetPVR()))

				fmt.Println(fmt.Sprintf("Description: %s",
					solution.Description))

				fmt.Println("Use flags:")
				for _, use := range solution.Package.UseFlags {
					fmt.Println("- " + use)
				}

				fmt.Println("Build deps:")
				for _, p := range solution.BuildDeps {
					fmt.Println("-" + p.GetPackageNameWithSlot())
				}

				fmt.Println("Runtime deps:")
				for _, p := range solution.RuntimeDeps {
					fmt.Println("-" + p.GetPackageNameWithSlot())
				}

				fmt.Println("Labels:")
				for k, v := range solution.Labels {
					fmt.Println(k + ": " + v)
				}

			}

		},
	}

	cmd.Flags().String("to", "", "Targer tree where bump new specs.")
	cmd.Flags().Bool("override", false, "Override existing specs if already present.")
	cmd.Flags().Bool("json", false, "JSON output.")
	cmd.Flags().StringP("pkg", "p", "", "Define the package to analyze.")
	cmd.PersistentFlags().Bool("with-portage-pkg", false, "Generate portage packages for every required package.")

	return cmd
}
